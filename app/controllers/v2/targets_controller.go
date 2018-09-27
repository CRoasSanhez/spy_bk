package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"spyc_backend/app/auth"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// TargetsController ...
type TargetsController struct {
	BaseController
}

// Show [/v2/steps/:id] GET
// gets the Target by the given id
func (c TargetsController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var target models.Target

	if Target, ok := app.Mapper.GetModel(&target); ok {
		// Validate ID valid
		if !bson.IsObjectIdHex(id) {
			return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
		}

		var query = []bson.M{
			bson.M{"_id": bson.ObjectIdHex(id)},
			bson.M{"status.name": core.StatusActive},
		}

		// Get target from DB
		if err := Target.FindWithOperator("$and", query).Exec(&target); err != nil {
			return c.ErrorResponse(nil, c.Message("error.notFound", "Target", ""), core.ValidationStatus[core.StatusNotFound])
		}

		// If Target is type options the go find options to retrieve in message
		if target.Type == core.GameTypeOptions {
			var targetSubcription models.TargetUser
			var question models.Question

			if err := targetSubcription.GetSubscription(id, c.CurrentUser.GetID(), core.StatusInit); err != nil {
				return c.ErrorResponse(nil, c.Message("error.notSubscribed", "Target"), core.ValidationStatus[core.StatusNotFound])
			}

			if err := question.GetQuestionByID(targetSubcription.QuestionID); err != nil {
				return c.ErrorResponse(nil, c.Message("error.retreive", "Questions"), core.ValidationStatus[core.StatusNotFound])
			}

			target.Question = question
		}

		return c.SuccessResponse(target, "success", core.ModelsType[core.ModelStep], serializers.TargetSerializer{Lang: c.Request.Locale})
	}

	return c.ServerErrorResponse()
}

// GetCurrentTargetByMissionID [/v2/steps/:idMission/current] POST
// returns the current active target based on the given MissionID
func (c TargetsController) GetCurrentTargetByMissionID(idMission string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	// Validate ID valid
	if !bson.IsObjectIdHex(idMission) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	// Get active targets from Mission
	var mission models.Mission
	var missionTargets []models.Target

	if m, err := mission.FindTargetsByMission(idMission); err != nil {
		return c.ErrorResponse(nil, c.Message("error.notFound", "Targets", ""), core.ValidationStatus[core.StatusNotFound])
	} else {
		missionTargets = m
	}

	// with the given id and return current Target
	var missionSubscription models.MissionUser

	// Get Mission Subscription
	if err := missionSubscription.GetSubscriptionByMissionUser(idMission, c.CurrentUser.GetID()); err != nil {
		return c.ErrorResponse(nil, c.Message("error.notSubscribed", "Mission"), core.ValidationStatus[core.StatusWrong])
	}

	// Find user Targets subscription in status
	// init, active,penidng Validation
	var targetSubcription models.TargetUser
	targetSubscriptions, errMission := targetSubcription.FindSubscriptionsByUser(c.CurrentUser.GetID())

	if errMission != nil {
		return c.ErrorResponse(nil, c.Message("error.notSubscribed", "Targets"), core.ValidationStatus[core.StatusNotFound])
	}

	var currentTarget models.Target

	// Get Target from array of Targets of mission
	for _, target := range missionTargets {
		for _, targetSubcription := range targetSubscriptions {
			// Verify if the Target from mission is equal to the targetSubcription ID
			if targetSubcription.TargetID == target.GetID() {
				currentTarget = target
				break
			}
		}
	}

	// Validate Target found
	if currentTarget.Type == "" {
		return c.ErrorResponse(nil, c.Message("error.notSubscribed", "Any target of this Mission"), core.ValidationStatus[core.StatusNotFound])
	}

	if currentTarget.Status.Name == core.StatusPendingValidation {
		if currentTarget.NextTargetID != "" {
			Target, _ := app.Mapper.GetModel(&currentTarget)
			Target.Find(currentTarget.NextTargetID).Exec(&currentTarget)
		}
	}
	return c.SuccessResponse(currentTarget, "success", core.ModelsType[core.ModelStep], serializers.TargetSerializer{Lang: c.Request.Locale})
}

