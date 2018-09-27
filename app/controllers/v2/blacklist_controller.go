package v2

import (
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// BlacklistController controller
type BlacklistController struct {
	BaseController
}

// Blacklist show blocked users
func (c BlacklistController) Blacklist() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	// if User, ok := app.Mapper.GetModel(&c.CurrentUser); ok {
	var blacklist []models.User
	var query = bson.M{"saas.name": bson.M{"$in": c.CurrentUser.BlackList}}

	if err := models.Query(c.CurrentUser.GetDocumentName(), query).Select(c.CurrentUser.GetPublicFields()).Exec(&blacklist); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(blacklist, "Success", 200, serializers.UserSerializer{})
}

// AddToBlackList add user to blacklist.
func (c BlacklistController) AddToBlackList(jid string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if jid == c.CurrentUser.Saas.Name {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	var selector = bson.M{"_id": c.CurrentUser.GetID()}
	var query = bson.M{"$addToSet": bson.M{"blacklist": jid}}

	if err := models.UpdateByQuery(&c.CurrentUser, selector, query, false); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("ok", "User was added to blacklist", 200, nil)
}

// RemoveFromBlacklist remove user from blacklist
func (c BlacklistController) RemoveFromBlacklist(jid string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if jid == c.CurrentUser.Saas.Name {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	var selector = bson.M{"_id": c.CurrentUser.GetID()}
	var query = bson.M{"$pull": bson.M{"blacklist": jid}}

	if err := models.UpdateByQuery(&c.CurrentUser, selector, query, false); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("ok", "User was removed from blacklist", 200, nil)
}

// CheckJIDInBlacklist returns code 9050 if jid is in blacklist else returns code 200
func (c BlacklistController) CheckJIDInBlacklist(jid string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	return c.ServerErrorResponse()
}
