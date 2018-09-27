package serializers

import (
	"spyc_backend/app/core"
	"spyc_backend/app/models"
)

// TargetSerializer ...
type TargetSerializer struct {
	ID          string          `json:"id"`
	Name        interface{}     `json:"name"`
	Description interface{}     `json:"description"`
	StartDate   string          `json:"start_date"`
	EndDate     string          `json:"end_date"`
	Image       interface{}     `json:"image"`
	WebURL      string          `json:"web_url"`
	Type        string          `json:"type"`
	Question    models.Question `json:"question"`
	Lang        string          `json:"lang"`
}

// Cast ...
func (s TargetSerializer) Cast(data interface{}) Serializer {
	serializer := new(TargetSerializer)

	if model, ok := data.(models.Target); ok {
		serializer.ID = model.GetID().Hex()
		serializer.Name = InternationalizeSerializer{}.Get(model.Name, s.Lang)
		serializer.Description = InternationalizeSerializer{}.Get(model.Description, s.Lang)
		serializer.StartDate = model.StartDate.Format(core.MXTimeFormat)
		serializer.EndDate = model.EndDate.Format(core.MXTimeFormat)
		serializer.WebURL = model.WebURL
		serializer.Type = model.Type
		serializer.Question = model.Question
		serializer.Image = Serialize(model.Attachment, AttachmentSerializer{
			Parent: model.GetDocumentName(), ParentID: model.GetID().Hex(), Field: "attachement", VerifyURL: true,
		})
	}

	return serializer
}
