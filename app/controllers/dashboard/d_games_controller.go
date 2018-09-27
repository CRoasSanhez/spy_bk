package dashboard

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

type DGamesController struct {
	DBaseController
}

// Index [/spyc_admin/webgames] GET
// return all games on DB
func (c DGamesController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var games []models.WebGame

	if WebGame, ok := app.Mapper.GetModel(&models.WebGame{}); ok {
		// err = Mission.Query(bson.M{"title": bson.M{"$regex": search}}).Sort([]string{"start_date", "title"}).Paginate(page-1, quantity).Exec(&missions)
		var match = bson.M{"$and": []bson.M{
			bson.M{"status.name": bson.M{"$ne": "inactive"}},
		}}

		var languages []models.Language
		if Language, ok := app.Mapper.GetModel(&models.Language{}); ok {
			if err := Language.FindWithOperator("$and", []bson.M{bson.M{"status.name": core.StatusActive}}).Exec(&languages); err != nil {
				revel.ERROR.Print("ERROR FIND Languages ---" + err.Error())
			}
		}

		var pipe = mgomap.Aggregate{}.Match(match)

		if err := WebGame.Pipe(pipe, &games); err != nil {
			revel.ERROR.Print(err)
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(games, "success", 200, serializers.WebGameSerializer{})
		}

		c.ViewArgs["Games"] = games
		c.ViewArgs["Languages"] = languages

		return c.Render()
	}

	return c.ServerErrorResponse()
}

// Show [/spyc_admin/webgames/:id] GET
// returns the HTML view of game Detail
func (c DGamesController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var webgame models.WebGame

	if WebGame, ok := app.Mapper.GetModel(&webgame); ok {
		if err := WebGame.Find(id).Exec(&webgame); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			return c.SuccessResponse(webgame, "success", 200, serializers.WebGameSerializer{})
		}

		c.ViewArgs["WebGame"] = webgame

		return c.Render()
	}

	return c.ServerErrorResponse()
}

// Create [/spyc_admin/webgames] POST
// inserts a webgame document in DB
func (c DGamesController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var game models.WebGame
	var err error

	if c.Params.Get("format") == "json" {
		if err = c.Params.BindJSON(&game); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
	} else {
		game.Name = c.Params.Form.Get("name")
		game.NameURL = c.Params.Form.Get("name_url")
		//game.Description = c.Params.Get("description")
		game.Langs = append(game.Langs, c.Params.Get("language"))
		game.Description = map[string]interface{}{c.Params.Get("language"): c.Params.Get("description")}
	}

	if WebGame, ok := app.Mapper.GetModel(&game); ok {
		if err = WebGame.Create(&game); err != nil {
			revel.ERROR.Print("ERROR CREATE Game --- " + err.Error())
			return c.Redirect(DGamesController.Index)
		}

		if len(c.Params.Files["webgame_picture"]) == 1 {

			if game.Attachment.PATH != "" {
				game.Attachment.Remove()
			}

			// var owner = &mgo.DBRef{
			// 	Id:         game.GetID(),
			// 	Collection: game.GetDocumentName(),
			// 	Database:   app.Mapper.DatabaseName,
			// }

			// Upoload multipart file (normal image)
			if err = game.Attachment.Init(models.AsDocumentBase(&game), c.Params.Files["webgame_picture"][0]); err != nil {
				revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
				return c.Redirect(DGamesController.Index)
			}

			if err = game.Attachment.Upload(); err != nil {
				revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
				return c.Redirect(DGamesController.Index)
			}

			// Upload Image(byteArray) thumbnail
			var byteArray = c.ResizeImage(100, "webgame_picture", core.BlurChallenge)
			if len(byteArray) <= 0 {
				return c.Redirect(DGamesController.Index)
			}
			if err = game.Thumbnail.UploadBytes(models.AsDocumentBase(&game), byteArray, c.Params.Files["webgame_picture"][0].Filename); err != nil {
				revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
				return c.Redirect(DGamesController.Index)
			}

			if err = WebGame.Update(&game); err != nil {
				revel.ERROR.Print("ERROR UPDATE Game --- " + err.Error())
				return c.Redirect(DGamesController.Index)
			}
		}
	}
	return c.Redirect(DGamesController.Index)
}

// Update [/spyc_admin/webgames/update/:id] POST
// updates a webgame based on the given id
func (c DGamesController) Update(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var game models.WebGame
	var err error

	if WebGame, ok := app.Mapper.GetModel(&game); ok {
		if err = WebGame.Find(id).Exec(&game); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if c.Params.Get("format") == "json" {
			if err := c.Params.BindJSON(&game); err != nil {
				return c.ErrorResponse(err, err.Error(), 400)
			}
		} else {
			game.Name = c.Params.Form.Get("name")
			game.NameURL = c.Params.Form.Get("name_url")
			//webgame.Description = c.Params.Form.Get("description")
			if c.Params.Get("language") != "" {
				if core.FindOnArray(game.Langs, c.Params.Get("language")) < 0 {
					game.Langs = append(game.Langs, c.Params.Get("language"))
				}
				game.Description[c.Params.Get("language")] = c.Params.Get("description")
			}
		}

		if len(c.Params.Files["webgame_picture"]) > 0 {

			if game.Attachment.PATH != "" {
				game.Attachment.Remove()
			}

			// var owner = &mgo.DBRef{
			// 	Id:         game.GetID(),
			// 	Collection: game.GetDocumentName(),
			// 	Database:   app.Mapper.DatabaseName,
			// }

			// Upoload multipart file (normal image)
			if err = game.Attachment.Init(models.AsDocumentBase(&game), c.Params.Files["webgame_picture"][0]); err != nil {
				revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
				return c.Redirect(DGamesController.Index)
			}

			if err = game.Attachment.Upload(); err != nil {
				revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
				return c.Redirect(DGamesController.Index)
			}

			// Upload Image(byteArray) thumbnail
			var byteArray = c.ResizeImage(100, "webgame_picture", core.BlurChallenge)
			if len(byteArray) <= 0 {
				return c.Redirect(DGamesController.Index)
			}

			if err = game.Thumbnail.UploadBytes(models.AsDocumentBase(&game), byteArray, c.Params.Files["webgame_picture"][0].Filename); err != nil {
				revel.ERROR.Print("ERROR UPLOAD FileBytes --- " + err.Error())
				return c.Redirect(DGamesController.Index)
			}
		}

		if err = WebGame.Update(&game); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
		//return c.SuccessResponse("success", "update success", 0, nil)
	}
	return c.Redirect(DGamesController.Index)
	//return c.ServerErrorResponse()
}

// ActivateGame [/spyc_admin/webgames/activate/:id] POST
// activates webgame from dashboard
func (c DGamesController) ActivateGame(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var game models.WebGame

	if WebGame, ok := app.Mapper.GetModel(&game); ok {
		if err := WebGame.Find(id).Exec(&game); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		game.SetStatus(core.StatusActive)

		if err := WebGame.Update(&game); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}

// Delete [/spyc_admin/webgames/:id] DELETE
// is a logical deletion of the WebGame with the given id
func (c DGamesController) Delete(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.Redirect("/spyc_admin/sessions/new")
	}

	var game models.WebGame

	if WebGame, ok := app.Mapper.GetModel(&game); ok {
		if err := WebGame.Find(id).Exec(&game); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		game.SetStatus(core.StatusInactive)
		game.Deleted = true

		if err := WebGame.Update(&game); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}
