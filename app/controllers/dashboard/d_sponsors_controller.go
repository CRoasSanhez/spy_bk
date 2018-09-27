package dashboard

import (
	"math"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"

	"github.com/revel/revel"
)

// DSponsorsController ...
type DSponsorsController struct {
	DBaseController
}

// Index [/spyc_admin/sp] GET
// returns all Sponsonrs on DB
func (c DSponsorsController) Index(page int, quantity int, search string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var sponsors []models.Sponsor

	if Sponsor, ok := app.Mapper.GetModel(&models.Sponsor{}); ok {
		if page == 0 {
			page = 1
		}

		// Get all sponsors
		if err := Sponsor.All().Exec(&sponsors); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		countCountries := len(sponsors)

		c.ViewArgs["Sponsors"] = sponsors

		if quantity > 0 {
			c.ViewArgs["Pages"] = int(math.Ceil(float64(countCountries) / float64(quantity)))
		} else {
			c.ViewArgs["Pages"] = 0
		}

		return c.Render()
	}

	return c.ServerErrorResponse()
}

// Show [/spyc_admin/sp/:id] GET
// returns sponsor name and a list of Campaigns
func (c DSponsorsController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var sponsor models.Sponsor

	if Sponsor, ok := app.Mapper.GetModel(&sponsor); ok {
		if err := Sponsor.Find(id).Exec(&sponsor); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(sponsor, "success", 200, nil)
		}

		// avilable campaigns per Sponsor
		var campaigns []models.Campaign

		if c, err := sponsor.GetAllCampaignsBySponsor(); err != nil {
			revel.ERROR.Print(err)
		} else {
			campaigns = c
		}

		/*
			err = Sponsor.Query(
				[{$lookup:
					{from:"campaigns",
						localField:"_id",
						foreignField:"sponsor",as:"campaigns"
					}
				},
				{ $match : {
					_id : ObjectId("598a606de7f4f82be60497dc")
				}
			}]).Exec(&sponsor)
		*/
		// Get all Countries
		var countries []models.Country

		if Country, ok := app.Mapper.GetModel(&models.Country{}); ok {
			if err := Country.All().Exec(&countries); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}

			c.ViewArgs["Sponsor"] = sponsor
			c.ViewArgs["Countries"] = countries
			c.ViewArgs["Campaigns"] = campaigns

			return c.Render()
		}
	}

	return c.ServerErrorResponse()
}

// Create [/spyc_admin/sp] POST
// inserts a reward document in DB
func (c DSponsorsController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var sponsor models.Sponsor
	var err error

	if Sponsor, ok := app.Mapper.GetModel(&sponsor); ok {

		if c.Params.Get("format") == "json" {
			if err := c.PermitParams(&sponsor, "name"); err != nil {
				revel.ERROR.Print("Error Binding Sponsor")
				return c.ErrorResponse(err, err.Error(), 400)
			}
		} else {
			sponsor.Name = c.Params.Form.Get("name")
		}

		if err = Sponsor.Create(&sponsor); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if sponsor.Attachment.PATH != "" {
			sponsor.Attachment.Remove()
		}

		// owner := &mgo.DBRef{
		// 	Id:         sponsor.GetID(),
		// 	Collection: sponsor.GetDocumentName(),
		// 	Database:   app.Mapper.DatabaseName,
		// }

		if err = sponsor.Attachment.Init(models.AsDocumentBase(&sponsor), c.Params.Files["sponsor_picture"][0]); err != nil {
			return c.ErrorResponse(err, err.Error(), 402)
		}

		if err = sponsor.Attachment.Upload(); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if err = Sponsor.Update(&sponsor); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.Redirect(DSponsorsController.Index)
	}

	return c.Redirect(DSponsorsController.Index)
}

// Update [/spyc_admin/sp/:id] PATCH
// updates a sponsor based on the given id
func (c DSponsorsController) Update(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var sponsor models.Sponsor

	if Sponsor, ok := app.Mapper.GetModel(&sponsor); ok {
		if err := Sponsor.Find(id).Exec(&sponsor); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			if err := c.Params.BindJSON(&sponsor); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}
		} else {
			sponsor.Name = c.Params.Form.Get("name")
		}

		// Validations
		if errs := sponsor.Validate(); errs != nil {
			revel.ERROR.Print(errs)
			return c.ErrorResponse(errs, errs.Error(), 400)
		}

		if err := Sponsor.Update(&sponsor); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("success", "update success", 0, nil)
	}

	return c.ServerErrorResponse()
}

// ActivateSponsor [/spyc_admin/sp/active/:id] POST
// changes Game status to active
func (c DSponsorsController) ActivateSponsor(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var sponsor models.Sponsor

	if Sponsor, ok := app.Mapper.GetModel(&sponsor); ok {
		if err := Sponsor.Find(id).Exec(&sponsor); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		sponsor.SetStatus(core.StatusActive)

		if err := Sponsor.Update(&sponsor); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}

// Delete [/spyc_admin/sp/:id] DELETE
// deletes logically the sponsor
func (c DSponsorsController) Delete(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var sponsor models.Sponsor

	if Sponsor, ok := app.Mapper.GetModel(&sponsor); ok {
		if err := Sponsor.Find(id).Exec(&sponsor); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		sponsor.SetStatus(core.StatusInactive)
		sponsor.Deleted = true

		if err := Sponsor.Update(&sponsor); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}
