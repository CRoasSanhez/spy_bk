package dashboard

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

type DCashHuntController struct {
	DBaseController
}

// Index [/spyc_admin/cashhunt] GET
func (c DCashHuntController) Index() revel.Result {

	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var err error
	var countries []models.Country
	if Country, ok := app.Mapper.GetModel(&models.Country{}); ok {

		pipe := mgomap.Aggregate{}.Sort(bson.M{"name": 1})
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

			var gamesTypes []models.WebGame
			if WebGames, ok := app.Mapper.GetModel(&models.WebGame{}); ok {
				if err = WebGames.FindBy("status.name", core.StatusActive).Exec(&gamesTypes); err != nil {
					return c.ErrorResponse(err, err.Error(), 400)
				}
			}

			var sponsors []models.Sponsor
			if Sponsor, ok := app.Mapper.GetModel(&models.Sponsor{}); ok {
				if err = Sponsor.FindBy("status.name", core.StatusActive).Exec(&sponsors); err != nil {
					return c.ErrorResponse(err, err.Error(), 400)
				}
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

			// target types
			var targetTypes = map[string]string{
				core.GameTypeCheckIn: "CheckIn",
				core.GameTypeGame:    "Game",
				core.GameTypeOptions: "Quiz",
				core.GameTypePhoto:   "Photo",
				core.GameTypeQR:      "QR image",
				core.GameTypeText:    "Text",
			}

			var languages []models.Language
			if Language, ok := app.Mapper.GetModel(&models.Language{}); ok {
				if err = Language.FindWithOperator("$and", []bson.M{bson.M{"status.name": core.StatusActive}}).Exec(&languages); err != nil {
					revel.ERROR.Print("ERROR FIND Languages ---" + err.Error())
				}
			}

			c.ViewArgs["JSapikey"] = core.GoogleMapsJSAPIKey
			c.ViewArgs["Countries"] = countries
			c.ViewArgs["Frequency"] = frequency
			c.ViewArgs["Priority"] = priority
			c.ViewArgs["MissionTypes"] = core.MissionTypes
			c.ViewArgs["Languages"] = languages
			c.ViewArgs["Types"] = targetTypes
			c.ViewArgs["WebGames"] = gamesTypes
			c.ViewArgs["Sponsors"] = sponsors
			return c.Render()

		}
	}

	return c.ServerErrorResponse()
}

