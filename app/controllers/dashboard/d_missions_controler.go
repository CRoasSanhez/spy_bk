package dashboard

import (
	"errors"
	"math"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// DMissionsController ...
type DMissionsController struct {
	DBaseController
}

// Index [/spyc_admin/missions] GET
// return all missions on DB
func (c DMissionsController) Index(page int, quantity int, search string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var missions []models.Mission
	var err error
	if Mission, ok := app.Mapper.GetModel(&models.Mission{}); ok {
		if page <= 0 {
			page = 1
		}

		if quantity <= 0 {
			quantity = 10
		}

		// err = Mission.Query(bson.M{"title": bson.M{"$regex": search}}).Sort([]string{"start_date", "title"}).Paginate(page-1, quantity).Exec(&missions)
		var match = bson.M{"$and": []bson.M{
			//{"title": bson.M{"$regex": search}},
			{"status.name": bson.M{"$ne": "inactive"}},
		}}

		var pipe = mgomap.Aggregate{}.Match(match).Sort(bson.M{"start_date": 1, "title": 1}).Skip((page - 1) * quantity).Limit(quantity)

		if err = Mission.Pipe(pipe, &missions); err != nil {
			revel.ERROR.Print(err)
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(missions, "success", 200, serializers.MissionSerializer{})
		}

		//pipe = mgomap.Aggregate{}.Match(match).Sort(bson.M{"start_date": 1, "title": 1}).Count("total")
		total := &struct{ Total int }{}

		if err = Mission.Pipe(pipe, &total); err != nil && err.Error() != "Not found" {
			revel.ERROR.Print(err)
			return c.ErrorResponse(err, err.Error(), 400)
		}

		total.Total = len(missions) // TEST THIS

		var languages []models.Language
		if Language, ok := app.Mapper.GetModel(&models.Language{}); ok {
			if err = Language.FindWithOperator("$and", []bson.M{bson.M{"status.name": core.StatusActive}}).Exec(&languages); err != nil {
				revel.ERROR.Print("ERROR FIND Languages ---" + err.Error())
			}
		}

		var countries []models.Country
		if Country, ok := app.Mapper.GetModel(&models.Country{}); ok {

			pipe = mgomap.Aggregate{}.Sort(bson.M{"name": 1})
			if err = Country.Pipe(pipe, &countries); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}

			var campaigns []models.Campaign
			if Campaign, ok := app.Mapper.GetModel(&models.Campaign{}); ok {

				var query = []bson.M{
					bson.M{"status.name": core.StatusActive},
				}

				if err = Campaign.FindWithOperator("$and", query).Exec(&campaigns); err != nil {
					return c.ErrorResponse(err, err.Error(), 400)
				}

				// Frequency
				var frequency = map[int]string{
					1: "Once a day",
					2: "Twice a day",
					3: "Three times a day",
					4: "Four times a day",
				}

				// Priority
				var priority = map[int]string{
					1: "One",
					2: "Two",
					3: "Three",
					4: "Four",
				}

				c.ViewArgs["JSapikey"] = core.GoogleMapsJSAPIKey
				c.ViewArgs["Missions"] = missions
				c.ViewArgs["Frequency"] = frequency
				c.ViewArgs["Priority"] = priority
				c.ViewArgs["MissionTypes"] = core.MissionTypes
				c.ViewArgs["Campaigns"] = campaigns
				c.ViewArgs["Countries"] = countries
				c.ViewArgs["Languages"] = languages
				c.ViewArgs["CurrentPage"] = page - 1
				if quantity > 0 {
					c.ViewArgs["Pages"] = int(math.Ceil(float64(total.Total) / float64(quantity)))
				} else {
					c.ViewArgs["Pages"] = 0
				}

				return c.Render()
			}
		}
	}

	return c.ServerErrorResponse()
}

// New [/spyc_admin/missions/new] GET
func (c DMissionsController) New() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	return c.Render()
}

