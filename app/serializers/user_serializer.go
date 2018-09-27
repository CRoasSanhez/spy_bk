package serializers

import "spyc_backend/app/models"

// UserSerializer ...
type UserSerializer struct {
	ID           string        `json:"id,omitempty"`
	Email        string        `json:"email,omitempty"`
	UserName     string        `json:"user_name,omitempty"`
	ImageURL     string        `json:"image_url,omitempty"`
	PersonalData interface{}   `json:"personal_data,omitempty"`
	Status       interface{}   `json:"status,omitempty"`
	Device       models.Device `json:"device,omitempty"`
	Saas         interface{}   `json:"saas,omitempty"`
	Token        string        `json:"token,omitempty"`
}

// Cast ...
func (s UserSerializer) Cast(data interface{}) Serializer {
	serializer := new(UserSerializer)

	if model, ok := data.(models.User); ok {
		if model.Attachment.HasExpired() {
			model.Attachment.UpdateURLParentField(model.GetDocumentName(), model.GetID().Hex(), "attachment")
		}

		serializer.ID = model.GetID().Hex()
		serializer.Email = model.Email
		serializer.UserName = model.UserName
		serializer.Device = model.Device
		serializer.PersonalData = model.PersonalData
		serializer.Saas = model.Saas
		serializer.Status = model.Status
		serializer.Token = s.Token
		serializer.ImageURL = model.Attachment.URL
	}

	return serializer
}

// UserBaseSerializer has the base user fields to show
type UserBaseSerializer struct {
	//ID           string        `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	UserName string `json:"user_name,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// Cast ...
func (s UserBaseSerializer) Cast(data interface{}) Serializer {
	serializer := new(UserBaseSerializer)

	if model, ok := data.(models.User); ok {
		if model.Attachment.HasExpired() {
			model.Attachment.UpdateURLParentField(model.GetDocumentName(), model.GetID().Hex(), "attachment")
		}

		serializer.Email = model.Email
		serializer.UserName = model.UserName
		serializer.ImageURL = model.Attachment.URL
	}
	return serializer
}
