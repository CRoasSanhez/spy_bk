package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"github.com/revel/revel"
)

// StatsController controller
type StatsController struct {
	BaseController
}

// Index returns current user stats
func (c StatsController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var stats models.Stats

	if Stats, ok := app.Mapper.GetModel(&stats); ok {
		if err := Stats.FindBy("owner", c.CurrentUser.GetID()).Exec(&stats); err != nil {
			return c.ErrorResponse(err, err.Error(), 200)
		}

		return c.SuccessResponse(stats, "ok", 200, serializers.StatsSerializer{})
	}

	return c.ServerErrorResponse()
}

func (c StatsController) Replicate() revel.Result {
	models.ReplicateStats()

	return c.SuccessResponse("ok", "ok", 200, nil)
}