// SubscribeToTarget [/v2/steps/:id/subscribe] POST
// Subscribes an user to the given target
func (c TargetsController) SubscribeToTarget(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var target models.Target
	if Target, ok := app.Mapper.GetModel(&target); ok {

		var match = bson.M{"status.name": core.StatusActive, "_id": bson.ObjectIdHex(id)}
		var match2 = bson.M{"subscriptions.user_id": c.CurrentUser.GetID()}
		var pipe = mgomap.Aggregate{}.Match(match).LookUp("targets_users", "_id", "target_id", "subscriptions").Match(match2)

		if err := Target.Pipe(pipe, &target); err != nil {
			revel.ERROR.Print(err)
			return c.ErrorResponse(nil, c.Message("error.notFound", "Target", ""), core.ValidationStatus[core.StatusNotFound])
		}

		if len(target.Subscriptions) >= core.MaxSubscriptions {
			return c.ErrorResponse(nil, c.Message("error.maxReached", "Subscriptions"), core.GameStatus[core.StatusMaxIntentsReached])
		}

		var questionID string
		var subscription models.TargetUser

		if target.Type == core.GameTypeOptions {
			question := target.GetRandomQuestion()
			target.Question = question
			questionID = target.Question.GetID().Hex()
		}

		if target.Type == core.GameTypeGame {
			subscription.Token = core.StatusInit
		}

		subscription.TargetID = bson.ObjectId(id)
		subscription.UserID = c.CurrentUser.GetID()
		subscription.Intents = 0
		subscription.QuestionID = questionID
		startDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
		subscription.StartDate = startDate

		if err := subscription.CreateSubscription(); err == nil {

			if err = models.StatsUpdateFields(bson.M{"missions.subscribed": 1, "missions.total": 1}, []string{"missions.first_access_date", "missions.last_access"}, c.CurrentUser.GetID()); err != nil {
			}

			return c.SuccessResponse(target, "success", core.ModelsType[core.ModelStep], serializers.TargetSerializer{Lang: c.Request.Locale})
		}
	}
	return c.ServerErrorResponse()
}

