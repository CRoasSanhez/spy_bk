package v2

import (
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

// V2SaasController controller
type V2SaasController struct {
	BaseController
}

// Info return information of saas name
func (c V2SaasController) Info(jid string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var saasUser models.User
	var fields = []string{"user_name", "device.messaging_token", "personal_data"}
	var query = bson.M{"saas.name": jid}

	if err := models.Query(c.CurrentUser.GetDocumentName(), query).Select(fields).Exec(&saasUser); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(saasUser, "Success", 200, serializers.UserSerializer{})
}