// Create [/spyc_admin/cashhunt] POST
// creates a simple and basic cashhunt flow based on the given json
func (c DCashHuntController) Create() revel.Result {

	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	// Parse json request
	var request map[string]interface{}
	//var CashHuntRequest models.CashHuntRequest

	// TODO FINISH THIS METHOS WITH CashHuntRequest and in dashboard

	if err := c.PermitParams(&request, "title", "description", "langs", "type", "geolocation", "start_date", "end_date",
		"countries", "frequency", "instructions", "target_type", "reward_title", "reward_description", "score", "game",
		"web_url", "sponsor", "campaign", "max_winners", "reward_multiple"); err != nil {
		revel.ERROR.Print("Error Binding Mission")
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if !bson.IsObjectIdHex(request["campaign"].(string)) || !bson.IsObjectIdHex(request["sponsor"].(string)) {
		return c.ErrorResponse("Error Invalid Campaign or sponsor", "Error Invalid Campaign or sponsor", 400)
	}

	var now = time.Now()
	stime, err1 := time.Parse(core.MXTimeFormatTZ, request["start_date"].(string))
	etime, err2 := time.Parse(core.MXTimeFormatTZ, request["end_date"].(string))

	stime = core.ChangeUTCTimeToLocalZone(stime)
	etime = core.ChangeUTCTimeToLocalZone(etime)

	// Validate if error parsing dates or startTime is before Now or endTime is before Now or endTime before startTime
	if err1 != nil || err2 != nil || stime.Before(now) || etime.Before(stime) {
		revel.ERROR.Print("Error Invalid dates")
		return c.ErrorResponse("Error Invalid date", "Error Invalid date", 400)
	}

	var strArray, strArray2 []string
	for _, v := range request["countries"].([]interface{}) {
		strArray = append(strArray, v.(string))
	}

	for _, v := range request["langs"].([]interface{}) {
		strArray2 = append(strArray2, v.(string))
	}

	var myInt = request["frequency"].(float64)
	var geoMap = request["geolocation"].(map[string]interface{})

	var coordsArray []float64
	for _, v := range geoMap["coordinates"].([]interface{}) {
		coordsArray = append(coordsArray, v.(float64))
	}

	var geoModel = &models.Geo{
		Type:        geoMap["type"].(string),
		Coordinates: coordsArray,
	}

	// Fill mission struct
	var mission = models.Mission{
		Title:       request["title"].(map[string]interface{}),
		Description: request["description"].(map[string]interface{}),
		Type:        request["type"].(string),
		Geolocation: geoModel,
		StartDate:   stime,
		EndDate:     etime,
		Countries:   strArray,
		Campaign:    bson.ObjectIdHex(request["campaign"].(string)),
		Frequency:   int(myInt),
		Langs:       strArray2,
	}

	// Save mission document
	if Mission, ok := app.Mapper.GetModel(&mission); ok {
		if err := Mission.Create(&mission); err != nil {
			revel.ERROR.Print("Error Creating Mission")
			return c.ErrorResponse(err, err.Error(), 400)
		}
	}

	// Fill target struct
	var target = models.Target{
		Name:        mission.Title,
		Description: request["instructions"].(map[string]interface{}),
		Geolocation: mission.Geolocation,
		Missions:    mission.GetID(),
		Order:       1,
		Score:       request["score"].(string),
		StartDate:   mission.StartDate,
		EndDate:     mission.EndDate,
		Type:        request["target_type"].(string),
		Langs:       mission.Langs,
		WebURL:      request["web_url"].(string),
	}

	// Save target document
	if Target, ok := app.Mapper.GetModel(&target); ok {
		if err := Target.Create(&target); err != nil {
			revel.ERROR.Print("Error Creating Target")
			return c.ErrorResponse(err, err.Error(), 400)
		}

		target.SetStatus(core.StatusActive)

		if target.Type == core.GameTypeGame {

			webgame := request["game"].(string)
			urlString := core.GetGameBasePath() + ":webGame/?igm=:missionID&istp=:targetID&actokn=:token&lng=:lang&tp=cshnt"

			urlString = strings.Replace(urlString, ":webGame", webgame, 1)
			urlString = strings.Replace(urlString, ":missionID", mission.GetID().Hex(), 1)
			urlString = strings.Replace(urlString, ":targetID", target.GetID().Hex(), 1)

			target.WebURL = urlString

		}
		if err := Target.Update(&target); err != nil {
			revel.ERROR.Printf("Error Updating Target: %s --- %s", target.GetID().Hex(), err.Error())
			return c.ErrorResponse(err, err.Error(), 400)
		}
	}

	var maxWinners, _ = strconv.Atoi(request["max_winners"].(string))
	var multi bool
	if multiple, err := strconv.ParseBool(request["reward_multiple"].(string)); err == nil {
		multi = multiple
	}

	// Fill reward struct
	var reward = models.Reward{
		Name:         request["reward_title"].(map[string]interface{}),
		Description:  request["reward_description"].(map[string]interface{}),
		IsVisible:    true,
		ResourceID:   target.GetID(),
		CampaignID:   mission.Campaign,
		ResourceType: target.Type,
		MaxWinners:   maxWinners,
		Multi:        multi,
		Langs:        target.Langs,
	}

	// Save reward document
	if Reward, ok := app.Mapper.GetModel(&reward); ok {
		if err := Reward.Create(&reward); err != nil {
			revel.ERROR.Print("Error Creating Reward")
			return c.ErrorResponse(err, err.Error(), 400)
		}
	}

	return c.SuccessResponse(bson.M{"mid": mission.GetID().Hex(), "tid": target.GetID().Hex()}, "update success", 0, nil)
}

// Picture [/spyc_admin/cashhunt/:idMission/:idTarget/picture] POST
// sets the cashhunt(mission and target) and reward picture
func (c DCashHuntController) Picture(idMission, idTarget string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	if len(c.Params.Files) <= 1 {
		revel.ERROR.Print("ERROR PARAMS Missing File")
		return c.ErrorResponse(nil, c.Message("error.badRequest"), 400)
	}

	var err error
	var mission models.Mission
	var target models.Target
	if Mission, ok := app.Mapper.GetModel(&mission); ok {

		Target, okT := app.Mapper.GetModel(&target)
		if !okT {
			return c.ServerErrorResponse()
		}

		Mission.Find(idMission).Exec(&mission)

		if err := Target.Find(idTarget).Exec(&target); err != nil {
			revel.ERROR.Print("ERROR FIND Target with missions: ---" + err.Error())
			return c.ErrorResponse(nil, c.Message("error.badRequest"), 400)
		}

		// ---- MISION ----
		// Resize image for Mission and target
		//var byteArray = c.ResizeImage(400, "cover_picture", core.BlurZero)
		//if len(byteArray) <= 0 {
		//	return c.ErrorResponse(nil, c.Message("error.badRequest"), 400)
		//}

		// Prepare attachment field
		if mission.Attachment.PATH != "" {
			mission.Attachment.Remove()
		}

		//if err = mission.Attachment.Init(models.AsDocumentBase(&mission), c.Params.Files["cover_picture"][0]); err != nil {
		if err = mission.Attachment.Init(models.AsDocumentBase(&models.User{}), c.Params.Files["cover_picture"][0]); err != nil {
			return c.Redirect(DMissionsController.Index)
		}

		if err = mission.Attachment.Upload(); err != nil {
			return c.Redirect(DMissionsController.Index)
		}

		// Upload MIssion file
		//if err = mission.Attachment.UploadBytes(owner, byteArray, c.Params.Files["cover_picture"][0].Filename); err != nil {
		//	revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
		//	return c.ServerErrorResponse()
		//}

		// Update mission
		if err = Mission.Update(&mission); err != nil {
			revel.ERROR.Printf("ERROR UPDATE Mission Image: %s --- %s", mission.GetID().Hex(), err.Error())
			return c.ServerErrorResponse()
		}

		// ---- TARGET ------
		// Prepare attachment field
		if target.Attachment.PATH != "" {
			target.Attachment.Remove()
		}

		// Upload Target file
		//if err = target.Attachment.UploadBytes(owner, byteArray, c.Params.Files["cover_picture"][0].Filename); err != nil {
		//	revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
		//	return c.ServerErrorResponse()
		//}

		if err = target.Attachment.Init(models.AsDocumentBase(&target), c.Params.Files["cover_picture"][0]); err != nil {
			return c.Redirect(DMissionsController.Index)
		}

		if err = target.Attachment.Upload(); err != nil {
			return c.Redirect(DMissionsController.Index)
		}

		// Update Target
		if err = Target.Update(&target); err != nil {
			revel.ERROR.Printf("ERROR UPDATE Target Image: %s --- %s", target.GetID().Hex(), err.Error())
			return c.ServerErrorResponse()
		}

		// ------ Reward ------
		var reward models.Reward
		if Reward, ok2 := app.Mapper.GetModel(&reward); ok2 {

			Reward.Query(bson.M{"resource_id": bson.ObjectIdHex(idTarget)}).Exec(&reward)

			if reward.Attachment.PATH != "" {
				reward.Attachment.Remove()
			}

			if err = reward.Attachment.Init(models.AsDocumentBase(&reward), c.Params.Files["reward_picture"][0]); err != nil {
				return c.Redirect(DMissionsController.Index)
			}

			if err = reward.Attachment.Upload(); err != nil {
				return c.Redirect(DMissionsController.Index)
			}

			reward.SetStatus(core.StatusActive)

			if err = Reward.Update(&reward); err != nil {
				revel.ERROR.Printf("ERROR UPDATE Reward Image: %s --- %s", reward.GetID().Hex(), err.Error())
				return c.ServerErrorResponse()
			}

			if c.Params.Get("format") != "json" {
				return c.Redirect("/spyc_admin/missions")
			}
			return c.SuccessResponse("success uploads", "update success", 0, nil)
		}
	}

	return c.ServerErrorResponse()
}