// UnsubscribeFromTarget [/v2/steps/:id/unsubscribe] POST
// Ussubcribe an user from the given target
func (c TargetsController) UnsubscribeFromTarget(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	subscription := models.TargetUser{}
	if TargetUser, ok := app.Mapper.GetModel(&subscription); ok {
		var selector = bson.M{"$and": []bson.M{
			bson.M{"user_id": c.CurrentUser.GetID()},
			bson.M{"target_id": bson.ObjectIdHex(id)},
			bson.M{"status.name": core.StatusInit},
		}}
		endDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
		var query = bson.M{"$set": bson.M{
			"completed":   true,
			"status.name": core.StatusInactive,
			"status.code": core.SubscriptionStatus[core.StatusInactive],
			"end_date":    endDate,
		}}

		if err := TargetUser.UpdateQuery(selector, query, false); err != nil {
			revel.ERROR.Printf("ERROR UPDATE TargetUser Target: %s --- %s", id, err.Error())
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if err := models.StatsUpdateFields(bson.M{"missions.uncompleted": 1}, []string{"missions.last_uncompleted", "missions.last_access"}, c.CurrentUser.GetID()); err != nil {
		}

		return c.SuccessResponse("success", "success", core.ModelsType[core.ModelSimpleResponse], nil)
	}
	return c.ServerErrorResponse()
}

// ValidateReward [/v2/validate/:idMission/:idTarget/:data] POST
// validates the given answer  for the target
func (c TargetsController) ValidateReward(file []byte, idMission, idTarget, data string) revel.Result {

	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(idMission) || !bson.IsObjectIdHex(idTarget) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var score string
	var success bool
	var err error
	var target models.Target

	if Target, ok := app.Mapper.GetModel(&target); ok {

		if target, err = target.FindTargetAndSubscriptions(idTarget, core.StatusActive, c.CurrentUser.GetID()); err != nil {
			return c.ErrorResponse(nil, c.Message("error.notFound", "Target", ""), core.ValidationStatus[core.StatusNotFound])
		}

		var mission models.Mission
		if mission, err = mission.FindMissionAndSubscriptions(idMission, core.StatusActive, c.CurrentUser.GetID()); err != nil {
			return c.ErrorResponse(nil, c.Message("error.notFound", "Mission", ""), core.ValidationStatus[core.StatusNotFound])
		}

		if len(mission.Subscriptions) > core.MaxSubscriptions {
			return c.ErrorResponse(nil, c.Message("error.maxReached", "Mission Subscriptions"), core.ValidationStatus[core.StatusMaxSubscriptionsReached])
		}

		var subscription = target.Subscriptions[0]

		TargetUser, ok2 := app.Mapper.GetModel(&subscription)
		if !ok2 {
			revel.ERROR.Print("ERROR MAPPING TargetUser")
			return c.ServerErrorResponse()
		}

		// Update last_access
		if err = models.StatsUpdateFields(nil, []string{"missions.last_access"}, c.CurrentUser.GetID()); err != nil {
		}

		// Validate max number ob subscriptions and update Status if true
		if subscription.Intents >= core.MaxGameIntents || subscription.Intents < 0 {
			subscription.CompleteSubscription(core.StatusIncomplete)

			if err = models.StatsUpdateFields(bson.M{"missions.uncompleted": 1}, []string{"missions.last_uncompleted"}, c.CurrentUser.GetID()); err != nil {
			}

			if err := mission.Subscriptions[0].CompleteSubscriptionByMissionUser(idMission, c.CurrentUser.GetID().Hex(), core.StatusCompleted); err != nil {
				return c.ErrorResponse(c.Message("error.update", "Subscription Status"), c.Message("error.update", "Subscription Status"), core.ValidationStatus[core.StatusError])
			}
			return c.ErrorResponse(c.Message("error.maxReached", "Attempts"), c.Message("error.maxReached", "Attempts"), core.GameStatus[core.StatusMaxIntentsReached])
		}

		if len(target.Subscriptions) > core.MaxSubscriptions {

			if subscription.Status.Name == core.StatusCompleted {
				return c.ErrorResponse(c.Message("error.complete"), c.Message("error.complete"), core.ValidationStatus[core.StatusCompleted])
			}
			return c.ErrorResponse(nil, c.Message("error.maxReached", "Subscriptions"), core.GameStatus[core.StatusMaxIntentsReached])
		}

		// Receive request by type
		switch target.Type {
		case core.GameTypeNIP, core.GameTypeText, core.GameTypeOptions:
			var answerReq models.AnswerRequest

			if err := c.Params.BindJSON(&answerReq); err != nil {
				revel.ERROR.Print("Error binding general request")
				return c.ErrorResponse(err, err.Error(), 400)
			}

			if target.Type == core.GameTypeNIP || target.Type == core.GameTypeText {

				TargetUser, ok := app.Mapper.GetModel(&subscription)
				if !ok {
					revel.ERROR.Print("Error Mapper.GetModel")
					return c.ErrorResponse("Internal Error", "Internal Error ", core.ValidationStatus[core.StatusError])
				}
				subscription.Score = answerReq.Data
				subscription.SetStatus(core.StatusPendingValidation)

				if err := TargetUser.Update(&subscription); err != nil {
					revel.ERROR.Print("Error updating subscription: " + subscription.GetID())
					return c.ErrorResponse(err, err.Error(), core.ValidationStatus[core.StatusError])
				}

				if err = mission.Subscriptions[0].CompleteSubscription(core.StatusPendingValidation); err != nil {
					return c.ServerErrorResponse()
				}

				core.Notify(core.ValidateEntry, "Answer Validation ", "There is an answer for validation", "Mission: "+mission.Title.GetString(c.Request.Locale),
					core.GetDashboardPath()+"validation", "", "Target: "+target.Name.GetString(c.Request.Locale), []interface{}{
						bson.M{"title": "Target Answer: " + target.Score, "value": "User Answer: " + subscription.Score, "short": true},
					})

				return c.SuccessResponse("success", "Your answer is in validation", core.ValidationStatus[core.StatusPendingValidation], nil)
			}

			if target.Type == core.GameTypeOptions {
				success = ValidateOption(subscription.QuestionID, &answerReq)
			}

			break
		case core.GameTypePhoto:
			// Validate user's location
			_, err := auth.GetHeaderCoordinates(c.Request)
			if err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}

			//if !ValidateAnswerLocation(geo.Coordinates, target.Geolocation.Coordinates) {
			//	break
			//}

			if len(file) <= 0 {
				return c.ErrorResponse(nil, "File not found", 400)
			}

			// Upoload multipart file (normal image)
			if err = subscription.ResponseImage.Init(models.AsDocumentBase(&subscription), c.Params.Files["file"][0]); err != nil {
				revel.ERROR.Print("ERROR INIT Image --- " + err.Error())
				return c.ErrorResponse(err, err.Error(), 400)
			}

			if err = subscription.ResponseImage.Upload(); err != nil {
				revel.ERROR.Print("ERROR UPLOAD Image --- " + err.Error())
				return c.ErrorResponse(err, err.Error(), 400)
			}

			subscription.SetStatus(core.StatusPendingValidation)

			// Only for This version change mission subscription to Pending validation
			if target.NextTargetID != "" {
				var targetResponse models.Target
				if err = Target.Find(target.NextTargetID).Exec(&targetResponse); err != nil {
					return c.ErrorResponse(c.Message("error.notFound", "Next Target"), c.Message("error.notFound", "Next Target"), core.ValidationStatus[core.StatusNotFound])
				}

				return c.SuccessResponse(targetResponse, "success", core.ModelsType[core.ModelStep], serializers.TargetSerializer{Lang: c.Request.Locale})
			}
			if err = mission.Subscriptions[0].CompleteSubscription(core.StatusPendingValidation); err != nil {
				return c.ServerErrorResponse()
			}

			if err = TargetUser.Update(&subscription); err != nil {
				revel.ERROR.Print("ERROR MAPPING TargetUser")
				return c.ServerErrorResponse()
			}

			core.Notify(core.ValidateEntry, "Photo Validation ", "There is a Photo for validation", "Mission: "+mission.Title.GetString(c.Request.Locale),
				core.GetDashboardPath()+"validation", "", "Target: "+target.Name.GetString(c.Request.Locale), []interface{}{
					bson.M{"title": "Target Score: " + target.Score, "value": "User Score: " + subscription.Score, "short": true},
				})

			return serializers.SuccessResponse("success", "Your answer is in validation", core.ValidationStatus[core.StatusPendingValidation], nil)
			break
		case core.GameTypeCheckIn:
			if geo, err := auth.GetHeaderCoordinates(c.Request); err != nil {
				revel.ERROR.Print(err.Error())
				return c.ErrorResponse(err, err.Error(), 400)
			} else {
				var answerRequest models.AnswerRequest

				if err := c.Params.BindJSON(&answerRequest); err != nil {
					success = false
					return c.ErrorResponse(err, err.Error(), 400)
				}

				answerRequest.Geolocation = geo

				success = ValidateAnswerLocation(geo.Coordinates, target.Geolocation.Coordinates)
			}

			break
		case core.GameTypeQR:
			if _, err := auth.GetHeaderCoordinates(c.Request); err != nil {
				revel.ERROR.Print(err.Error())
				return c.ErrorResponse(err, err.Error(), 400)
			} else {
				// VAlidate if user makes the checkin in the place
				//if !ValidateAnswerLocation(geo.Coordinates, target.Geolocation.Coordinates) {
				//	break
				//}
				// Validate data parameter in URL
				success = ValidateQRAnswer(data, mission.GetID().Hex(), target.GetID().Hex())
			}

			break
			//
			// For this target just validate token
			// and returns new token and gameInfo
		case core.GameTypeGame:
			if subscription.Token != core.StatusInit {
				return c.ErrorResponse(c.Message("error.invalid", "Game"), c.Message("error.invalid", "Game"), core.ValidationStatus[core.StatusError])
			}

			subscription.Score = target.Score
			gameInfo, errG := GetGameInfo(subscription, c.CurrentUser)
			if errG != nil {
				return c.ErrorResponse(nil, c.Message("error.generateNew", "Token"), core.ValidationStatus[core.StatusError])
			}

			return c.SuccessResponse(gameInfo, "success", core.ModelsType[core.ModelWebGameInfo], nil)
		default:
			success = false
			break
		}

		// Update lastResult
		if err = models.StatsSetField(c.CurrentUser.GetID(), bson.M{"missions.last_result": strconv.FormatBool(success)}); err != nil {
			revel.ERROR.Printf("ERROR UPDATE last_access for owner: %s --- %s", c.CurrentUser.GetID().Hex(), err.Error())
		}

		// If unsuccessfully validation
		if !success {
			if err := subscription.AddIntent(); err != nil {
				revel.ERROR.Printf("ERROR ADD Target Intent: %s --- %s", target.GetID(), err.Error())
				return c.ErrorResponse(err, err.Error(), 400)
			}
			return c.ErrorResponse(c.Message("error.wrongAnswer"), c.Message("error.wrongAnswer"), core.ValidationStatus[core.StatusWrong])
		}

		if err = models.StatsUpdateFields(bson.M{"missions.completed": 1}, []string{"missions.last_completed"}, c.CurrentUser.GetID()); err != nil {
		}

		// Change current target subscription to complete
		if err = subscription.CompleteSubscription(core.StatusCompleted); err != nil {
			return serializers.ErrorResponse(c.Message("error.update", "Subscription status"), c.Message("error.update", "Subscription status"), core.ValidationStatus[core.StatusError])
		}

		// If next_target_id exists return next target
		if target.NextTargetID != "" {

			var targetResponse models.Target
			if err := Target.Find(target.NextTargetID).Exec(&targetResponse); err != nil {
				return c.ErrorResponse(c.Message("error.notFound", "Next Target"), c.Message("error.notFound", "Next Target"), core.ValidationStatus[core.StatusNotFound])
			}

			return c.SuccessResponse(targetResponse, "success", core.ModelsType[core.ModelStep], serializers.TargetSerializer{Lang: c.Request.Locale})
		}

		// Get available reward for the current Target
		var reward models.Reward
		rewards := reward.GetRewardsByResource(idTarget)
		if len(rewards) == 0 {
			return c.SuccessResponse(c.Message("error.notAvailable", "Reward"), "success", core.ModelsType[core.ModelSimpleResponse], nil)
		}

		reward = rewards[0]
		win := models.Win{}

		_, errWinner := win.CreateWin(reward.GetID(), c.CurrentUser.GetID(), core.WinTypeCashHunt, score)
		if errWinner != nil {
			return c.ErrorResponse(c.Message("error.create", "Winner"), c.Message("error.create", "Winner"), core.GameStatus[core.StatusInactive])
		}

		// if reward is active then change its status to complete
		if ok := reward.CompleteReward(c.CurrentUser.GetID().Hex()); !ok {
			return c.ErrorResponse(c.Message("error.update", "Reward status"), c.Message("error.update", "Reward status"), core.GameStatus[core.StatusInactive])
		}

		// Change missionSubscription to complete
		if err = mission.Subscriptions[0].CompleteSubscriptionByMissionUser(idMission, c.CurrentUser.GetID().Hex(), core.StatusCompleted); err != nil {
			return serializers.ErrorResponse(c.Message("error.update", "Mission status"), c.Message("error.update", "Mission status"), core.ValidationStatus[core.StatusError])
		}

		reward.GetCoinsByMission(idMission)

		if (len(rewards) == 1 && !reward.Multi) || (reward.Multi && reward.Status.Name == core.StatusCompleted) {
			if mission.Complete() {
				c.NotifyMissionComplete(mission)
			}
		}

		// Update rewardID to winner subscription
		if err = TargetUser.UpdateQuery(bson.M{"_id": subscription.GetID()}, bson.M{"$set": bson.M{"reward_id": reward.GetID()}}, false); err != nil {
			revel.ERROR.Printf("ERROR UPDATE Target Subscription: %s --- %s", subscription.GetID(), err.Error())
		}

		// Notify to Slack that the user has obtained the reward
		core.Notify(core.NewEntry, "Reward Obtained", "A Reward was Obtained", mission.Title.GetString(c.CurrentUser.Device.Language),
			core.GetDashboardPath()+"users/"+c.CurrentUser.GetID().Hex(), "", mission.Description.GetString(c.CurrentUser.Device.Language), []interface{}{
				bson.M{"title": mission.Type, "value": core.ConcatArray(mission.Countries), "short": true},
			})

		return c.SuccessResponse(reward, "success", core.ModelsType[core.ModelReward], serializers.RewardSerializer{Lang: c.Request.Locale})
	}

	return c.ServerErrorResponse()
}

