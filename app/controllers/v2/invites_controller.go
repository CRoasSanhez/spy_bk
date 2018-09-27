package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

// InvitesController Controller definition
type InvitesController struct {
	BaseController
}

// Index ...
func (c InvitesController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var invitations []models.Invitation
	var query = bson.M{"$and": []bson.M{
		bson.M{"receiver.$id": c.CurrentUser.GetID()},
		bson.M{"status.name": "init"},
	}}

	if err := models.GetDocumentBy(query, &invitations); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(invitations, "Success", 200, serializers.InviteSerializer{})
}

// Respond ...
func (c InvitesController) Respond(id, respond string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	var invite models.Invitation
	var query = bson.M{"$and": []bson.M{
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"receiver.$id": c.CurrentUser.GetID()},
	}}

	if err := models.GetDocumentBy(query, &invite); err != nil {
		return c.ErrorResponse(err, "Invite not found", 400)
	}

	if invite.Status.Name != "init" {
		return c.ErrorResponse(nil, "Invite has a respond", 400)
	}

	switch respond {
	case "accept":
		if invite.Type == "friend" {
			// Make relationship Resource[CurrentUser] -> Sender[TargetUser]
			c.MakeFriend(invite.Resource.Id, invite.Sender.Id)
			// Make relationship Sender[TargetUser] -> Resource[CurrentUser]
			c.MakeFriend(invite.Sender.Id, invite.Resource.Id)
		}

		invite.Status.Name = "accept"

		if User, ok := app.Mapper.GetModel(&models.User{}); ok {
			var user1 models.User

			if err := User.FindRef(invite.Resource).Exec(&user1); err != nil {
				revel.ERROR.Println(err)
			}

			var user2 models.User
			if err := User.FindRef(invite.Sender).Exec(&user2); err != nil {
				revel.ERROR.Println(err)
			}

			if lang := core.FindOnArray(core.CurrentLocales, c.CurrentUser.Device.Language); lang >= 0 {
				c.Request.Locale = c.CurrentUser.Device.Language
			}

			var title = c.Message("notification.friend_request_confirmed_title")
			var message = c.Message("notification.friend_request_confirmed", user1.UserName)
			c.NewNotification("are friends now", user1.Attachment, message, "profile", title, core.FrienRequestConfirmed, user2, nil, true)

		}

		if err := models.UpdateDocument(&invite); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
	case "deny":
		invite.Status.Name = "deny"
		if err := models.UpdateDocument(&invite); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
	default:
		return c.ErrorResponse(nil, "Invalid respond", 400)
	}

	return c.SuccessResponse(invite, "Invitation accepted", 200, serializers.InviteSerializer{})
}

// MakeFriend ...
func (c InvitesController) MakeFriend(user1, user2 interface{}) error {
	var selector = bson.M{"_id": user1.(bson.ObjectId)}
	var query = bson.M{"$addToSet": bson.M{"friends": user2.(bson.ObjectId).Hex()}}

	if err := models.UpdateByQuery(&models.User{}, selector, query, false); err != nil {
		return err
	}

	return nil
}
