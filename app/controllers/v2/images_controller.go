package v2

import (
	"spyc_backend/app/core"
	"spyc_backend/app/models"

	"github.com/revel/revel"
)

type ImagesController struct {
	BaseController
}

// GetResource returns the image url fr challenge based on the given type
// and size (normal, thumb)
func (c *ImagesController) GetResource(gameType, size string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var webgame models.WebGame

	webgame = webgame.FindByNameURL(gameType)

	if webgame.Name == "" {
		return c.ErrorResponse(nil, c.Message("error.notFound", "Game", ""), core.ValidationStatus[core.StatusNotFound])
	}

	url := webgame.FindImageURL(size)

	return c.SuccessResponse(url, size, 0, nil)
}
