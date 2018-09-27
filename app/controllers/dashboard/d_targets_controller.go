package dashboard

import (
	"bytes"
	"image/png"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strconv"
	"strings"
	"time"

	"github.com/Reti500/mgomap"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// DTargetsController ...
type DTargetsController struct {
	DBaseController
}

// Index [/spyc_admin/targets] GET
// Not in use for the moment
func (c DTargetsController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var targets []models.Target

	if Target, ok := app.Mapper.GetModel(&models.Target{}); ok {
		if err := Target.All().Exec(&targets); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(targets, "success", 200, serializers.TargetSerializer{})
	}

	return c.ServerErrorResponse()
}

// Create [/spyc_admin/targets/:missionID] POST
// a target in DB
func (c DTargetsController) Create(missionID string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var params = c.Params
	var target models.Target
	var err error
	var now = time.Now()

	// Convert latitud and longitud from string to float64
	lat, _ := strconv.ParseFloat(params.Get("lat"), 64)
	lng, _ := strconv.ParseFloat(params.Get("long"), 64)
	order, _ := strconv.Atoi(params.Get("order"))

	// Convert dates from string to time
	startDate, err1 := time.Parse(core.MXTimeFormatTZ, c.Params.Form.Get("start_date")+" "+c.Params.Form.Get("start_time"))
	endDate, err2 := time.Parse(core.MXTimeFormatTZ, c.Params.Form.Get("end_date")+" "+c.Params.Form.Get("end_time"))

	if err1 != nil || err2 != nil || startDate.Before(now) || endDate.Before(startDate) {
		revel.ERROR.Print("Error Invalid dates")
		return c.ErrorResponse("Error Invalid date", "Error Invalid date", 400)
	}

	if len(c.Params.Files["target_image"]) == 0 {
		return c.Redirect("../missions/%s", missionID)
	}

	if c.Params.Get("format") == "json" {
		if err = c.PermitParams(&target, "name", "description", "start_date", "end_date", "type", "order", "next_step_id", "geolocation", "order", "web_url"); err != nil {
			revel.ERROR.Print("Error Binding Target")
			return c.ErrorResponse(err, err.Error(), 400)
		}
	} else {
		// Fill Target struct
		target.Langs = append(target.Langs, c.Params.Get("language"))
		target.Name = map[string]interface{}{c.Params.Get("language"): c.Params.Get("name")}
		target.Description = map[string]interface{}{c.Params.Get("language"): c.Params.Get("description")}
		target.StartDate = core.ChangeUTCTimeToLocalZone(startDate)
		target.EndDate = core.ChangeUTCTimeToLocalZone(endDate)
		target.Type = params.Get("type")
		target.Order = order
		target.WebURL = params.Get("web_url")
		target.Geolocation = &models.Geo{Type: "Point", Coordinates: []float64{lng, lat}}
		target.Missions = bson.ObjectIdHex(missionID)
		target.Score = params.Get("score")
	}

	if Target, ok := app.Mapper.GetModel(&target); ok {
		if err = Target.Create(&target); err != nil {
			revel.ERROR.Printf("ERROR CREATE Target --- %s", err.Error())
			return c.Redirect("../missions/%s", missionID)
		}

		// Iftarget type is webgame type then create the webURL
		if target.Type == core.GameTypeGame {

			webgame := params.Get("webgame")
			urlString := core.GetGameBasePath() + ":webGame/?igm=:missionID&istp=:targetID&actokn=:token&lng=:lang&tp=cshnt"

			urlString = strings.Replace(urlString, ":webGame", webgame, 1)
			urlString = strings.Replace(urlString, ":missionID", missionID, 1)
			urlString = strings.Replace(urlString, ":targetID", target.GetID().Hex(), 1)

			target.WebURL = urlString
		}

		// Open Multipart
		// var byteArray = c.ResizeImage(400, "target_image", core.BlurZero)
		// if len(byteArray) <= 0 {
		// 	return c.Redirect("../missions/%s", missionID)
		// }

		if target.Attachment.PATH != "" {
			target.Attachment.Remove()
		}

		if err = target.Attachment.Init(models.AsDocumentBase(&target), c.Params.Files["target_image"][0]); err != nil {
			return c.Redirect("../missions/%s", missionID)
		}

		if err = target.Attachment.Upload(); err != nil {
			return c.Redirect("../missions/%s", missionID)
		}

		// if err = target.Attachment.UploadBytes(owner, byteArray, c.Params.Files["target_image"][0].Filename); err != nil {
		// 	revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
		// 	return c.Redirect("../missions/%s", missionID)
		// }

		target.SetStatus(core.StatusActive)

		if err = Target.Update(&target); err != nil {
			revel.ERROR.Print("Error updating target: " + target.GetID())
			return c.Redirect("../missions/%s", missionID)
		}

		return c.Redirect("../targets/%s", target.GetID().Hex())
	}

	return c.ServerErrorResponse()
}

// Show [/spyc_admin/targets/:id] GET
// returns the HTML view of Target Detail
func (c DTargetsController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var target models.Target

	if Target, ok := app.Mapper.GetModel(&target); ok {
		if err := Target.Find(id).Exec(&target); err != nil {
			return c.Redirect(DMissionsController.Index)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(target, "success", 200, serializers.TargetSerializer{Lang: c.Request.Locale})
		}

		var q models.Question

		if questions, err := q.FindQuestionsByTarget(target.GetID().Hex()); err != nil {
			revel.ERROR.Print(err)
			return c.Redirect(DMissionsController.Index)
		} else {
			c.ViewArgs["Target"] = target
			c.ViewArgs["Questions"] = questions
		}
		return c.Render()
	}
	return c.Redirect(DMissionsController.Index)
}

// Update [/spyc_admin/targets/update/:id] POST
// a target with the given id
func (c DTargetsController) Update(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var target models.Target
	var err error
	if Target, ok := app.Mapper.GetModel(&target); ok {
		if err = Target.Find(id).Exec(&target); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			if err = c.Params.BindJSON(&target); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}
		} else {
			/*
				lat, _ := strconv.ParseFloat(c.Params.Get("lat"), 64)
				lng, _ := strconv.ParseFloat(c.Params.Get("lng"), 64)
				stime, _ := time.Parse(core.MXTimeFormat, c.Params.Form.Get("start_date")+" "+c.Params.Form.Get("start_time"))
				etime, _ := time.Parse(core.MXTimeFormat, c.Params.Form.Get("end_date")+" "+c.Params.Form.Get("end_time"))
			*/
			//target.Name = c.Params.Form.Get("name")
			target.WebURL = c.Params.Form.Get("web_url")
			//target.Description = c.Params.Form.Get("description")

			if c.Params.Get("language") != "" {
				if core.FindOnArray(target.Langs, c.Params.Get("language")) < 0 {
					target.Langs = append(target.Langs, c.Params.Get("language"))
				}
				target.Description[c.Params.Get("language")] = c.Params.Get("description")
				target.Name[c.Params.Get("language")] = c.Params.Get("name")
			}

		}

		// If there's picture then update it
		if len(c.Params.Files["target_image"]) > 0 {

			var byteArray = c.ResizeImage(400, "target_image", core.BlurZero)
			if len(byteArray) <= 0 {
				return c.Redirect(DTargetsController.Index)
			}

			if target.Attachment.PATH != "" {
				target.Attachment.Remove()
			}

			// owner := &mgo.DBRef{
			// 	Id:         target.GetID(),
			// 	Collection: target.GetDocumentName(),
			// 	Database:   app.Mapper.DatabaseName,
			// }

			if err = target.Attachment.UploadBytes(models.AsDocumentBase(&target), byteArray, c.Params.Files["target_image"][0].Filename); err != nil {
				revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
				return c.Redirect(DTargetsController.Index)
			}
		}

		if err = Target.Update(&target); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {

			return c.SuccessResponse("success", "update success", 0, nil)
		}
		return c.Redirect("../../missions/%s", target.Missions.Hex())
	}

	return c.ServerErrorResponse()
}

