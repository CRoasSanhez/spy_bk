package serializers

import (
	"spyc_backend/app"
	"spyc_backend/app/models"
)

// FriendSerializer ...
type FriendSerializer struct {
	ID           string              `json:"id"`
	ImageURL     string              `json:"image_url"`
	UserName     string              `json:"user_name"`
	Device       models.Device       `json:"device"`
	PersonalData models.PersonalData `json:"personal_data"`
	Saas         models.Saas         `json:"saas"`
}

// Cast ...
func (s FriendSerializer) Cast(data interface{}) Serializer {
	serializer := new(FriendSerializer)

	if model, ok := data.(string); ok {
		var user models.User
		if User, ok := app.Mapper.GetModel(&user); ok {
			if err := User.Find(model).Select([]string{
				"attachment", "user_name", "device.number", "personal_data", "saas.name"}).Exec(&user); err == nil {
				if user.Attachment.HasExpired() {
					user.Attachment.UpdateURLParentField(user.GetDocumentName(), user.GetID().Hex(), "attachment")
				}

				serializer.ID = user.GetID().Hex()
				serializer.ImageURL = user.Attachment.URL
				serializer.UserName = user.UserName
				serializer.Device.Number = user.Device.Number
				serializer.PersonalData.FirstName = user.PersonalData.FirstName
				serializer.PersonalData.LastName = user.PersonalData.LastName
				serializer.Saas.Name = user.Saas.Name
			}
		}
	}

	return serializer
}
