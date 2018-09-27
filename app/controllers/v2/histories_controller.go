package v2

import (
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

// HistoriesController ...
type HistoriesController struct {
	BaseController
}

// Index returns all histories for current user
func (c HistoriesController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	histories, _ := models.History{}.FindHistories(c.CurrentUser.Friends, c.CurrentUser.GetID().Hex())

	return c.SuccessResponse(histories, core.StatusSuccess, 200, serializers.GroupHistoriesSerializer{})
}

// Create ...
func (c HistoriesController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var history = models.History{
		Anim:   c.Params.Get("anim"),
		Color:  c.Params.Get("color"),
		Owner:  c.CurrentUser.GetID().Hex(),
		Public: c.Params.Get("public") == core.TrueString,
		Text:   c.Params.Get("text"),
		Type:   c.Params.Get("type"),
		Users:  strings.Split(c.Params.Get("users"), ","),
	}

	if err := history.Save(); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if len(c.Params.Files["attachment"]) > 0 {
		history.UpdateAttachment(c.CurrentUser, c.Params.Files["attachment"][0])
	}

	return c.SuccessResponse(history, core.StatusSuccess, 200, serializers.HistoriesSerializer{})
}

// AddView count one view for history
func (c HistoriesController) AddView(viewedtime int, id, completed string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "invalid", 400)
	}

	err := models.History{}.CountView(completed == core.TrueString, viewedtime, id, c.CurrentUser.GetID().Hex())
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("ok", "ok", 200, nil)
}