// Delete [/spyc_admin/targets/:id] DELETE
// changes target status to inactive and deleted = true
func (c DTargetsController) Delete(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var success = false
	var target models.Target

	if Target, ok := app.Mapper.GetModel(&target); ok {
		if err := Target.Find(id).Exec(&target); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		switch target.Type {
		case core.GameTypeNIP, core.GameTypeText:
			break
		case core.GameTypeOptions:
			if err := DeleteTargetQuestions(target); err != nil {
				revel.ERROR.Print(err)
				return c.ErrorResponse(err, err.Error(), 400)
			}
			success = true
		case core.GameTypeCheckIn:
		}

		if !success {
		}

		// Deleting Target
		target.SetStatus(core.StatusInactive)
		target.Deleted = true

		if err := Target.Update(&target); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}

// AddQuestion [/spyc_admin/questions/:idTarget] POST
// inserts a question into Target type options and updates target
func (c DTargetsController) AddQuestion(idTarget string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var question models.Question

	if err := c.Params.BindJSON(&question); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if Question, ok := app.Mapper.GetModel(&question); ok {
		if err := Question.Create(&question); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		var target models.Target

		if Target, ok := app.Mapper.GetModel(&target); ok {
			if err := Target.Find(idTarget).Exec(&target); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}

			if err := target.AddQuestion(question.GetID().Hex()); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}

			return c.SuccessResponse("OK", "Success", 200, nil)
		}
	}

	return c.ServerErrorResponse()
}

