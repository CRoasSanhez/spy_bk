package models

import (
	"errors"

	"github.com/Reti500/mgomap"
	"gopkg.in/mgo.v2"
)

// Chest model
type Chest struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Attachments []Attachment `bson:"attachments"`
	Geolocation *Geo         `bson:"geolocation" validate:"nonzero"`
	PINHash     string       `bson:"pin_hash"`
	Name        string       `bson:"name" validate:"nonzero" validate:"regexp=^[a-zA-Z0-9.():¿?=!¡ ]*$"`
	Owner       *mgo.DBRef   `bson:"owner"`
	Tags        []string     `bson:"tags" validate:"min=1"`
	Type        ChestType    `bson:"chest_type"`
	Size        int64        `bson:"size"`
	Current     int64        `bson:"current"`
	SecretToken string       `bson:"secret_token"`

	// Not saved fields
	Token string `json:"-" bson:"-"`
	PIN   string `json:"pin" bson:"-" validate:"nonzero"`
}

// GetDocumentName ...
func (m *Chest) GetDocumentName() string {
	return "chests"
}

// GeneratePIN create a hash for pin
func (m *Chest) GeneratePIN() (err error) {
	if len(m.PIN) <= 0 {
		return errors.New("Invalid pin")
	}

	if hash, err := MD5Crypt(m.PIN); err == nil {
		m.PINHash = string(hash)
		return err
	}

	return
}
