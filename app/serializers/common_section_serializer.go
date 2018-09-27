package serializers

import "spyc_backend/app/models"

// CommonSectionSerializer ...
type CommonSectionSerializer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	AccessTime int    `json:"access_time"`
}

// Cast ...
func (s CommonSectionSerializer) Cast(data interface{}) Serializer {
	var serializer = new(CommonSectionSerializer)

	if model, ok := data.(models.CommonSection); ok {
		serializer.ID = model.ID.Hex()
		serializer.Name = model.Name
		serializer.AccessTime = model.AccessTime
	}

	return serializer
}