// GetPreviousSubscriptions Gets user subscriptions from Target
func GetPreviousSubscriptions(idUser, idTarget bson.ObjectId) (subscriptions []models.TargetUser, ok bool) {

	TargetUser, _ := app.Mapper.GetModel(&models.TargetUser{})
	//mStatus := []string{core.StatusActive, core.StatusInit}

	// Find all user subscriptions
	err := TargetUser.Query(
		bson.M{
			"$and": []bson.M{
				bson.M{"target_id": idTarget},
				bson.M{"user_id": idUser},
				//bson.M{"completed": false},
				//bson.M{"status.name": bson.M{"$in": mStatus}},
			},
		}).Exec(&subscriptions)
	if err != nil {
		return subscriptions, false
	}
	return subscriptions, true
}

// ValidateOption validates target answer type multiple option
func ValidateOption(idQuestion string, answerReq *models.AnswerRequest) bool {

	question := models.Question{}
	Question, _ := app.Mapper.GetModel(&question)

	// Find answer reference in Document
	if err := Question.Find(idQuestion).Exec(&question); err != nil {
		revel.ERROR.Print(err.Error())
		return false
	}
	// Compare both answers DB vs Request
	for _, o := range question.Options {
		if o.IsAnswer && strconv.Itoa(o.ID) == answerReq.Data {
			return true
		}
	}
	return false
}

