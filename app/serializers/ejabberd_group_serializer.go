package serializers

import "spyc_backend/app/models"

// EjabberdGroupSerializer serializer
type EjabberdGroupSerializer struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	JID      string      `json:"jid"`
	ImageURL string      `json:"image_url"`
	Members  interface{} `json:"members"`
}

// Cast make interface to EjabberdGroup model
func (s EjabberdGroupSerializer) Cast(data interface{}) Serializer {
	serializer := new(EjabberdGroupSerializer)

	if model, ok := data.(models.EjabberdGroup); ok {
		if model.Attachment.HasExpired() {
			model.Attachment.UpdateURLParentField(model.GetDocumentName(), model.GetID().Hex(), "attachment")
		}

		serializer.ID = model.GetID().Hex()
		serializer.Name = model.Name
		serializer.JID = model.JID
		serializer.ImageURL = model.Attachment.URL
		serializer.Members = model.Members
	}

	return serializer
}
