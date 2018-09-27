package serializers

import (
	"spyc_backend/app"
	"spyc_backend/app/models"
)

// PlayerSerializer ...
type PlayerSerializer struct {
	ID           string              `json:"user_id"`
	GameName     string              `json:"game"`
	Status       string              `json:"status"`
	Goal         string              `json:"goal"`
	Score        string              `json:"score"`
	ImageURL     string              `json:"user_image"`
	UserName     string              `json:"user_name"`
	PersonalData models.PersonalData `json:"personal_data"`
	Rewards      interface{}         `json:"rewards"`
}

// Cast ...
func (s PlayerSerializer) Cast(data interface{}) Serializer {
	serializer := new(PlayerSerializer)

	if model, ok := data.(models.PlayerInfo); ok {
		var user models.User
		if User, ok := app.Mapper.GetModel(&user); ok {
			if err := User.Find(model.UserID).Exec(&user); err == nil {
				if user.Attachment.HasExpired() {
					user.Attachment.UpdateURLParentField(user.GetDocumentName(), user.GetID().Hex(), "attachment")
				}

				serializer.ID = user.GetID().Hex()
				serializer.Status = model.Status
				serializer.ImageURL = user.Attachment.URL
				serializer.UserName = user.UserName
				serializer.Goal = model.Goal
				serializer.GameName = model.Game
				serializer.Score = model.Score
				serializer.Rewards = Serialize(model.Rewards, RewardSerializer{})
				serializer.PersonalData.FirstName = user.PersonalData.FirstName
				serializer.PersonalData.LastName = user.PersonalData.LastName
			}
		}
	}

	return serializer
}
