package serializers

import (
	"spyc_backend/app/models"
)

// NotificationSerializer ...
type NotificationSerializer struct {
	From      interface{} `json:"from"`
	Resource  interface{} `json:"resource"`
	Id        string      `json:"id"`
	ImageURL  string      `json:"image_url"`
	Type      string      `json:"type"`
	Action    string      `json:"action"`
	Message   string      `json:"message"`
	Screen    string      `json:"screen"`
	Title     string      `json:"title"`
	Status    interface{} `json:"status"`
	CreatedAt interface{} `json:"created_at"`
}

// Cast ...
func (s NotificationSerializer) Cast(data interface{}) Serializer {
	serializer := new(NotificationSerializer)

	if model, ok := data.(models.Notification); ok {
		serializer.Id = model.GetID().Hex()
		serializer.From = model.From
		serializer.Resource = model.Resource
		serializer.Type = model.Type
		serializer.Action = model.Action
		serializer.Message = model.Message
		serializer.Screen = model.Screen
		serializer.Title = model.Title
		serializer.Status = model.Status
		serializer.CreatedAt = model.CreatedAt.Unix()

		if model.Attachment.PATH != "" {
			serializer.ImageURL = model.Attachment.GetURL()
		} else {
			serializer.ImageURL = model.ImageURL
		}
	}

	return serializer
}
