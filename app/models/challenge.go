package models

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"time"

	"github.com/Reti500/mgomap"
	"gopkg.in/mgo.v2/bson"
)

// Challenge model
type Challenge struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Title      string         `json:"title" bson:"title" required:"true" validate:"nonzero,max=120,min=6,regexp=^[a-zA-Z0-9 '!?:]*$"`
	UserID     bson.ObjectId  `json:"user_id" bson:"user_id" required:"true"`
	User       interface{}    `json:"user" bson:"user"`
	Rating     int            `json:"-" bson:"rating"`
	SecretText string         `json:"secret_text" bson:"secret_text" validate:"regexp=^[a-zA-Z0-9 ]*$"`
	Properties GameProperties `json:"properties" bson:"properties" validate:"nonzero"`
	Type       string         `json:"type" bson:"type" validate:"nonzero" validate:"regexp=^[a-zA-Z_ ]*$"`
	EndDate    time.Time      `json:"end_date" bson:"end_date" required:"true"`
	Players    []PlayerInfo   `json:"players" bson:"players"`
	IsPublic   bool           `json:"public" bson:"public"`
	Country    string         `json:"country" bson:"country"`

	// Used to retreive collections from DB in aggregate
	Rewards []Reward `json:"rewards" bson:"rewards"`

	// Player is the challenger
	Player PlayerInfo `json:"player" bson:"-"`
	//
	// Ischallenger token
	PlayerToken string `json:"token" bson:"-"`
}

// ChallengeFiles model
type ChallengeFiles struct {
	Attachment Attachment `json:"file" bson:"attachement"`
	Thumbnail  Attachment `json:"thumbnail" bson:"thumbnail"`
}

// ChallengeResponse ...
type ChallengeResponse struct {
	ID  string `json:"id" bson:"-"`
	URL string `json:"url" bson:"-"`
}

// GameProperties ...
type GameProperties struct {
	Intents   int        `json:"intents" bson:"intents" validate:"regexp=^[0-9]*$"`
	Score     string     `json:"score" bson:"score" validate:"regexp=^[a-zA-Z0-9]*$"`
	StartDate CustomTime `json:"start_date" bson:"start_date"`
	EndDate   CustomTime `json:"end_date" bson:"end_date"`
	GameType  string     `json:"type" bson:"game_type" validate:"regexp=^[a-zA-Z0-9]{0,25}$"`
}

// GetDocumentName ...
func (m *Challenge) GetDocumentName() string {
	return "challenges"
}

// GetUser returns the user from the challenge
func (m *Challenge) GetUser() User {
	user := User{}
	User, _ := app.Mapper.GetModel(&user)
	User.Find(m.UserID.Hex()).Exec(&user)
	return user
}

// SetStatus changes status of challenge
func (m *Challenge) SetStatus(status string) {
	m.Status.Code = core.SubscriptionStatus[status]
	m.Status.Name = status
}
