package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

// PINController controller
type PINController struct {
	BaseController
}

// Index return all pin for current user
func (c PINController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	pins, err := models.PINAdvertisement{}.GetPinsForUser(c.CurrentUser)
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(pins, core.StatusSuccess, 200, serializers.PINSerializer{Lang: c.Request.Locale})
}

// Create generate a new PIN
func (c PINController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var pin models.PINAdvertisement
	if err := c.PermitParams(&pin, false, "titles", "descriptions", "advertisement_id",
		"link", "target_min_age", "target_max_age", "target_gender", "scope", "tags", "geolocation"); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if Pin, ok := app.Mapper.GetModel(&pin); ok {
		if err := Pin.Create(&pin); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(pin, core.StatusSuccess, 200, serializers.PINSerializer{Lang: c.Request.Locale})
	}

	return c.ServerErrorResponse()
}

// Update ...
func (c PINController) Update(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	var pin models.PINAdvertisement
	if Pin, ok := app.Mapper.GetModel(&pin); ok {
		if err := Pin.Find(id).Exec(&pin); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if err := c.PermitParams(&pin, false, "titles", "descriptions",
			"link", "target_min_age", "target_max_age", "target_gender", "scope", "tags", "geolocation"); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		if err := Pin.Update(&pin); err != nil {
			return c.ErrorResponse(err, err.Error(), 200)
		}

		return c.SuccessResponse(pin, core.StatusSuccess, 200, serializers.PINSerializer{})
	}

	return c.ServerErrorResponse()
}

// Delete remove pin
func (c PINController) Delete(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	var pin models.PINAdvertisement
	if Pin, ok := app.Mapper.GetModel(&pin); ok {
		if err := Pin.Find(id).Exec(&pin); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		pin.Deleted = true

		if err := Pin.Update(&pin); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}

		return c.SuccessResponse(pin, core.StatusSuccess, 200, serializers.PINSerializer{})
	}

	return c.ServerErrorResponse()
}

// PinClicked count a one click for pin
func (c PINController) PinClicked(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	err := models.PINAdvertisement{}.React(core.ActionClick, id, c.CurrentUser.GetID().Hex())
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("ok", "ok", 200, nil)
}

// PinViewed count a one click for pin
func (c PINController) PinViewed(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	err := models.PINAdvertisement{}.React(core.ActionView, id, c.CurrentUser.GetID().Hex())
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("ok", "ok", 200, nil)
}