// ValidateAnswerLocation validates if the user coordinates are close to the answer coordinates
// Coordinates are [longitud,latitud]
// in GPS coordinates (+/-) 0.00015 ~ 15m ::: (+/-) 0.0005 ~ 50m
func ValidateAnswerLocation(userLocation, targetCoords []float64) (success bool) {

	maxLong := targetCoords[0] + 0.0005
	minLong := targetCoords[0] - 0.0005
	maxLat := targetCoords[1] + 0.0005
	minLat := targetCoords[1] - 0.0005

	if userLocation[0] > maxLong || userLocation[0] < minLong {
		revel.ERROR.Print("longitud not in range")
		return false
	}
	if userLocation[0] > maxLat || userLocation[1] < minLat {
		revel.ERROR.Print("latitud not in range")
		return false
	}
	return true
}

// ValidateQRAnswer validates tue data parameter from url
func ValidateQRAnswer(data, missionID, targetID string) (ok bool) {

	if data == "" {
		return false
	}
	decryptedStr := core.DecryptCipher(data)
	dataArr := strings.Split(decryptedStr, "-")

	if len(dataArr) != 3 {
		return false
	}
	if dataArr[0] != missionID || dataArr[1] != targetID || dataArr[2] != core.GameTypeQR {
		return false
	}
	return true
}

// GetGameInfo return the gameInfo and the token needed for validations
func GetGameInfo(targetSubscription models.TargetUser, user models.User) (gameInfo models.WebGameInfo, err error) {

	newToken := ""
	targetSubscriptionID := targetSubscription.GetID().Hex()
	userID := user.GetID().Hex()
	TargetUser, _ := app.Mapper.GetModel(&targetSubscription)

	// Create the token with action active
	newToken, err = auth.GenerateGameToken(targetSubscriptionID, userID, core.StatusActive)
	if err != nil {
		revel.ERROR.Print(err.Error())
		return
	}

	targetSubscription.Token = newToken
	if err = TargetUser.Update(&targetSubscription); err != nil {
		revel.ERROR.Print(err)
		return
	}

	// Create the data of the game
	gameData := models.GameData{
		Score: targetSubscription.Score, Lang: user.Device.Language, Intents: strconv.Itoa(core.MaxCashHuntGamesIntent),
	}

	playerInfo := models.PlayerInfo{UserName: user.UserName}

	// Create game info
	gameInfo = models.WebGameInfo{
		Player: playerInfo, GameData: gameData, Token: newToken,
	}
	return
}
