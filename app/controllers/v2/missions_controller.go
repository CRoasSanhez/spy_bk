package v2

import (
	"reflect"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"time"

	"gopkg.in/mgo.v2/bson"

	"spyc_backend/app/serializers"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// MissionsController ...
type MissionsController struct {
	BaseController
}

// Show [/v2/games/:id] GET
// Returns mission for the user based on the given id
func (c MissionsController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var mission models.Mission
	if Mission, ok := app.Mapper.GetModel(&mission); ok {

		// Find mission in server
		if err := Mission.Find(id).Exec(&mission); err != nil {
			return c.ErrorResponse(nil, c.Message("error.notFound", "Mission", ""), core.ValidationStatus[core.StatusNotFound])
		}

		// Validate mission status
		if mission.Status.Name != core.StatusActive {
			return c.ErrorResponse(nil, c.Message("error.inactive", "Mission"), core.GameStatus[core.StatusInactive])
		}

		var missionSubscription models.MissionUser
		if MissionUser, ok := app.Mapper.GetModel(&missionSubscription); ok {

			var query = []bson.M{bson.M{"mission_id": mission.GetID()}, bson.M{"user_id": c.CurrentUser.GetID()}}
			MissionUser.FindWithOperator("$and", query).Exec(&missionSubscription)

			// Verify if the users is already subscribed
			if missionSubscription.Status.Name != "" {
				mission.Status = missionSubscription.Status
			}

			var missionTarget models.Target

			if err := missionTarget.GetByOrderAndMission(1, id, false); err != nil {
				return c.ErrorResponse(nil, c.Message("error.notFound", "Targets"), core.ValidationStatus[core.StatusWrong])
			}

			mission.Targets = append(mission.Targets, missionTarget)

			return c.SuccessResponse(mission, "success", core.ModelsType[core.ModelGame], serializers.MissionSerializer{Lang: c.Request.Locale})
		}
	}

	return c.ServerErrorResponse()
}

// ShowAllMissions [/v2/all/games] GET
// returns all missions in the DB
func (c MissionsController) ShowAllMissions() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	// TODO: User rol validation

	var missions []models.Mission

	if Mission, ok := app.Mapper.GetModel(&models.Mission{}); ok {
		if err := Mission.All().Exec(&missions); err != nil {
			return c.ErrorResponse(nil, c.Message("error.notFound", "Missions", ""), 0)
		}

		return c.SuccessResponse(missions, "success", core.ModelsType[core.ModelGame], serializers.MissionSerializer{Lang: c.Request.Locale})
	}

	return c.ServerErrorResponse()
}

// GetSubscribedMissions [/v2/games] GET
// returns all Missions in the DB
func (c MissionsController) GetSubscribedMissions() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var myMissionsIDS []bson.ObjectId
	var gu *models.MissionUser

	//Find user subscriptions
	subscribedMissions := gu.GetActiveSubscriptionsByUser(c.CurrentUser.GetID())
	if len(subscribedMissions) <= 0 {
		return c.ErrorResponse(nil, c.Message("error.notFound", "Current", "Missions"), core.ValidationStatus[core.StatusNotFound])
	}

	for i := 0; i < len(subscribedMissions); i++ {
		myMissionsIDS = append(myMissionsIDS, subscribedMissions[i].MissionID)
	}

	var missions models.Missions

	if Mission, ok := app.Mapper.GetModel(&models.Mission{}); ok {
		// find Mission that fit in user subscriptions
		err := Mission.Query(
			bson.M{
				"$and": []bson.M{
					bson.M{"status.name": core.StatusActive},
					bson.M{"_id": bson.M{"$in": myMissionsIDS}},
				},
			}).Exec(&missions.Missions)
		if err != nil {
			return c.ErrorResponse(nil, c.Message("error.notFound", "Missions", ""), 0)
		}

		return c.SuccessResponse(missions, "success", core.ModelsType[core.ModelGames], serializers.MissionsSerializer{Lang: c.Request.Locale})
	}

	return c.ServerErrorResponse()
}

