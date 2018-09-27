package models

import "github.com/Reti500/mgomap"

// Message model for chat messages
type Message struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Sender     string `json:"sender" bson:"sender"`
	Receiver   string `json:"receiver" bson:"receiver"`
	Stanza     string `json:"stanza" bson:"stanza"`
	Message    string `json:"message" bson:"message"`
	EjabberdID string `json:"ejabberd_id" bson:"ejabberd_id"`
}

// GetDocumentName ...
func (u *Message) GetDocumentName() string {
	return "messages"
}
