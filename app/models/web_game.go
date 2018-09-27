package models

import (
	"regexp"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
	validation "github.com/go-ozzo/ozzo-validation"
)

// WebGame ...
type WebGame struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name    string `json:"name" bson:"name"`
	NameURL string `json:"url_name" bson:"name_url"`
	//Description string     `json:"description" bson:"description" validate:"nonzero,regexp=^[a-zA-Z0-9 ]*$"`
	Description Internationalized `json:"description" bson:"description"`
	Attachment  Attachment        `json:"-" bson:"attachment"`
	Thumbnail   Attachment        `json:"-" bson:"thumbnail"`
	Type        string            `json:"type" bson:"type" validate:"regexp=^[a-zA-Z0-9 ]*$"`
	Langs       []string          `json:"langs" bson:"languages"`
}

// WebGames ...
type WebGames struct {
	WebGames []WebGame `json:"webgames"`
}

// GetDocumentName return collection name
func (m *WebGame) GetDocumentName() string {
	return "games"
}

// WebGameFinish ...
type WebGameFinish struct {
	Data   string `json:"data" bson:"-" validate:"min=0,max=30,regexp=^[a-zA-Z0-9]*$"`
	Status string `json:"status" bson:"-" validate:"max=20,regexp=^[a-zA-Z0-9]*$"`
}

// WebGameInfo is the response to gameInfo request
type WebGameInfo struct {
	User     PlayerInfo `json:"user" bson:"-"`
	Player   PlayerInfo `json:"player" bson:"-"`
	GameData GameData   `json:"game_data" bson:"-"`
	Token    string     `json:"token" bson:"-"`
}

// PlayerInfo is the general info of player
type PlayerInfo struct {
	UserName    string    `json:"userName" bson:"-"`
	UserID      string    `json:"user_id" bson:"user_id"`
	Status      string    `json:"status" bson:"status"`
	CompletedAt time.Time `json:"completed_at" bson:"completed_at"`
	Score       string    `json:"score" bson:"score"`

	// Used for chat
	Goal  string `json:"goal"`
	Game  string `json:"game"`
	Token string `json:"-" bson:"token"`

	// Used to return possible rewards when asking challenge status
	Rewards []Reward `json:"-" bson:"-"`
}

// GameData is the data of the requested game
type GameData struct {
	Score   string `json:"score"`
	Lang    string `json:"lang"`
	Intents string `json:"intents"`
}

// Validate ...
func (m WebGame) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z0-9 ]*$"))),
		validation.Field(&m.NameURL, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z0-9 ]*$"))),
		validation.Field(&m.Type, validation.Match(regexp.MustCompile("^[a-zA-Z0-9 ]*$"))),
	)
}

// SetStatus changes status of challenge
func (m *WebGame) SetStatus(status string) {
	m.Status.Code = core.GameStatus[status]
	m.Status.Name = status
}

// FindWebGamesByStatus retreives webgames on the given status
func (m *WebGame) FindWebGamesByStatus(status string) []WebGame {
	webgames := []WebGame{}
	WebGame, _ := app.Mapper.GetModel(&WebGame{})
	var query = []bson.M{
		bson.M{"status.name": status},
	}
	WebGame.FindWithOperator("$and", query).Exec(&webgames)
	return webgames
}

// FindImageURL returns the image url  based on the given type
func (m *WebGame) FindImageURL(imgType string) (imageURL string) {

	if imgType == "normal" {
		m.Attachment.UpdateURLParentField(m.GetDocumentName(), m.GetID().Hex(), "attachment")
		return m.Attachment.URL
	}

	m.Thumbnail.UpdateURLParentField(m.GetDocumentName(), m.GetID().Hex(), "thumbnail")
	return m.Thumbnail.URL
}

// FindByNameURL returns the active game with the givven nameURL
func (m *WebGame) FindByNameURL(nameURL string) WebGame {
	webgame := WebGame{}
	WebGame, _ := app.Mapper.GetModel(&webgame)
	var query = []bson.M{
		bson.M{"status.name": core.StatusActive},
		bson.M{"name_url": nameURL},
	}
	WebGame.FindWithOperator("$and", query).Exec(&webgame)
	return webgame
}

// Validate player Info
func (m PlayerInfo) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Status, validation.Match(regexp.MustCompile("^[a-zA-Z0-9 ]*$"))),
		validation.Field(&m.Score, validation.Match(regexp.MustCompile("^[a-zA-Z0-9 ]*$"))),
	)
}

// Validate WebGameFinish
func (m WebGameFinish) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Data, validation.Length(0, 30), validation.Match(regexp.MustCompile("^[a-zA-Z0-9]*$"))),
		validation.Field(&m.Status, validation.Required, validation.Length(1, 20), validation.Match(regexp.MustCompile("^[a-zA-Z0-9]*$"))),
	)
}
