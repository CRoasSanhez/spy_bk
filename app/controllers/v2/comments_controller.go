package v2

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

// CommentsController controller
type CommentsController struct {
	BaseController
}

// React [POST, /comment/:id/react]
func (c CommentsController) React(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	return c.ServerErrorResponse()
}
