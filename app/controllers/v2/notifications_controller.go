package v2

import (
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"gopkg.in/mgo.v2/bson"

	"spyc_backend/app/core"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// NotificationsController ...
type NotificationsController struct {
	BaseController
}

// Index ...
func (c NotificationsController) Index() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var notifications []models.Notification
	var invitations []models.Invitation

	var query = bson.M{"$and": []bson.M{
		{"$or": []bson.M{
			{"to": c.CurrentUser.GetID()},
			{"ids": bson.M{"$in": []string{c.CurrentUser.Device.MessagingToken}}},
		}},
		{"status.name": bson.M{"$ne": core.StatusCompleted}},
		{"type": bson.M{"$ne": core.Message}},
	}}

	var pipe = mgomap.Aggregate{}.Match(query).Limit(20).Sort(bson.M{"created_at": -1})
	if err := models.AggregateQuery(&models.Notification{}, pipe, &notifications); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	//for k, v := range notifications {
	//notifications[k] = notifications[k].Internationalize(c.(revel.Controller))
	//}

	query = bson.M{"$and": []bson.M{
		{"receiver.$id": c.CurrentUser.GetID()},
		{"status.name": "init"},
	}}

	pipe = mgomap.Aggregate{}.Match(query).Limit(20).Sort(bson.M{"created_at": -1})
	if err := models.AggregateQuery(&models.Invitation{}, pipe, &invitations); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	data := &struct {
		Notifications interface{} `json:"notifications"`
		Invitations   interface{} `json:"invitations"`
	}{
		Notifications: serializers.Serialize(notifications, serializers.NotificationSerializer{}),
		Invitations:   serializers.Serialize(invitations, serializers.InviteSerializer{}),
	}

	return c.SuccessResponse(data, "Notifications", 200, nil)
}

// MarkAsReceived set notification status as received
func (c NotificationsController) MarkAsReceived() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var selector = bson.M{"to": c.CurrentUser.GetID(), "status.name": core.StatusInit}
	var query = bson.M{"$set": bson.M{"status.name": core.StatusReceived}}

	if err := models.UpdateByQuery(&models.Notification{}, selector, query, true); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("OK", "Updated", 200, nil)
}

// MarkAsSeen set notification status as seen
func (c NotificationsController) MarkAsSeen(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	var selector = bson.M{"_id": bson.ObjectIdHex(id), "to": c.CurrentUser.GetID()}
	var query = bson.M{"$set": bson.M{"status.name": core.StatusSeen}}

	if err := models.UpdateByQuery(models.Notification{}, selector, query, true); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	selector = bson.M{"_id": bson.ObjectIdHex(id)}
	query = bson.M{"$pull": bson.M{"ids": c.CurrentUser.Device.MessagingToken}}

	if err := models.UpdateByQuery(models.Notification{}, selector, query, true); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("OK", "Updated", 200, nil)
}

// MarkAsCompleted set notification status as completed
func (c NotificationsController) MarkAsCompleted(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(id) {
		return c.ErrorResponse(nil, "Invalid", 400)
	}

	var selector = bson.M{"_id": bson.ObjectIdHex(id), "to": c.CurrentUser.GetID()}
	var query = bson.M{"$set": bson.M{"status.name": core.StatusCompleted}}

	if err := models.UpdateByQuery(&models.Notification{}, selector, query, true); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	selector = bson.M{"_id": bson.ObjectIdHex(id)}
	query = bson.M{"$pull": bson.M{"ids": c.CurrentUser.Device.MessagingToken}}

	if err := models.UpdateByQuery(&models.Notification{}, selector, query, true); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse("OK", "Updated", 200, nil)
}
