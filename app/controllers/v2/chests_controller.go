package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"gopkg.in/mgo.v2/bson"

	"spyc_backend/app/auth"
	"strings"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
)

// ChestsController ...
type ChestsController struct {
	BaseController
}

// Show ...
func (c ChestsController) Show(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var query = []bson.M{
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"owner.$id": c.CurrentUser.GetID()},
		bson.M{"secret_token": c.Request.Header.Get(core.ChestAccess)},
	}

	var chest models.Chest
	if Chest, ok := app.Mapper.GetModel(&chest); ok {
		if err := Chest.FindWithOperator("$and", query).Exec(&chest); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
	}

	return c.SuccessResponse(chest, "Chest!!", 200, serializers.ChestSerializer{})
}

// Search ...
func (c ChestsController) Search(tags string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if Chest, ok := app.Mapper.GetModel(&models.Chest{}); ok {
		var chests []models.Chest
		var query = []bson.M{
			bson.M{"tags": bson.M{"$in": strings.Split(tags, ",")}},
			bson.M{"owner.$id": c.CurrentUser.GetID()},
		}

		if err := Chest.FindWithOperator("$and", query).Exec(&chests); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(chests, "Success!!", 200, serializers.ChestSerializer{})
	}

	return c.ServerErrorResponse()
}

// Create ...
func (c ChestsController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var chest models.Chest

	if err := c.PermitParams(&chest, true, "name", "pin", "tags", "geolocation"); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	chest.GeneratePIN()
	chest.Size = 1 * core.GB
	chest.Current = 0

	chest.Owner = &mgo.DBRef{
		Id:         c.CurrentUser.GetID(),
		Collection: c.CurrentUser.GetDocumentName(),
		Database:   app.Mapper.DatabaseName,
	}

	if Chest, ok := app.Mapper.GetModel(&chest); ok {
		if err := Chest.Create(&chest); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(chest, "Created!!", 200, serializers.ChestSerializer{})
	}

	return c.ErrorResponse(nil, "Server error", 500)
}

// Update ...
func (c ChestsController) Update(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if len(c.Params.Get("name")) <= 0 {
		return c.ErrorResponse(nil, "Missing paramas", 400)
	}

	var chest models.Chest
	var query = []bson.M{bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"owner": c.CurrentUser.GetID().Hex()}}

	if Chest, ok := app.Mapper.GetModel(&chest); ok {
		if err := Chest.FindWithOperator("$and", query).Exec(&chest); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		chest.Name = c.Params.Get("name")

		if err := Chest.Update(&chest); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
	}

	return c.SuccessResponse(chest, "Updated!!", 200, serializers.ChestSerializer{})
}

// AddFile ...
func (c ChestsController) AddFile(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if len(c.Params.Files["file"]) <= 0 {
		return c.ErrorResponse(nil, "File not found", 400)
	}

	var chest models.Chest

	if Chest, ok := app.Mapper.GetModel(&chest); ok {
		var query = []bson.M{
			bson.M{"_id": bson.ObjectIdHex(id)},
			bson.M{"owner.$id": c.CurrentUser.GetID()},
			bson.M{"secret_token": c.Request.Header.Get(core.ChestAccess)},
		}

		if err := Chest.FindWithOperator("$and", query).Exec(&chest); err != nil {
			return c.ErrorResponse(err, err.Error(), 401)
		}

		var attachment = models.Attachment{}

		if err := attachment.Init(models.AsDocumentBase(&c.CurrentUser), c.Params.Files["file"][0]); err != nil {
			return c.ErrorResponse(err, err.Error(), 402)
		}

		if (chest.Current + attachment.Size) > chest.Size {
			return c.ErrorResponse(nil, "The chest does not have enough storage space", 403)
		}

		attachment.Upload()

		chest.Attachments = append(chest.Attachments, attachment)
		chest.Current += attachment.Size

		if err := Chest.Update(&chest); err != nil {
			return c.ErrorResponse(err, err.Error(), 404)
		}

		return c.SuccessResponse("OK", "Success", 200, nil)
	}

	return c.ServerErrorResponse()
}

// RemoveFile ...
func (c ChestsController) RemoveFile(id, file string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var chest models.Chest
	if Chest, ok := app.Mapper.GetModel(&chest); ok {
		var query = []bson.M{
			bson.M{"_id": bson.ObjectIdHex(id)},
			bson.M{"owner.$id": c.CurrentUser.GetID()},
			bson.M{"secret_token": c.Request.Header.Get(core.ChestAccess)},
			bson.M{"attachments": bson.M{"$elemMatch": bson.M{"current_name": file}}},
		}

		if err := Chest.FindWithOperator("$and", query).Exec(&chest); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		var selector = bson.M{"_id": chest.GetID()}
		var removeQuery = bson.M{"$pull": bson.M{"attachments": bson.M{"current_name": file}}}
		if err := Chest.UpdateQuery(selector, removeQuery, false); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(chest, "Deleted!!", 200, serializers.ChestSerializer{})
	}

	return c.ServerErrorResponse()
}

// Lock ...
func (c ChestsController) Lock(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var chest models.Chest
	if Chest, ok := app.Mapper.GetModel(&chest); ok {
		var query = []bson.M{
			bson.M{"_id": bson.ObjectIdHex(id)},
			bson.M{"owner.$id": c.CurrentUser.GetID()},
			bson.M{"secret_token": c.Request.Header.Get(core.ChestAccess)},
		}

		if err := Chest.FindWithOperator("$and", query).Exec(&chest); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if token, err := auth.GenerateToken(c.CurrentUser.GetID().Hex(), core.ActionLock); err == nil {
			chest.SecretToken = token

			Chest.Update(&chest)
		}

		return c.SuccessResponse("OK", "Locked!!", 200, nil)
	}

	return c.ServerErrorResponse()
}

// Unlock ...
func (c ChestsController) Unlock(id, pin string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var chest models.Chest
	if Chest, ok := app.Mapper.GetModel(&chest); ok {
		if err := Chest.Find(id).Exec(&chest); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(chest.PINHash), []byte(pin)); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if token, err := auth.GenerateToken(c.CurrentUser.GetID().Hex(), core.ActionUnlock); err == nil {
			chest.SecretToken = token

			Chest.Update(&chest)
		}

		return c.SuccessResponse(chest, "Success!!", 200, serializers.ChestSerializer{})
	}

	return c.ErrorResponse(nil, "Server error", 500)
}