// Create [/spyc_admin/missions] POST
func (c DMissionsController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var mission models.Mission
	var err error
	var now = time.Now()

	if c.Params.Get("format") == "json" {
		if err := c.PermitParams(&mission, "title", "description", "type", "geolocation", "start_date", "end_date", "country", "campaign", "frequency"); err != nil {
			revel.ERROR.Print("Error Binding Mission")
			return c.ErrorResponse(err, err.Error(), 400)
		}
	} else {
		lat, _ := strconv.ParseFloat(c.Params.Get("lat"), 64)
		lng, _ := strconv.ParseFloat(c.Params.Get("lng"), 64)
		stime, err1 := time.Parse(core.MXTimeFormatTZ, c.Params.Form.Get("start_date")+" "+c.Params.Form.Get("start_time"))
		etime, err2 := time.Parse(core.MXTimeFormatTZ, c.Params.Form.Get("end_date")+" "+c.Params.Form.Get("end_time"))

		if err1 != nil || err2 != nil || stime.Before(now) || etime.Before(stime) {
			revel.ERROR.Print("Error Parsing date", c.Params.Form.Get("start_date")+" "+c.Params.Form.Get("start_time"))
			revel.ERROR.Print(c.Params.Form.Get("end_date") + " " + c.Params.Form.Get("end_time"))
			return c.Redirect(DMissionsController.Index)
		}

		mission.Type = c.Params.Form.Get("type")
		mission.Langs = append(mission.Langs, c.Params.Get("language"))
		mission.Title = map[string]interface{}{c.Params.Get("language"): c.Params.Get("title")}
		mission.Description = map[string]interface{}{c.Params.Get("language"): c.Params.Get("description")}
		mission.StartDate = core.ChangeUTCTimeToLocalZone(stime)
		mission.EndDate = core.ChangeUTCTimeToLocalZone(etime)
		mission.Countries = strings.Split(c.Params.Form.Get("countries"), ",")

		if len(mission.Countries) <= 0 {
			revel.ERROR.Print("No countries selected")
			return c.Redirect(DMissionsController.Index)
		}

		mission.Campaign = bson.ObjectIdHex(c.Params.Form.Get("campaign"))
		frequency, _ := strconv.Atoi(c.Params.Form.Get("frequency"))
		mission.Frequency = frequency
		priority, _ := strconv.Atoi(c.Params.Form.Get("priority"))
		mission.Priority = priority
		mission.Geolocation = &models.Geo{
			Type:        "Point",
			Coordinates: []float64{lng, lat},
		}
	}

	if len(c.Params.Files) <= 0 {
		revel.ERROR.Print("ERROR PARAMS File not found")
		return c.Redirect(DMissionsController.Index)
	}

	// From 09hrs to 21hrs = 12hrs
	periods := []models.Period{}
	startDate := mission.StartDate

	for i := 0; i < mission.Frequency; i++ {

		// Add frequncy in hours
		if i != 0 {
			startDate = startDate.Add(+time.Duration(mission.Frequency) * time.Hour)
		}

		// append into periods slice
		periods = append(periods, models.Period{
			StartDate: startDate,
			EndDate:   startDate.Add(+time.Duration(1 * time.Hour)),
		})
	}
	mission.Periods = periods

	if Mission, ok := app.Mapper.GetModel(&mission); ok {
		if err = Mission.Create(&mission); err != nil {
			revel.ERROR.Print("ERROR CREATING Mission --- " + err.Error())
			return c.Redirect(DMissionsController.Index)
		}

		// UPLOAD IMAGE
		// Open Multipart
		// var byteArray = c.ResizeImage(400, "cover_picture", core.BlurZero)
		// if len(byteArray) <= 0 {
		// 	return c.Redirect(DMissionsController.Index)
		// }

		if mission.Attachment.PATH != "" {
			mission.Attachment.Remove()
		}

		if err = mission.Attachment.Init(models.AsDocumentBase(&mission), c.Params.Files["cover_picture"][0]); err != nil {
			return c.Redirect(DMissionsController.Index)
		}

		if err = mission.Attachment.Upload(); err != nil {
			return c.Redirect(DMissionsController.Index)
		}

		// if err = mission.Attachment.UploadBytes(owner, byteArray, c.Params.Files["cover_picture"][0].Filename); err != nil {
		// 	revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
		// 	return c.Redirect(DMissionsController.Index)
		// }

		// Validate there's an advertisement
		if c.Params.Form.Get("advertisement") != "" {
			mission.Advertisement = c.Params.Form.Get("advertisement")
		}

		if err = Mission.Update(&mission); err != nil {
			revel.ERROR.Printf("ERROR UPDATE Mission: %s --- %s", mission.GetID().Hex(), err.Error())
			return c.Redirect(DMissionsController.Index)
		}

		if mission.Priority == 3 {
			var topicMision = mission.GetDocumentName() + mission.GetID().Hex()
			n := models.Notification{}
			n.SendToTopic(topicMision, mission.GetID().Hex(), []string{})
		}
		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(mission, "success", 200, serializers.MissionSerializer{})
		}

		return c.Redirect("missions/%s", mission.GetID().Hex())
	}

	return c.Redirect(DMissionsController.Index)
}

