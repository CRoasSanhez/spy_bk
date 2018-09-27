package serializers

import (
	"spyc_backend/app"
	"spyc_backend/app/models"
)

// PINSerializer serializer
type PINSerializer struct {
	ID            string      `json:"id"`
	Title         interface{} `json:"title"`
	Description   interface{} `json:"decription"`
	Attachment    interface{} `json:"attachment"`
	Advertisement interface{} `json:"advertisement"`
	Link          string      `json:"link"`
	Tags          []string    `json:"tags"`
	Lang          string      `json:"lang"`
}

// Cast ...
func (s PINSerializer) Cast(data interface{}) Serializer {
	var serializer = new(PINSerializer)

	if model, ok := data.(models.PINAdvertisement); ok {
		serializer.ID = model.ID.Hex()
		serializer.Title = model.Titles.Get(s.Lang)
		serializer.Description = model.Descriptions.Get(s.Lang)
		serializer.Link = model.Link
		serializer.Tags = model.Tags
		serializer.Lang = s.Lang

		var advertisement models.Advertisement
		if Advertisement, ok := app.Mapper.GetModel(&advertisement); ok {
			if err := Advertisement.Find(model.AdvertisementID).Exec(&advertisement); err == nil {
				serializer.Advertisement = Serialize(advertisement, AdvertisementSerializer{})
			}
		}

		if model.Attachments.Get(s.Lang) != nil {
			serializer.Attachment = Serialize(model.Attachments.Get(s.Lang), AttachmentSerializer{})
		}
	}

	return serializer
}
