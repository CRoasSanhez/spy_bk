package models

import (
	"time"

	"spyc_backend/app/core"

	"github.com/Reti500/mgomap"
)

// Token model
type Token struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Text        string      `bson:"text"`
	Type        string      `bson:"type"`
	ResouceType string      `bson:"resource_type"`
	ResourceID  interface{} `bson:"resource_id"`
	Expire      time.Time   `bson:"expire"`
	Resource    interface{} `bson:"resource"`
}

// GetDocumentName ...
func (u *Token) GetDocumentName() string {
	return "token"
}

// NewToken  generate a new token and return the token
func NewToken(tokenType string, resource string, id interface{}, exprire time.Time) (string, error) {
	var t = Token{
		Text:        core.GenerateToken(core.WebLetterRunes, 64),
		Type:        tokenType,
		ResouceType: resource,
		ResourceID:  id,
		Expire:      exprire,
	}

	if err := CreateDocument(&t); err != nil {
		return "", err
	}

	return t.Text, nil
}