// Show [/spyc_admin/missions/:id] GET
// returns the VIEW for Mission Details
func (c DMissionsController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var mission models.Mission
	var err error
	if Mission, ok := app.Mapper.GetModel(&mission); ok {
		if err = Mission.Find(id).Exec(&mission); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(mission, "success", 200, serializers.MissionSerializer{Lang: c.Request.Locale})
		}

		// target types
		targetTypes := map[string]string{
			core.GameTypeCheckIn: "CheckIn",
			core.GameTypeGame:    "Game",
			core.GameTypeOptions: "Question",
			core.GameTypePhoto:   "Photo",
			core.GameTypeQR:      "QR image",
			core.GameTypeText:    "Text",
		}

		var gamesTypes []models.WebGame

		if WebGames, ok := app.Mapper.GetModel(&models.WebGame{}); ok {
			if err := WebGames.FindBy("status.name", core.StatusActive).Exec(&gamesTypes); err != nil {
				return c.Redirect(DMissionsController.Index)
			}

			// avilable target per Mission
			var g models.Mission
			var targets []models.Target

			if t, err := g.GetAllTargetsByMission(id); err != nil {
				revel.ERROR.Print(err)
			} else {
				targets = t
			}

			var languages []models.Language
			if Language, ok := app.Mapper.GetModel(&models.Language{}); ok {
				if err = Language.FindWithOperator("$and", []bson.M{bson.M{"status.name": core.StatusActive}}).Exec(&languages); err != nil {
					revel.ERROR.Print("ERROR FIND Languages ---" + err.Error())
				}
			}

			// Google maps JS API key
			jsapikey := core.GoogleMapsJSAPIKey

			c.ViewArgs["Mission"] = mission
			c.ViewArgs["Types"] = targetTypes
			c.ViewArgs["WebGames"] = gamesTypes
			c.ViewArgs["Targets"] = targets
			c.ViewArgs["Languages"] = languages
			c.ViewArgs["JSapikey"] = jsapikey

			return c.Render()
		}
	}

	return c.ServerErrorResponse()
}

