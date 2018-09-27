package serializers

import "spyc_backend/app/models"

// CommentOwnerSerializer serializer
type CommentOwnerSerializer struct {
	ID       string `json:"id,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

// Cast ...
func (s CommentOwnerSerializer) Cast(data interface{}) Serializer {
	serializer := new(CommentOwnerSerializer)

	if model, ok := data.(models.User); ok {
		if model.Attachment.HasExpired() {
			model.Attachment.UpdateURLParentField(model.GetDocumentName(), model.GetID().Hex(), "attachment")
		}

		serializer.ID = model.GetID().Hex()
		serializer.ImageURL = model.Attachment.URL
		serializer.UserName = model.UserName
	}

	return serializer
}
