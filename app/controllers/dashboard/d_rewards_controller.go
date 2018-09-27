package dashboard

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strconv"

	"github.com/Reti500/mgomap"

	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

// DRewardsController ...
type DRewardsController struct {
	DBaseController
}

// Index [/spyc_admin/rewards] GET
// returns all Rewards on DB
func (c DRewardsController) Index(page int, quantity int, idCampaign string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var rewards []interface{}
	if Reward, ok := app.Mapper.GetModel(&models.Reward{}); ok {
		if page == 0 {
			page = 1
		}

		var match = []bson.M{
			bson.M{"status.name": bson.M{"$in": []string{core.StatusActive, core.StatusInit}}},
			bson.M{"resource_type": bson.M{"$nin": []string{core.ModelTypeChallenge}}},
		}

		if idCampaign != "" && bson.IsObjectIdHex(idCampaign) {
			match = append(match, bson.M{"campaign_id": bson.ObjectIdHex(idCampaign)})
		}

		var project = bson.M{"$project": bson.M{"id": "$_id", "name": "$name", "description": "$description", "status": "$status.name",
			"mission": "$missions.title", "target": "$targets.description", "langs": "$langs"}}

		var pipe = mgomap.Aggregate{}.Match(bson.M{"$and": match}).
			LookUp("targets", "resource_id", "_id", "targets").
			Add(bson.M{"$unwind": "$targets"}).
			LookUp("missions", "targets.missions", "_id", "missions").
			Add(bson.M{"$unwind": "$missions"}).Add(project)

		if err := Reward.Pipe(pipe, &rewards); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(rewards, "Rewards", core.ModelsType[core.ModelReward], nil)
		}

		var languages []models.Language
		if Language, ok := app.Mapper.GetModel(&models.Language{}); ok {
			if err := Language.FindWithOperator("$and", []bson.M{bson.M{"status.name": core.StatusActive}}).Exec(&languages); err != nil {
				revel.ERROR.Print("ERROR FIND Languages ---" + err.Error())
			}
		}

		var campaigns []models.Campaign
		if Campaign, ok := app.Mapper.GetModel(&models.Campaign{}); ok {
			if err := Campaign.FindBy("status.name", core.StatusActive).Exec(&campaigns); err != nil {
				return c.Redirect(DashboardController.Index)
			}
		}

		c.ViewArgs["JSapikey"] = core.GoogleMapsJSAPIKey
		c.ViewArgs["Campaigns"] = campaigns
		c.ViewArgs["Languages"] = languages

		return c.Render()
	}

	return c.ServerErrorResponse()
}

// Show [/spyc_admin/rewards/:id] GET
// returns the HTML view of Reward Detail
func (c DRewardsController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var reward models.Reward

	if Reward, ok := app.Mapper.GetModel(&reward); ok {
		if err := Reward.Find(id).Exec(&reward); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(reward, "success", 200, serializers.RewardSerializer{})
		}

		c.ViewArgs["Reward"] = reward

		return c.Render()
	}

	return c.ServerErrorResponse()
}

// Create [/spyc_admin/rewards] POST
// inserts a reward document in DB
func (c DRewardsController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var reward models.Reward
	var err error

	if c.Params.Get("format") == "json" {
		if err := c.Params.BindJSON(&reward); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
	} else {
		resourceID := c.Params.Form.Get("resource_id")
		reward.ResourceID = bson.ObjectIdHex(resourceID)
		reward.ResourceType = c.Params.Form.Get("resource_type")
		reward.Name = map[string]interface{}{c.Params.Get("language"): c.Params.Get("name")}
		reward.Description = map[string]interface{}{c.Params.Get("language"): c.Params.Get("description")}
		reward.Langs = append(reward.Langs, c.Params.Get("language"))
		reward.Multi, _ = strconv.ParseBool(c.Params.Get("multiple"))
		reward.CampaignID = bson.ObjectIdHex(c.Params.Get("campaign"))
		if reward.Multi {
			reward.MaxWinners, _ = strconv.Atoi(c.Params.Get("max_winners"))
		} else {
			reward.MaxWinners = 1
		}
	}
	reward.IsVisible = true

	if Reward, ok := app.Mapper.GetModel(&reward); ok {
		if err := Reward.Create(&reward); err != nil {
			return c.Redirect(DRewardsController.Index)
		}

		if reward.Attachment.PATH != "" {
			reward.Attachment.Remove()
		}

		// Upoload multipart file (normal image)
		if err = reward.Attachment.Init(models.AsDocumentBase(&reward), c.Params.Files["reward_picture"][0]); err != nil {
			revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
			return c.Redirect(DRewardsController.Index)
		}

		if err = reward.Attachment.Upload(); err != nil {
			revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
			return c.Redirect(DRewardsController.Index)
		}

		if err = Reward.Update(&reward); err != nil {
			return c.Redirect(DRewardsController.Index)
		}

		return c.Redirect(DRewardsController.Index)
	}

	return c.ServerErrorResponse()
}

// Update [/spyc_admin/rewards/update/:id] POST
// updates a reward based on the given id
func (c DRewardsController) Update(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var reward models.Reward
	var err error
	if Reward, ok := app.Mapper.GetModel(&reward); ok {
		if err = Reward.Find(id).Exec(&reward); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			if err = c.Params.BindJSON(&reward); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}
		} else {
			if c.Params.Get("language") != "" {
				if core.FindOnArray(reward.Langs, c.Params.Get("language")) < 0 {
					reward.Langs = append(reward.Langs, c.Params.Get("language"))
				}
				reward.Name[c.Params.Get("language")] = c.Params.Get("name")
				reward.Description[c.Params.Get("language")] = c.Params.Get("description")
			}

			// If there's picture then update it
			if len(c.Params.Files["reward_picture"]) > 0 {

				if reward.Attachment.PATH != "" {
					reward.Attachment.Remove()
				}

				if err = reward.Attachment.Init(models.AsDocumentBase(&reward), c.Params.Files["reward_picture"][0]); err != nil {
					return c.Redirect(DRewardsController.Index)
				}

				if err = reward.Attachment.Upload(); err != nil {
					return c.Redirect(DRewardsController.Index)
				}
			}
		}

		if err := Reward.Update(&reward); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse("success", "update success", 0, nil)
		}
		return c.Redirect(DRewardsController.Index)
	}

	return c.ServerErrorResponse()
}

// Delete [/spyc_admin/rewards/:id] DELETE
// is a logical deletion of the Reward with the given id
func (c DRewardsController) Delete(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var reward models.Reward

	if Reward, ok := app.Mapper.GetModel(&reward); ok {
		if err := Reward.Find(id).Exec(&reward); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		reward.SetStatus(core.StatusInactive)
		reward.Deleted = true

		if err := Reward.Update(&reward); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}

// ActivateReward [/spyc_admin/rewards/active/:id] POST
// change Reward status to active
func (c DRewardsController) ActivateReward(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var reward models.Reward

	if Reward, ok := app.Mapper.GetModel(&reward); ok {
		if err := Reward.Find(id).Exec(&reward); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		reward.SetStatus(core.StatusActive)

		if err := Reward.Update(&reward); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}