// ShowMissionsByLocation [/v2/games/geo] POST
// returns the Missions near by user locaation
// -99.13320799999997, 19.4326077
func (c *MissionsController) ShowMissionsByLocation(lat, lng float64) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if lng > 180 || lng < -180 || lat > 90 || lat < -90 {
		return c.ErrorResponse(nil, c.Message("error.invalid", "Coordinates"), core.ValidationStatus[core.StatusError])
	}

	// Get Missions by location
	var missions []models.Mission
	Mission, ok := app.Mapper.GetModel(&models.Mission{})
	if !ok {
		return c.ErrorResponse(nil, c.Message("error.badRequest"), core.ValidationStatus[core.StatusError])
	}
	var near = bson.M{"type": "Point", "coordinates": []float64{lng, lat}}
	var match = bson.M{"status.name": core.StatusActive}
	var pipe = mgomap.Aggregate{}.GeoNear(near, "geolocation", core.MaxDistanceForPin, nil, "geolocation", 0).LookUp("targets", "_id", "missions", "targets").Match(match).Sort(bson.M{"start_date": 1, "title": 1})

	if err := Mission.Pipe(pipe, &missions); err != nil {
		revel.ERROR.Print(err)
		return c.ErrorResponse(nil, c.Message("error.notFound", "Near", "Missions"), core.ValidationStatus[core.StatusNotFound])
	}

	// Get user active subscriptions to Missions to set Subcribed field in response
	ms := models.MissionUser{}
	subscriptions := ms.FindSubscriptionsByUserStatus(c.CurrentUser.GetID())

MainList:
	for i := 0; i < len(missions); i++ {
		for j := 0; j < len(subscriptions); j++ {
			if subscriptions[j].MissionID == missions[i].GetID() {
				missions[i].Subscribed = true
				missions[i].Targets = append([]models.Target{}, FindByID(subscriptions[j].CurrentTarget, missions[i].Targets))
				switch subscriptions[j].Status.Name {

				case core.StatusInit, core.StatusActive:
					missions[i].SetStatus(core.StatusActive)

				case core.StatusIncomplete:
					missions[i].SetStatus(core.StatusIncomplete)

				case core.StatusCompletedNotWinner, core.StatusCompletedWinner, core.StatusCompleted:
					missions[i].SetStatus(core.StatusCompleted)
				case core.StatusPendingValidation:
					missions[i].Status.Name = core.StatusPendingValidation
					missions[i].Status.Code = core.ValidationStatus[core.StatusPendingValidation]
				default:
					missions[i].SetStatus(core.StatusInit)
				}
				continue MainList
			}
		}
		missions[i].Targets = append([]models.Target{}, FindByOrder(1, missions[i].Targets))
	}

	// Get pending rewards
	r := models.Reward{}
	var resources = []string{}
	rewards := r.GetRewardsByUserAndResourceType(c.CurrentUser.GetID().Hex(), resources, true, core.LimitRewards)

	chr := models.CashHuntResponse{
		Missions: missions,
		Rewards:  rewards,
	}
	return c.SuccessResponse(chr, "success", core.ModelsType[core.ModelGame], serializers.CashHuntSerializer{Lang: c.Request.Locale})
}

