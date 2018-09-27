package serializers

import (
	"spyc_backend/app/models"
)

// TargetValidationSerializer ...
type TargetValidationSerializer struct {
	ID                string      `json:"id"`
	TargetID          string      `json:"step_id"`
	UserID            string      `json:"user_id"`
	UserName          string      `json:"user_name"`
	UserScore         string      `json:"user_score"`
	Intents           int         `json:"intents"`
	ResponseImage     interface{} `json:"response_image"`
	TargetScore       string      `json:"step_score"`
	TargetDescription string      `json:"description"`
	Email             string      `json:"email"`
	CompletedAt       interface{} `json:"completed_at"`
}

// Cast ...
func (s TargetValidationSerializer) Cast(data interface{}) Serializer {
	serializer := new(TargetValidationSerializer)

	if model, ok := data.(models.TargetUser); ok {
		serializer.ID = model.GetID().Hex()
		serializer.TargetID = model.TargetID.Hex()
		serializer.UserScore = model.Score
		serializer.Intents = model.Intents
		target := model.FindTarget()
		serializer.TargetScore = target.Score
		serializer.TargetDescription = target.Description.GetString(target.Langs[0])
		user := model.FindUser()
		serializer.UserID = model.UserID.Hex()
		serializer.UserName = user.PersonalData.FirstName + user.PersonalData.LastName
		serializer.Email = user.Email
		serializer.CompletedAt = model.UpdatedAt
		serializer.ResponseImage = Serialize(model.ResponseImage, AttachmentSerializer{
			Parent: model.GetDocumentName(), ParentID: model.GetID().Hex(), Field: "response_image", VerifyURL: true,
		})
	}

	return serializer
}
