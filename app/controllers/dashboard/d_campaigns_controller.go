package dashboard

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

type DCampaignsController struct {
	DBaseController
}

// Index [/spyc_admin/campaigns] GET
// returns all Sponsonrs on DB
// Not in use
func (c DCampaignsController) Index(page int, quantity int, search string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var campaigns []models.Campaign
	if Campaign, ok := app.Mapper.GetModel(&models.Campaign{}); ok {
		if page == 0 {
			page = 1
		}

		// Get all sponsors
		if err := Campaign.All().Exec(&campaigns); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.Render()
	}

	return c.ServerErrorResponse()
}

// Show [/spyc_admin/campaigns/:id] GET
// returns the object or the View according to request param ("json")
func (c DCampaignsController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var campaign models.Campaign
	if Campaign, ok := app.Mapper.GetModel(&campaign); ok {
		if !bson.IsObjectIdHex(id) {
			return c.ErrorResponse(c.Message("error.invalid"), c.Message("error.invalid"), 400)
		}

		var query = []bson.M{
			bson.M{"status.name": bson.M{"$ne": core.StatusInactive}},
		}

		if err := Campaign.FindWithOperator("$and", query).Exec(&campaign); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(campaign, "success", 0, nil)
		} else {
			return c.Render()
		}
	}
	return c.ServerErrorResponse()
}

// Create [/spyc_admin/campaigns/:sponsorID] POST
// inserts a reward document in DB
func (c DCampaignsController) Create(sponsorID string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var campaign models.Campaign

	if c.Params.Get("format") == "json" {
		if err := c.PermitParams(&campaign, "name", "start_date", "end_date", "sponsor"); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
	} else {
		params := c.Params
		startDate, _ := time.Parse(core.MXTimeFormat, params.Get("start_date"))
		endDate, _ := time.Parse(core.MXTimeFormat, params.Get("end_date"))

		campaign.Name = params.Get("name")
		budget, _ := strconv.Atoi(params.Get("budget"))
		campaign.Budget = budget
		campaign.StartDate = startDate
		campaign.EndDate = endDate
		campaign.Sponsor = bson.ObjectIdHex(sponsorID)
	}

	if Campaign, ok := app.Mapper.GetModel(&campaign); ok {
		if err := Campaign.Create(&campaign); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.Redirect("../sp/%s", sponsorID)
	}

	return c.ServerErrorResponse()
}

// ActivateCampaign [/spyc_admin/campaigns/active/:id] POST
// change Campaign status to active
func (c DCampaignsController) ActivateCampaign(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var campaign models.Campaign
	if Campaign, ok := app.Mapper.GetModel(&campaign); ok {
		if err := Campaign.Find(id).Exec(&campaign); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		campaign.SetStatus(core.StatusActive)

		if err := Campaign.Update(&campaign); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}

// Delete [/spyc_admin/campaigns/:id] DELETE
// deletes logically the Campaign
func (c DCampaignsController) Delete(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var campaign models.Campaign

	if Campaign, ok := app.Mapper.GetModel(&campaign); ok {
		if err := Campaign.Find(id).Exec(&campaign); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		campaign.SetStatus(core.StatusInactive)
		campaign.Deleted = true

		if err := Campaign.Update(&campaign); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}

// GetCampaignsBySponsor [/spyc_admin/campaigns/sponsor/:sponsorID] GET
func (c DCampaignsController) GetCampaignsBySponsor(sponsorID string) revel.Result {

	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	// Validate ID valid
	if !bson.IsObjectIdHex(sponsorID) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var campaigns []models.Campaign
	if Campaign, ok := app.Mapper.GetModel(&models.Campaign{}); ok {
		var query = []bson.M{
			bson.M{"status.name": core.StatusActive},
			bson.M{"sponsor": bson.ObjectIdHex(sponsorID)},
		}
		if err := Campaign.FindWithOperator("$and", query).Exec(&campaigns); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
		return c.SuccessResponse(campaigns, "success", 200, serializers.CampaignSerializer{})
	}
	return c.ServerErrorResponse()
}
