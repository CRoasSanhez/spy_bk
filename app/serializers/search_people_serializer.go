package serializers

import "spyc_backend/app/models"

// SearchPeopleSerializer ...
type SearchPeopleSerializer struct {
	ID         string      `json:"id"`
	Email      string      `json:"email"`
	ImageURL   string      `json:"image_url"`
	Name       string      `json:"name"`
	Number     string      `json:"number"`
	AreFriends interface{} `json:"are_friends"`
}

// Cast ...
func (s SearchPeopleSerializer) Cast(data interface{}) Serializer {
	serializer := new(SearchPeopleSerializer)

	if model, ok := data.(models.User); ok {
		if model.Attachment.HasExpired() {
			model.Attachment.UpdateURLParentField(model.GetDocumentName(), model.GetID().Hex(), "attachment")
		}

		serializer.ID = model.GetID().Hex()
		serializer.Name = model.FullName()
		serializer.Number = model.Device.Number
		serializer.ImageURL = model.Attachment.URL
		serializer.AreFriends = model.ExtraParameters["are_friends"]
	}

	return serializer
}
