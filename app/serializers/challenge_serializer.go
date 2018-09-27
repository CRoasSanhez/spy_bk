package serializers

import (
	"spyc_backend/app"
	"spyc_backend/app/models"

	"github.com/revel/revel"
)

// ChallengeSerializer ...
type ChallengeSerializer struct {
	Title    string      `json:"title"`
	Rating   int         `json:"rating"`
	IsPublic bool        `json:"public"`
	Players  interface{} `json:"players"`
}

// Cast ...
func (s ChallengeSerializer) Cast(data interface{}) Serializer {
	serializer := new(ChallengeSerializer)
	if model, ok := data.(models.Challenge); ok {

		serializer.Title = model.Title
		serializer.Rating = model.Rating
		serializer.IsPublic = model.IsPublic
		serializer.Players = Serialize(model.Players, PlayerSerializer{})
	}
	return serializer
}

// ChallengeStatusSerializer ...
type ChallengeStatusSerializer struct {
	Title   string      `json:"title"`
	Owner   interface{} `json:"owner"`
	Rating  int         `json:"rating"`
	Players interface{} `json:"players"`
	Goal    string      `json:"goal"`
}

// Cast ...
func (s ChallengeStatusSerializer) Cast(data interface{}) Serializer {
	serializer := new(ChallengeStatusSerializer)
	if model, ok := data.(models.Challenge); ok {

		serializer.Title = model.Title
		serializer.Rating = model.Rating
		serializer.Goal = model.Properties.Score
		serializer.Owner = Serialize(FindUser(model.UserID.Hex(), model.GetID().Hex()), UserBaseSerializer{})
		serializer.Players = Serialize(model.Players, PlayerSerializer{})
	}
	return serializer
}

// FindUser finds the owner for the challenge
func FindUser(userID, challengeID string) (user models.User) {
	if User, ok := app.Mapper.GetModel(&models.User{}); ok {
		if err := User.Find(userID).Exec(&user); err != nil {
			revel.ERROR.Printf("ERROR FIND User: %s for Challenge: %s --- %s", userID, challengeID, err.Error())
		}
	}
	return user
}