// GenerateQR [/spyc_admin/targets/qr/:id] GET
// returns a PNG image as the QR for the given Target
func (c DTargetsController) GenerateQR(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var target models.Target
	if Target, ok := app.Mapper.GetModel(&target); ok {
		if err := Target.Find(id).Exec(&target); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		// If target type is QR
		if target.Type != core.GameTypeQR {
			return c.ErrorResponse(nil, "", 400)
		}

		strMissionID := target.Missions.Hex()
		strData := strMissionID + "-" + target.GetID().Hex() + "-" + target.Type
		encryptedStr := core.EncryptCipher(strData)
		qrCode, _ := qr.Encode(encryptedStr, qr.M, qr.Auto)

		// Scale the barcode to 200x200 pixels
		qrCode, _ = barcode.Scale(qrCode, 300, 300)

		buffer := new(bytes.Buffer)
		if err := png.Encode(buffer, qrCode); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.FileResponse(buffer.Bytes(), id)
	}

	return c.ServerErrorResponse()
}

// GetTargetsByMission [/spyc_admin/targets/game/:id] GET
// returns a slice of targets by the given Mission ID
func (c DTargetsController) GetTargetsByMission(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var g models.Mission

	targets, err := g.GetAllTargetsByMission(id)
	if err != nil {
		revel.ERROR.Print(err)
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(targets, "success", core.ModelsType[core.ModelStep], serializers.TargetSerializer{Lang: c.Request.Locale})
}

// Validation [/spyc_admin/validation] GET
// shows the VIEW to validate user's text target answers
func (c DTargetsController) Validation(page int, quantity int, targetType string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect(DSessionsController.New)
	}

	if page <= 0 {
		page = 1
	}

	if quantity <= 0 {
		quantity = 10
	}

	if targetType == "" {
		targetType = core.GameTypeOptions
	}

	var campaigns []models.Campaign
	if Campaign, ok := app.Mapper.GetModel(&models.Campaign{}); ok {
		if err := Campaign.FindBy("status.name", core.StatusActive).Exec(&campaigns); err != nil {
			return c.Redirect(DashboardController.Index)
		}
	}

	var countries []models.Country
	if Country, ok := app.Mapper.GetModel(&models.Country{}); ok {
		if err := Country.All().Exec(&countries); err != nil {
			return c.Redirect(DashboardController.Index)
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

		// target types
		targetStatus := map[string]string{
			core.StatusInit:         "Init",
			core.StatusActive:       "Active",
			core.StatusCompleted:    "Completed",
			core.StatusUnsubscribed: "Inactive",
		}

		c.ViewArgs["Types"] = targetTypes
		c.ViewArgs["Status"] = targetStatus
		c.ViewArgs["Campaigns"] = campaigns
		c.ViewArgs["Countries"] = countries

		return c.Render()
	}

	return c.Redirect(DashboardController.Index)
}

// Validate [/spyc_admin/validation/:idMission/:idSubscription] POST
// validates the response of user in Validate view
func (c DTargetsController) Validate(idMission, idSubscription string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var answerReq models.WebGameFinish
	var err error
	var targetSubsciption = models.TargetUser{}
	TargetUser, _ := app.Mapper.GetModel(&targetSubsciption)
	var target = models.Target{}
	Target, _ := app.Mapper.GetModel(&target)

	if err = c.Params.BindJSON(&answerReq); err != nil {
		revel.ERROR.Printf("ERROR BIND Request --- %s", err.Error())
		return c.ErrorResponse(nil, c.Message("error.badRequest"), core.ValidationStatus[core.StatusWrong])
	}

	// Find target Subscription
	if err = TargetUser.Find(idSubscription).Exec(&targetSubsciption); err != nil {
		revel.ERROR.Print("ERROR FIND Subscription ID:" + idSubscription)
		return c.ErrorResponse(err, err.Error(), 400)
	}

	// Find user from target Subscription
	var user models.User
	User, _ := app.Mapper.GetModel(&user)

	if err = User.Find(targetSubsciption.UserID.Hex()).Exec(&user); err != nil {
		revel.ERROR.Print("ERROR FIND user ID:" + targetSubsciption.UserID.Hex())
		return c.ServerErrorResponse()
	}

	c.Request.Locale = user.Device.Language
	mission := models.Mission{}
	Mission, _ := app.Mapper.GetModel(&mission)

	var match = bson.M{"$and": []bson.M{
		bson.M{"status.name": bson.M{"$in": []string{core.StatusActive, core.StatusCompleted}}},
		bson.M{"targets._id": targetSubsciption.TargetID},
	}}
	var pipe = mgomap.Aggregate{}.LookUp("targets", "_id", "missions", "targets").Match(match)

	// Find Mission and its targets
	if err = Mission.Pipe(pipe, &mission); err != nil {
		revel.ERROR.Printf("ERROR PIPE Mission with Target: %s --- %s", targetSubsciption.TargetID.Hex(), err.Error())
		return c.ErrorResponse(nil, c.Message("error.badRequest"), core.ValidationStatus[core.StatusNotFound])
	}

	// Find target in nested targets in Mission by bsonID from target subscription
	if target = mission.FindTargetByID(targetSubsciption.TargetID); target.Type == "" {
		return c.ErrorResponse(nil, c.Message("error.notFound", "Target", targetSubsciption.TargetID.Hex()), core.ValidationStatus[core.StatusNotFound])
	}

	// Updating last result date
	if err = models.StatsSetField(user.GetID(), bson.M{"missions.last_result": targetSubsciption.Score}); err != nil {
		revel.ERROR.Printf("ERROR UPDATE last_access for owner: %s --- %s", user.GetID().Hex(), err.Error())
	}

	// Complete subscription to sucess or wrong
	if answerReq.Data == core.StatusSuccess {
		if err = models.StatsUpdateFields(bson.M{"missions.completed": 1}, []string{"missions.last_completed", "missions.last_access"}, c.CurrentUser.GetID()); err != nil {
		}
		err = targetSubsciption.CompleteSubscription(core.StatusCompleted)
	} else {
		if err = models.StatsUpdateFields(bson.M{"missions.uncompleted": 1}, []string{"missions.last_uncompleted", "missions.last_access"}, c.CurrentUser.GetID()); err != nil {
		}
		err = targetSubsciption.CompleteSubscription(core.StatusCompletedWrong)
	}

	if target.NextTargetID != "" {

		targetResponse := models.Target{}

		// Find next Target
		if err = Target.Find(target.NextTargetID).Exec(&targetResponse); err != nil {
			revel.ERROR.Print("ERROR FIND Target ID:" + target.NextTargetID)
			return c.ServerErrorResponse()
		}

		// Notify to the user that the target was completed and send the next one
		go c.NotifyUser(answerReq.Data, mission, user, targetResponse)

		return c.SuccessResponse(bson.M{"data": "success"}, "Complete next target", core.ModelsType[core.ModelSimpleResponse], nil)
	}

	// Complete Mission subscription
	var missionSubscription models.MissionUser
	if answerReq.Data != core.StatusSuccess {
		if err = missionSubscription.CompleteSubscriptionByMissionUser(mission.GetID().Hex(), targetSubsciption.UserID.Hex(), core.StatusCompletedNotWinner); err != nil {
			return c.ServerErrorResponse()
		}
		c.NotifyUser(answerReq.Data, mission, user, mission)
		return c.SuccessResponse("OK", "Validated Wrong answer", 200, nil)
	}

	if err = missionSubscription.CompleteSubscriptionByMissionUser(mission.GetID().Hex(), targetSubsciption.UserID.Hex(), core.StatusCompletedWinner); err != nil {
		return c.ServerErrorResponse()
	}

	var reward models.Reward
	rewards := reward.GetRewardsByResource(targetSubsciption.TargetID.Hex())
	if len(rewards) == 0 {

		// Notify the user has completed the target but there are no rewards left
		c.NewNotification("play", mission.Attachment, c.Message("reward.notObtained"), core.WinTypeCashHunt, c.Message("reward.notObtained"), core.CashHuntFinished, mission.GetDocumentName(), user, mission.GetID())

		return c.SuccessResponse("OK", "There are no rewards left", core.GameStatus[core.StatusError], nil)
	}

	// Create winner
	reward = rewards[0]
	win := models.Win{}
	_, errWinner := win.CreateWin(reward.GetID(), targetSubsciption.UserID, core.WinTypeCashHunt, targetSubsciption.Score)
	if errWinner != nil {
		revel.ERROR.Printf("ERROR CREATE win for rewardID: %s --- %s", reward.GetID().Hex(), err.Error())
		return c.ServerErrorResponse()
	}

	// if reward is active then change its status to complete
	if ok := reward.CompleteReward(targetSubsciption.UserID.Hex()); !ok {
		revel.ERROR.Printf("ERROR UPDATING reward status: %s --- %s", reward.GetID().Hex(), err.Error())
		return c.ServerErrorResponse()
	}

	if err = TargetUser.UpdateQuery(bson.M{"_id": targetSubsciption.GetID()}, bson.M{"$set": bson.M{"reward_id": reward.GetID()}}, false); err != nil {
		revel.ERROR.Printf("ERROR UPDATE Target Subscription: %s --- %s", targetSubsciption.GetID(), err.Error())
	}

	reward.GetCoinsByMission(idMission)

	if (len(rewards) == 1 && !reward.Multi) || (reward.Multi && reward.Status.Name == core.StatusCompleted) {
		if mission.Complete() {
			c.NotifyMissionComplete(mission)
		}
	}

	// Notify the user has btained the reward
	c.NotifyUser(answerReq.Data, mission, user, reward)

	return c.SuccessResponse(bson.M{"data": "success"}, "success", core.ModelsType[core.ModelSimpleResponse], nil)
}

// GetSubscriptionsByTarget [/spyc_admin/validation/:idTarget] GET
// returns a list of TargetUser by te given id
func (c DTargetsController) GetSubscriptionsByTarget(idTarget string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var subs []models.TargetUser

	if TargetUser, ok := app.Mapper.GetModel(&models.TargetUser{}); ok {
		// subscriptions in pending and active Targets with target Mission equals to the given idMission
		match := bson.M{"$and": []bson.M{
			bson.M{"status.name": core.StatusPendingValidation},
			//bson.M{"targets.missions": bson.ObjectIdHex(idMission)},
			bson.M{"targets._id": bson.ObjectIdHex(idTarget)},
			bson.M{"targets.status.name": core.StatusActive},
		}}

		pipe := mgomap.Aggregate{}.LookUp("targets", "target_id", "_id", "targets").Match(match).Sort(bson.M{"updated_at": 1, "title": 1})

		if err := TargetUser.Pipe(pipe, &subs); err != nil {
			revel.ERROR.Print(err)
			return c.ErrorResponse(err, err.Error(), 400)
		}
		return c.SuccessResponse(subs, "success", core.ModelsType[core.ModelStep], serializers.TargetValidationSerializer{})
	}
	return c.ServerErrorResponse()
}

// NotifyUser notifies the target status to the user and send the given target
func (c DTargetsController) NotifyUser(status string, from models.Mission, toUser models.User, resourceModel interface{}) {

	var msg, resource, screen, answerMsg string
	var url models.Attachment

	if rewardRes, ok := resourceModel.(models.Reward); ok {
		msg = c.Message("reward.obtained")
		resource = rewardRes.GetID().Hex()
		screen = core.ModelReward
		url = rewardRes.Attachment
		answerMsg = c.Message("reward.obtained")

		// Notify to Slack that the user has obtained the reward
		core.Notify(core.NewEntry, "Reward Obtained", "A Reward was Obtained", from.Title.GetString(toUser.Device.Language),
			core.GetDashboardPath()+"users/"+toUser.GetID().Hex(), "", from.Description.GetString(toUser.Device.Language), []interface{}{
				bson.M{"title": from.Type, "value": core.ConcatArray(from.Countries), "short": true},
			})

		goto SendNotification
	}

	if targetRes, ok2 := resourceModel.(models.Target); ok2 {
		msg = c.Message("target.nextTarget")
		resource = targetRes.GetID().Hex()
		screen = "target"
		url = targetRes.Attachment
		answerMsg = c.Message("target.nextTarget")
		goto SendNotification
	}

	if missionRes, ok3 := resourceModel.(models.Mission); ok3 {
		msg = c.Message("answer.wrong")
		resource = missionRes.GetID().Hex()
		screen = core.WinTypeCashHunt
		url = missionRes.Attachment
		answerMsg = c.Message("answer.wrong")
		goto SendNotification
	}

SendNotification:
	c.NewNotification("play", url, msg, screen, answerMsg, core.CashHuntReward, from.GetDocumentName(), toUser, resource)
}

// UpdateTargetsOrder [/spyc_admin/targets/order/:idMission] POST
func (c DTargetsController) UpdateTargetsOrder() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var jsonPrams map[string]interface{}
	var target models.Target

	if err := c.Params.BindJSON(&jsonPrams); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if Target, ok := app.Mapper.GetModel(&target); ok {
		for i, v := range jsonPrams {

			var idString, _ = v.(string)

			if !bson.IsObjectIdHex(idString) {
				revel.ERROR.Print("ERROR invalid ID: " + idString)
				continue
			}

			orderInt, _ := strconv.Atoi(i)
			var orderStr string
			var next_target string

			// Verify if actual loop is not last
			if orderInt != len(jsonPrams) {
				orderStr = strconv.Itoa(orderInt + 1)
				next_target, _ = jsonPrams[orderStr].(string)
			} else {
				next_target = ""
			}

			var selector = bson.M{"_id": bson.ObjectIdHex(idString)}
			var query = bson.M{
				"$set": bson.M{"order": orderInt, "next_target_id": next_target},
			}

			if err := Target.UpdateQuery(selector, query, false); err != nil {
				revel.ERROR.Print("ERROR UPDATING target ID: " + idString)
			}
		}
		return c.SuccessResponse("success", "success", core.ValidationStatus[core.StatusSuccess], nil)
	}

	return c.ServerErrorResponse()
}

// DeleteTargetQuestions change all questions from target to status inactive
func DeleteTargetQuestions(target models.Target) error {
	q := models.Question{}

	// Get target questions
	questions, err := q.FindQuestionsByTarget(target.GetID().Hex())
	if err != nil {
		return err
	}

	// Delete questions
	for _, q := range questions {
		if !q.DeleteQuestion() {
			revel.ERROR.Printf("Question %s not deleted", q.GetID().Hex())
			continue
		}
	}
	return nil
}