// Update [/spyc_admin/missions/update/:id] POST
// updates a mission based on the given id
func (c DMissionsController) Update(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var mission models.Mission
	var err error
	if Mission, ok := app.Mapper.GetModel(&mission); ok {
		if err = Mission.Find(id).Exec(&mission); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			if err = c.PermitParams(&mission, "title", "description"); err != nil {
				revel.ERROR.Print("Error Binding Mission")
				return c.ErrorResponse(err, err.Error(), 400)
			}
		} else {

			mission.Type = c.Params.Form.Get("type")
			if c.Params.Get("language") != "" {
				if core.FindOnArray(mission.Langs, c.Params.Get("language")) < 0 {
					mission.Langs = append(mission.Langs, c.Params.Get("language"))
				}
				mission.Description[c.Params.Get("language")] = c.Params.Get("description")
				mission.Title[c.Params.Get("language")] = c.Params.Get("title")
			}

		}

		// If there's picture then update it
		if len(c.Params.Files["cover_picture"]) > 0 {

			var byteArray = c.ResizeImage(400, "cover_picture", core.BlurZero)
			if len(byteArray) <= 0 {
				return c.Redirect(DMissionsController.Index)
			}

			if mission.Attachment.PATH != "" {
				mission.Attachment.Remove()
			}

			// owner := &mgo.DBRef{
			// 	Id:         mission.GetID(),
			// 	Collection: mission.GetDocumentName(),
			// 	Database:   app.Mapper.DatabaseName,
			// }

			if err = mission.Attachment.UploadBytes(models.AsDocumentBase(&mission), byteArray, c.Params.Files["cover_picture"][0].Filename); err != nil {
				revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
				return c.Redirect(DMissionsController.Index)
			}
		}

		if err = Mission.Update(&mission); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse("success", "update success", 0, nil)
		}
		return c.Redirect(DMissionsController.Index)

	}

	return c.ServerErrorResponse()
}

