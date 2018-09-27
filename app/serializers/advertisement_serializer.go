package serializers

import "spyc_backend/app/models"

// AdvertisementSerializer serializer
type AdvertisementSerializer struct {
	ID         string      `json:"id"`
	Attachment interface{} `json:"attachment"`
	Sponsor    string      `josn:"sponsor"`
	Views      interface{} `json:"views"`
	Tags       interface{} `json:"tags"`
	Name       string      `json:"name"`
	AdType     string      `json:"type"`
	Delay      int         `json:"delay"`
}

// Cast ...
func (s AdvertisementSerializer) Cast(data interface{}) Serializer {
	var serializer = new(AdvertisementSerializer)

	if model, ok := data.(models.Advertisement); ok {
		serializer.ID = model.ID.Hex()
		serializer.AdType = model.AdType
		serializer.Delay = model.Delay
		serializer.Name = model.Name
		serializer.Sponsor = model.Sponsor
		serializer.Tags = model.Tags
		serializer.Views = model.Views
		serializer.Attachment = Serialize(model.Attachment, AttachmentSerializer{
			Parent: model.GetDocumentName(), ParentID: model.GetID().Hex(), Field: "attachment", VerifyURL: true,
		})
	}

	return serializer
}
