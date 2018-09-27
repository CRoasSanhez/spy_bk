package serializers

import "spyc_backend/app/models"

// ColorSerializer ...
type ColorSerializer struct {
	Name       string      `json:"color"`
	Hex        string      `json:"hex"`
	Attachment interface{} `json:"attachment"`
}

// Cast ...
func (c ColorSerializer) Cast(data interface{}) Serializer {
	var serializer = new(ColorSerializer)

	if model, ok := data.(models.Color); ok {
		serializer.Name = model.Name
		serializer.Hex = model.Hexadecimal

		serializer.Attachment = Serialize(model.Attachment, AttachmentSerializer{
			Parent: model.GetDocumentName(), ParentID: model.GetID().Hex(), Field: "attachment", VerifyURL: true,
		})
	}

	return serializer
}