// Delete [/spyc_admin/missions/:id] DELETE
// is a logical deletion of the mission with the given id
func (c DMissionsController) Delete(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var mission models.Mission

	if Mission, ok := app.Mapper.GetModel(&mission); ok {
		if err := Mission.Find(id).Exec(&mission); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		mission.SetStatus(core.StatusInactive)
		mission.Deleted = true

		if err := Mission.Update(&mission); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		core.Notify(core.InfoEntry, "Cash Hunt was Deactivated", "A Cash Hunt was Dectivated", mission.Title.GetString(c.Request.Locale),
			core.GetDashboardPath()+"missions/"+mission.GetID().Hex(), "", mission.Description.GetString(c.Request.Locale), []interface{}{
				bson.M{"title": mission.Type, "value": core.ConcatArray(mission.Countries), "short": true},
			})

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}

// ActivateMission [/spyc_admin/missions/active/:id] POST
// change Mission status to active
func (c DMissionsController) ActivateMission(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var mission models.Mission
	if Mission, ok := app.Mapper.GetModel(&mission); ok {
		if err := Mission.Find(id).Exec(&mission); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		code, err := c.MissionActivation(mission)
		switch code {
		case 200:
			break
		case 400:
			return c.ErrorResponse(err.Error(), err.Error(), code)
		case 500:
			return c.ErrorResponse(err.Error(), err.Error(), code)
		default:

		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}
	return c.ServerErrorResponse()
}

// MissionActivation activates the mission
func (c DMissionsController) MissionActivation(mission models.Mission) (int, error) {

	if Mission, ok := app.Mapper.GetModel(&mission); ok {

		var selector = bson.M{
			"_id":         mission.GetID(),
			"status.name": core.StatusInit,
		}

		var query = bson.M{"$set": bson.M{
			"status.name": core.StatusActive,
			"status.code": core.SubscriptionStatus[core.StatusActive],
		}}

		if err := Mission.UpdateQuery(selector, query, false); err != nil {
			revel.ERROR.Printf("ERROR UPDATE Mission %s --- %s", mission.GetID().Hex(), err.Error())
			return 400, err
		}

		var missionLang = mission.Langs[0]

		var ok bool
		var notificaMessage = c.Message("cashhunt.new", serializers.InternationalizeSerializer{}.Get(mission.Title, missionLang))
		var notificaMessage2 string
		var deviceList []struct {
			ID    string   `bson:"_id"`
			Users []string `bson:"users"`
		}
		var deviceList2 = deviceList

		switch mission.Type {
		case core.TypeCountry:
			deviceList, ok = GetUsersInCountry(mission.Countries).([]struct {
				ID    string   `bson:"_id"`
				Users []string `bson:"users"`
			})
		case core.TypeBroadcast:
			deviceList, ok = GetUsersInCountry(mission.Countries).([]struct {
				ID    string   `bson:"_id"`
				Users []string `bson:"users"`
			})

			deviceList2, ok = GetUsersNotInCountry(mission.Countries).([]struct {
				ID    string   `bson:"_id"`
				Users []string `bson:"users"`
			})
			notificaMessage2 = c.Message("cashhunt.newNotYourCountry", serializers.InternationalizeSerializer{}.Get(mission.Title, missionLang), core.ConcatArray(mission.Countries))
		default:
			deviceList, ok = GetUsersByCoordinates(mission.Geolocation).([]struct {
				ID    string   `bson:"_id"`
				Users []string `bson:"users"`
			})
		}

		if !ok || deviceList == nil || len(deviceList) == 0 {
			return 400, errors.New("No Users found")
		}

		var notification = models.Notification{
			From:        mission.GetDocumentName(),
			To:          mission.GetID(),
			Resource:    mission.GetID().Hex(),
			Action:      "play",
			Type:        core.CashHuntNew,
			Message:     notificaMessage,
			FullMessage: notificaMessage,
			Screen:      "cash_hunt",
			Title:       c.Message("cashhunt.newTitle"),
			Device:      "",
			Attachment:  mission.Attachment,
		}

		// Send normal notificas
		for i := 0; i < len(deviceList); i++ {
			if deviceList[i].ID == "Android" && len(deviceList[i].Users) > 0 {
				SendRegIds(deviceList[i].Users, "Android", &notification)
			}
			if deviceList[i].ID == "IOS" && len(deviceList[i].Users) > 0 {
				SendRegIds(deviceList[i].Users, "IOS", &notification)
			}
		}

		// If Mission type is broadcast
		if mission.Type == core.TypeBroadcast {
			notification.Message = notificaMessage2
			notification.FullMessage = notificaMessage2
			for i := 0; i < len(deviceList); i++ {
				if deviceList[i].ID == "Android" && len(deviceList2[i].Users) > 0 {
					SendRegIds(deviceList2[i].Users, "Android", &notification)
				}
				if deviceList[i].ID == "IOS" && len(deviceList2[i].Users) > 0 {
					SendRegIds(deviceList2[i].Users, "IOS", &notification)
				}
			}
		}

		// Sending to Slack
		core.Notify(core.InfoEntry, "Cash Hunt was Activated", "A Cash Hunt was Activated", mission.Title.GetString(missionLang),
			core.GetDashboardPath()+"missions/"+mission.GetID().Hex(), "", mission.Description.GetString(missionLang), []interface{}{
				bson.M{"title": mission.Type, "value": core.ConcatArray(mission.Countries), "short": true},
			})

		return 0, nil
	}
	return 500, errors.New("ERROR Mapping Mission")
}

// GetMissionsByCampaign [] GET
// Retreives a list of missions by the given campaign
func (c DMissionsController) GetMissionsByCampaign(idCampaign string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	if !bson.IsObjectIdHex(idCampaign) {
		return c.ErrorResponse("Invalid ID", c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var misisons []models.Mission
	if Mission, ok := app.Mapper.GetModel(&models.Mission{}); ok {
		var query = []bson.M{
			bson.M{"status.name": bson.M{"$nin": []string{core.StatusInactive, core.StatusInit}}},
			bson.M{"campaign": bson.ObjectIdHex(idCampaign)},
		}
		if err := Mission.FindWithOperator("$and", query).Exec(&misisons); err != nil {
			revel.ERROR.Print("ERROR FIND Missions --- " + err.Error())
		}
		return c.SuccessResponse(misisons, "Success", core.ModelsType[core.ModelGame], serializers.MissionSerializer{Lang: c.CurrentUser.Device.Language})
	}
	revel.ERROR.Print("ERROR MAPPING Mission")
	return c.ServerErrorResponse()
}

// GetUsersByCoordinates ...
func GetUsersByCoordinates(geo *models.Geo) interface{} {
	if User, ok := app.Mapper.GetModel(&models.User{}); ok {
		var deviceList []struct {
			ID    string   `bson:"_id"`
			Users []string `bson:"users"`
		}

		var near = bson.M{
			"type":        "Point",
			"coordinates": []float64{geo.Coordinates[0], geo.Coordinates[1]},
		}
		var match = bson.M{"$group": bson.M{
			"_id":   "$device.os",
			"users": bson.M{"$push": "$$ROOT.device.messaging_token"},
		}}
		var pipe = mgomap.Aggregate{}.GeoNear(near, "geolocation", core.MaxDistanceForPin, nil, "geolocation", 0).Add(match)

		if err := User.Pipe(pipe, &deviceList); err != nil {
			var coords = strconv.FormatFloat(geo.Coordinates[0], 'E', -1, 64) + "," + strconv.FormatFloat(geo.Coordinates[1], 'E', -1, 64)
			revel.ERROR.Printf("ERROR FIND Users for Coordinates: %s --- %s", coords, err.Error())
			return nil
		}
		return deviceList
	}
	return nil
}

// GetUsersInCountry ...
func GetUsersInCountry(countries []string) interface{} {
	if User, ok := app.Mapper.GetModel(&models.User{}); ok {
		var deviceList []struct {
			ID    string   `bson:"_id"`
			Users []string `bson:"users"`
		}

		var match = bson.M{
			"$match": bson.M{"personal_data.address.country": bson.M{"$in": countries}},
		}
		var group = bson.M{"$group": bson.M{
			"_id":   "$device.os",
			"users": bson.M{"$push": "$device.messaging_token"},
		}}

		var pipe = mgomap.Aggregate{}.Add(match).Add(group)

		if err := User.Pipe(pipe, &deviceList); err != nil {
			revel.ERROR.Printf("ERROR FIND Users for Mission Countries --- %s", err.Error())
			return nil
		}
		return deviceList
	}
	return nil
}

// GetUsersNotInCountry returns a list of users not in the given country list
func GetUsersNotInCountry(countries []string) interface{} {
	if User, ok := app.Mapper.GetModel(&models.User{}); ok {
		var deviceList []struct {
			ID    string   `bson:"_id"`
			Users []string `bson:"users"`
		}

		var match = bson.M{
			"$match": bson.M{"personal_data.address.country": bson.M{"$nin": countries}},
		}
		var group = bson.M{"$group": bson.M{
			"_id":   "$device.os",
			"users": bson.M{"$push": "$device.messaging_token"},
		}}

		var pipe = mgomap.Aggregate{}.Add(match).Add(group)

		if err := User.Pipe(pipe, &deviceList); err != nil {
			revel.ERROR.Printf("ERROR FIND Users for Mission Countries --- %s", err.Error())
			return nil
		}
		return deviceList
	}
	return nil
}

// SendRegIds send the notification
func SendRegIds(tokenList []string, device string, notification *models.Notification) {

	if len(tokenList) > 0 {
		var sliceTokens = make([]string, core.MaxPushRegIds)
		if len(tokenList) > core.MaxPushRegIds {
			sliceTokens = tokenList[0:core.MaxPushRegIds]
		} else {
			sliceTokens = tokenList[0:len(tokenList)]
		}
		if _, ok := notification.Resource.(string); ok {
			notification.SendRegsIds(sliceTokens, device)
			if len(tokenList) > core.MaxPushRegIds {
				SendRegIds(tokenList[(core.MaxPushRegIds+1):len(tokenList)], device, notification)
			}
		}
	}
}
