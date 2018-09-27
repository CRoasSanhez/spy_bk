package serializers

import "spyc_backend/app/models"

// ChestSerializer serializer
type ChestSerializer struct {
	ID string `json:"id"`

	Geo         *models.Geo `json:"geolocation"`
	Name        string      `json:"name"`
	Total       int64       `json:"total"`
	Current     int64       `json:"current"`
	Left        int64       `json:"left"`
	Token       string      `json:"token"`
	Attachments interface{} `json:"attachments"`
}

// Cast ...
func (s ChestSerializer) Cast(data interface{}) Serializer {
	serializer := new(ChestSerializer)

	if model, ok := data.(models.Chest); ok {
		serializer.ID = model.GetID().Hex()
		serializer.Name = model.Name
		serializer.Total = model.Size
		serializer.Current = model.Current
		serializer.Left = model.Size - model.Current
		serializer.Geo = model.Geolocation
		serializer.Token = model.SecretToken
		serializer.Attachments = Serialize(model.Attachments, AttachmentSerializer{})
	}

	return serializer
}