// SubscribeToMission [/v2/games/:id/subscribe] POST
// Subscribes User to Mission and create UserTarget registry
func (c MissionsController) SubscribeToMission(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	// Validate ID valid
	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	// Mission status Verification
	var err error
	var mission models.Mission

	if Mission, ok := app.Mapper.GetModel(&mission); ok {
		var match = bson.M{"$and": []bson.M{
			bson.M{"_id": bson.ObjectIdHex(id)},
			bson.M{"status.name": core.StatusActive},
		}}
		pipe := mgomap.Aggregate{}.Match(match).LookUp("missions_users", "_id", "mission_id", "subscriptions")

		if err = Mission.Pipe(pipe, &mission); err != nil {
			return c.ErrorResponse(nil, c.Message("error.notFound", "Mission", ""), 400)
		}

		// Validate if mission status is inactive or complete
		if mission.Status.Name == core.StatusInactive || mission.Status.Name == core.StatusCompleted {
			return c.ErrorResponse(nil, c.Message("error.inactive", "Mission"), core.GameStatus[core.StatusInactive])
		}

		// Validate if missions has advertisment and user has viewed it
		if mission.Advertisement != "" {

			var advert models.Advertisement
			if Ad, ok := app.Mapper.GetModel(&advert); ok {
				if err := Ad.Find(mission.Advertisement).Exec(&advert); err != nil {
					revel.ERROR.Printf("ERROR FIND Advertisement %s --- %s", id, err.Error())
					return c.ErrorResponse(nil, c.Message("error.notFound", "Mission", ""), 400)
				}
				if !mission.FindUserInViews(c.CurrentUser.GetID().Hex()).Completed {
					return c.ErrorResponse(c.Message("error.subscription", "Mission"), "User hasn't seen Ad", core.ValidationStatus[core.StatusError])
				}

			}

		}

		// Create Missions_users document for user subscription to mission
		var missionSubscription models.MissionUser
		var subscriptions []models.MissionUser

		if MissionUser, ok := app.Mapper.GetModel(&missionSubscription); ok {
			// Get previous user subscription to mission
			for _, sub := range mission.Subscriptions {
				if sub.UserID == c.CurrentUser.GetID() {
					subscriptions = append(subscriptions, sub)
				}
			}

			// Validte if there're previous subscriptions
			if len(subscriptions) >= core.MaxSubscriptions {
				return c.ErrorResponse(nil, c.Message("error.maxReached", "subscriptions"), core.ValidationStatus[core.StatusMaxSubscriptionsReached])
			}

			// Finding first Target of mission to indicate it in mission_user
			var target models.Target

			if err = target.GetByOrderAndMission(1, id, true); err != nil || len(target.Rewards) <= 0 {
				return c.ErrorResponse(nil, c.Message("error.notFound", "Target", "or Reward"), core.ValidationStatus[core.StatusNotFound])
			}

			// Create mission subcription
			missionSubscription.UserID = c.CurrentUser.GetID()
			missionSubscription.MissionID = mission.GetID()
			missionSubscription.Completed = false
			missionStartDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
			missionSubscription.StartDate = missionStartDate
			missionSubscription.CurrentTarget = target.GetID()

			if err = MissionUser.Create(&missionSubscription); err != nil {
				return c.ErrorResponse(nil, c.Message("error.subscription", "Mission"), 400)
			}

			// Create user subscription to target
			var targetSubscription models.TargetUser

			if TargetUser, ok := app.Mapper.GetModel(&targetSubscription); ok {
				targetSubscription.UserID = c.CurrentUser.GetID()
				targetSubscription.TargetID = target.GetID()
				targetSubscription.Completed = false
				targetSubscription.Intents = 0
				startDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
				targetSubscription.StartDate = startDate

				// If target is type WebGame then find random quiestion and insert it into set Subscription
				if target.Type == core.GameTypeOptions {
					question := target.GetRandomQuestion()

					// Insert the selected one
					target.Question = question
					targetSubscription.QuestionID = question.GetID().Hex()
				}

				if target.Type == core.GameTypeGame {
					targetSubscription.Token = core.StatusInit
				}

				if err = TargetUser.Create(&targetSubscription); err != nil {
					return c.ErrorResponse(nil, c.Message("error.create", "Target"), core.ValidationStatus[core.StatusError])
				}

				return c.SuccessResponse(target, "success", core.ModelsType[core.ModelStep], serializers.TargetSerializer{Lang: c.Request.Locale})
			}
		}
	}

	return c.ServerErrorResponse()
}

// UnsubscribeFromMission [/v2/games/:id/unsubscribe] POST
// unsubscribes an user from a Mission
func (c MissionsController) UnsubscribeFromMission(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var missionSubscription models.MissionUser

	// Validate ID valid
	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	// Find subcription(MissionUser) and set logical deletion
	if err := missionSubscription.GetSubscriptionByMissionUser(id, c.CurrentUser.GetID()); err != nil {
		return c.ErrorResponse(nil, c.Message("error.notFound", "Subscription"), core.ValidationStatus[core.StatusNotFound])
	}

	if err := missionSubscription.CompleteUnsubscription(); err != nil {
		return c.ErrorResponse(nil, c.Message("error.unsubscription", "Mission"), core.ValidationStatus[core.StatusError])
	}

	// Unsubscribe target subscriptions from mission
	if errs := UnsubscribeTargetsByMission(id, c.CurrentUser.GetID()); len(errs) > 0 {
		return c.ErrorResponse(nil, c.Message("error.unsubscription", "Targets"), 400)
	}

	// Create GORUTINE to save User data analytics

	return c.SuccessResponse("success", "success", core.ModelsType[core.ModelSimpleResponse], nil)
}

