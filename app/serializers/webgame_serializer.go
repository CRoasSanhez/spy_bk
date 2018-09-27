package serializers

import (
	"spyc_backend/app/models"
)

// WebGameSerializer ...
type WebGameSerializer struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	CoverPicture interface{} `json:"cover_picture"`
	Thumbnail    interface{} `json:"thumbnail"`
	Type         string      `json:"type"`
	NameURL      string      `json:"url_name"`
	//Description  string `json:"description"`
	Description interface{} `json:"description"`

	// Lang is the request language
	Lang string `json:"lang"`
}

// WebGamesSerializer ...
type WebGamesSerializer struct {
	WebGames interface{} `json:"webgames"`
	Lang     string      `json:"-"`
}

// Cast ...
func (s WebGameSerializer) Cast(data interface{}) Serializer {
	serializer := new(WebGameSerializer)

	if model, ok := data.(models.WebGame); ok {
		serializer.ID = model.GetID().Hex()
		serializer.Name = model.Name
		serializer.Type = model.Type
		serializer.NameURL = model.NameURL
		serializer.Description = InternationalizeSerializer{}.Get(model.Description, s.Lang)
		serializer.CoverPicture = Serialize(model.Attachment, AttachmentSerializer{
			Parent: model.GetDocumentName(), ParentID: model.GetID().Hex(), Field: "attachment", VerifyURL: true,
		})
		serializer.Thumbnail = Serialize(model.Thumbnail, AttachmentSerializer{
			Parent: model.GetDocumentName(), ParentID: model.GetID().Hex(), Field: "thumbnail", VerifyURL: true,
		})
	}

	return serializer
}

// Cast ...
func (s WebGamesSerializer) Cast(data interface{}) Serializer {
	serializer := new(WebGamesSerializer)

	if model, ok := data.(models.WebGames); ok {
		serializer.WebGames = Serialize(model.WebGames, WebGameSerializer{Lang: s.Lang})
	}
	return serializer
}
