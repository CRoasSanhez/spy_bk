package dashboard

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"

	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

type DAppPathsController struct {
	DBaseController
}

// Index [/spyc_admin/app_paths] GET
// returns all routes in DB
func (c DAppPathsController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var appPaths []models.AppPath

	if AppPath, ok := app.Mapper.GetModel(&models.AppPath{}); ok {

		var query = []bson.M{
			{"status.name": bson.M{"$ne": "inactive"}},
		}

		if err := AppPath.FindWithOperator("$and", query).Exec(&appPaths); err != nil {
			revel.ERROR.Print(err)
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(appPaths, "success", 200, nil)
		}

		c.ViewArgs["AppPaths"] = appPaths
		c.ViewArgs["token"] = c.CurrentUser.Saas.Token

		return c.Render()

	}

	return c.ServerErrorResponse()
}

// Show [/spyc_admin/app_paths/:id] GET
// Returns the detail VIew
func (c DAppPathsController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var appPath models.AppPath

	if AppPath, ok := app.Mapper.GetModel(&models.AppPath{}); ok {

		var query = []bson.M{
			{"_id": bson.ObjectIdHex(id)},
		}

		if err := AppPath.FindWithOperator("$and", query).Exec(&appPath); err != nil {
			revel.ERROR.Print(err)
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(appPath, "success", 200, nil)
		}

		c.ViewArgs["AppPath"] = appPath

		return c.Render()

	}

	return c.ServerErrorResponse()
}

// Create [/spyc_admin/app_paths] POST
// Creates a route document in DB
func (c DAppPathsController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var appPath models.AppPath
	var err error

	if c.Params.Get("format") == "json" {
		if err := c.PermitParams(&appPath, "name", "section", "url", "enviroment"); err != nil {
			revel.ERROR.Print("Error Binding Mission")
			return c.ErrorResponse(err, err.Error(), 400)
		}
	} else {
		appPath.Name = c.Params.Get("name")
		appPath.Section = c.Params.Get("section")
		appPath.URL = c.Params.Get("url")
		appPath.Enviroment = c.Params.Get("enviroment")
	}

	if AppPath, ok := app.Mapper.GetModel(&appPath); ok {
		if err = AppPath.Create(&appPath); err != nil {
			revel.ERROR.Print("ERROR CREATING AppPath --- " + err.Error())
			return c.Redirect(DAppPathsController.Index)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(appPath, "success", 200, nil)
		}

		return c.Redirect(DAppPathsController.Index)
	}

	return c.Redirect(DAppPathsController.Index)
}

// Update [/spyc_admin/app_paths/:id] PATCH
// updates a mission based on the given id
func (c DAppPathsController) Update(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var appPath models.AppPath
	if AppPath, ok := app.Mapper.GetModel(&appPath); ok {
		if err := AppPath.Find(id).Exec(&appPath); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			if err := c.PermitParams(&appPath, "name", "section", "url", "enviroment"); err != nil {
				revel.ERROR.Print("Error Binding AppPath")
				return c.ErrorResponse(err, err.Error(), 400)
			}
		} else {

			appPath.Name = c.Params.Form.Get("name")
			appPath.Section = c.Params.Form.Get("section")
			appPath.URL = c.Params.Form.Get("url")
			appPath.Enviroment = c.Params.Form.Get("enviroment")
		}

		if err := AppPath.Update(&appPath); err != nil {

			if c.Params.Get("format") == "json" {
				return c.ErrorResponse(err, err.Error(), 400)
			}
			return c.ServerErrorResponse()
		}
		return c.SuccessResponse("success", "update success", 0, nil)
	}
	return c.ServerErrorResponse()
}

// Delete [/spyc_admin/app_paths] DELETE
// deletes logically the route by the given id
func (c DAppPathsController) Delete(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var appPath models.AppPath

	if AppPath, ok := app.Mapper.GetModel(&appPath); ok {
		if err := AppPath.Find(id).Exec(&appPath); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		appPath.SetStatus(core.StatusInactive)
		appPath.Deleted = true

		if err := AppPath.Update(&appPath); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}
