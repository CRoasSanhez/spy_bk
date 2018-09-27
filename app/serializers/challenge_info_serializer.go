package serializers

import (
	"log"
	"reflect"
	"spyc_backend/app/models"
)

// ChallengeInfoSerializer ...
type ChallengeInfoSerializer struct {
	User     models.PlayerInfo     `json:"user"`
	Player   models.PlayerInfo     `json:"player"`
	GameData models.GameProperties `json:"game_data"`
	Token    string                `json:"token"`
}

// Cast ...
func (s ChallengeInfoSerializer) Cast(data interface{}) Serializer {
	serializer := new(ChallengeInfoSerializer)
	if model, ok := data.(models.Challenge); ok {

		serializer.Player = model.Player
		serializer.GameData = model.Properties
		serializer.Token = model.PlayerToken
		user := model.GetUser()
		serializer.User = models.PlayerInfo{UserName: user.UserName}
	}
	return serializer
}

// EncryptModel ...
func EncryptModel(model models.Challenge) {

	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	t := reflect.TypeOf(model)

	s := reflect.ValueOf(&model).Elem()

	// Get the type and kind of our user variable
	//fmt.Println("Type:", t.Name())
	//fmt.Println("Kind:", t.Kind())

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field
		field := t.Field(i)

		// get the field
		f := s.Field(i)
		log.Print(field)

		log.Print(f.Interface())

		//fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}
}