// AddUserView [/v2/games/:id/ad] PATCH
// Adds an user view for the mission advertisement
func (c MissionsController) AddUserView(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var view models.Views
	var err error
	if err = c.PermitParams(&view, true, "close_time", "completed"); err != nil {
		return c.ErrorResponse(c.Message("error.badRequest"), c.Message("error.badRequest"), 400)
	}

	if Mission, ok := app.Mapper.GetModel(&models.Mission{}); ok {

		var selector = bson.M{
			"_id":         bson.ObjectIdHex(id),
			"status.name": core.StatusActive,
			"ad_views":    bson.M{"$ne": nil},
		}
		var query = bson.M{"$push": bson.M{"ad_views": bson.M{
			"user_id":    c.CurrentUser.GetID().Hex(),
			"close_time": view.CloseTime,
			"completed":  view.Completed,
		},
		}}

		if err = Mission.UpdateQuery(selector, query, false); err != nil {
			revel.ERROR.Printf("ERROR UPDATE MissionViews %s --- %s", id, err.Error())
			return c.ErrorResponse(nil, c.Message("error.update", "Mission"), 400)
		}
		return c.SuccessResponse("success", "success", 200, nil)
	}

	return c.ServerErrorResponse()
}

// UnsubscribeTargetsByMission sets logical deletion to user subscriptions
func UnsubscribeTargetsByMission(missionID string, userID bson.ObjectId) (errs []string) {

	missionTargets := []models.Target{}
	TargetUser, _ := app.Mapper.GetModel(&models.TargetUser{})
	targetSubscriptions := []models.TargetUser{}
	subscriptionsToDelete := []models.TargetUser{}

	var err error
	mission := models.Mission{}
	missionTargets, errTarget := mission.FindTargetsByMission(missionID)

	// Find active User Target subscriptions
	targetSubcription := models.TargetUser{}
	targetSubscriptions, errMission := targetSubcription.FindSubscriptionsByUser(userID)

	if errTarget != nil || errMission != nil {
		return append(errs, errTarget.Error(), errMission.Error())
	}

	// Iterate in Targets of mission
	for i := 0; i < len(missionTargets); i++ {

		// Iterate into user Targets subrcriptions
		for j := 0; j < len(targetSubscriptions); j++ {

			// Validate the same Target ID
			if missionTargets[i].GetID().Hex() == targetSubscriptions[j].GetID().Hex() {
				subscriptionsToDelete = append(subscriptionsToDelete, targetSubscriptions[j])
			}
		}
	}

	// Set active missionTargets status to incomplete
	for _, subsc := range subscriptionsToDelete {
		subsc.SetStatus(core.StatusUnsubscribed)
		if err = TargetUser.Update(&subsc); err != nil {
			errs = append(errs, err.Error())
		}
	}
	return errs
}

// FindByID ...
func FindByID(value bson.ObjectId, array []models.Target) models.Target {
	for i := 0; i < len(array); i++ {
		if array[i].GetID() == value {
			return array[i]
		}
	}
	return models.Target{}
}

func FindByOrder(value int, array []models.Target) models.Target {
	for i := 0; i < len(array); i++ {
		if array[i].Order == value {
			return array[i]
		}
	}
	return models.Target{}
}

func FindByProperty(value interface{}, array []models.Target, field string) models.Target {

	v := reflect.ValueOf(value)

	for i := 0; i < len(array); i++ {
		r := reflect.ValueOf(&array[i])
		f := reflect.Indirect(r).FieldByName(field)

		switch v.Kind() {
		case reflect.Int:
			if f.Int() == v.Int() {
				return array[i]
			}
		case reflect.Int32:
		case reflect.Int64:
		case reflect.Uint:
			if f.Uint() == v.Uint() {
				return array[i]
			}
		case reflect.Float32:
			if f.Float() == v.Float() {
				return array[i]
			}
		case reflect.Float64:
			if f.Float() == v.Float() {
				return array[i]
			}
		case reflect.String:
			if f.String() == v.String() {
				return array[i]
			}
		case reflect.Bool:
			if f.Bool() == v.Bool() {
				return array[i]
			}
		}

		if array[i].Order == value {
			return array[i]
		}
	}

	return models.Target{}
}
