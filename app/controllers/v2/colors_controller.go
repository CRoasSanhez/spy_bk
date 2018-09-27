package v2

import (
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"github.com/revel/revel"
)

// ColorsController ...
type ColorsController struct {
	BaseController
}

// Index returns all color available
func (c ColorsController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	colors, err := models.Color{}.GetAll()
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(colors, core.StatusSuccess, 200, serializers.ColorSerializer{})
}
