package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/models"

	"spyc_backend/app/core"

	"github.com/revel/revel"
)

// SanctionsController ...
type SanctionsController struct {
	BaseController
}

// SactionUser enroutes the sanction for the user based on the given sanctionType
// Current user is thebad user
// userID: is the user wich the screenshot is taken from
func (c SanctionsController) SactionUser(sanctionType, userID string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var affectedUser models.User

	if User, ok := app.Mapper.GetModel(&affectedUser); ok {
		var ok = false

		if err := User.Find(userID).Exec(&affectedUser); err != nil {
			return c.ErrorResponse(nil, c.Message("error.notFound", "User"), core.ValidationStatus[core.StatusError])
		}

		switch sanctionType {
		case core.SanctionScreenShoot:
			c.CurrentUser.AddTransaction(-200, "min")
			affectedUser.AddTransaction(200, "add")
			ok = true
		case core.SanctionMissionQuit:
			c.CurrentUser.AddTransaction(-200, "min")
			ok = true
		default:
			break
		}

		if !ok {
			return c.ErrorResponse(nil, c.Message("error.update", "user"), core.ValidationStatus[core.StatusError])
		}

		return c.SuccessResponse("success", "success", core.ValidationStatus[core.StatusSuccess], nil)
	}

	return c.ServerErrorResponse()
}
