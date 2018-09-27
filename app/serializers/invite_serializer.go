package serializers

import (
	"spyc_backend/app"
	"spyc_backend/app/models"

	"gopkg.in/mgo.v2/bson"
)

// InviteSerializer ...
type InviteSerializer struct {
	ID     string `json:"id"`
	Sender struct {
		UserName string `json:"user_name"`
		FullName string `json:"full_name"`
		ImageURL string `json:"image_url"`
	} `json:"sender"`
	Type         string      `json:"type"`
	ResourceID   string      `json:"resource_id"`
	ResourceType string      `json:"resource_type"`
	Status       interface{} `json:"status"`
	CreatedAt    interface{} `json:"created_at"`
}

// Cast ...
func (s InviteSerializer) Cast(data interface{}) Serializer {
	serializer := new(InviteSerializer)

	if model, ok := data.(models.Invitation); ok {
		var user models.User
		if User, ok := app.Mapper.GetModel(&user); ok {
			if err := User.FindRef(model.Sender).Exec(&user); err == nil {
				if user.Attachment.HasExpired() {
					user.Attachment.UpdateURLParentField(model.GetDocumentName(), model.GetID().Hex(), "attachment")
				}

				serializer.Sender.UserName = user.UserName
				serializer.Sender.FullName = user.FullName()
				serializer.Sender.ImageURL = user.Attachment.URL
			}
		}

		serializer.ID = model.GetID().Hex()
		serializer.Type = model.Type
		serializer.ResourceID = model.Resource.Id.(bson.ObjectId).Hex()
		serializer.ResourceType = model.ResourceType
		serializer.Status = model.Status
		serializer.CreatedAt = model.CreatedAt.Unix()
	}

	return serializer
}
