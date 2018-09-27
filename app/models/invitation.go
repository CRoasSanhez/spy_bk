package models

import (
	"spyc_backend/app/core"

	"github.com/Reti500/mgomap"
	mgo "gopkg.in/mgo.v2"
)

// Invitation model
type Invitation struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Sender       *mgo.DBRef `bson:"sender"`
	Receiver     *mgo.DBRef `bson:"receiver"`
	Resource     *mgo.DBRef `bson:"resource"`
	ResourceType string     `bson:"resource_type"`
	Type         string     `bson:"type"`
}

// GetDocumentName Requeired Method ...
func (m *Invitation) GetDocumentName() string {
	return "invitations"
}

// SendPush ...
func (m Invitation) SendPush(to string, ids []string) (*core.FcmResponseStatus, error) {
	if resp, err := core.NewFCMClient(
		to,
		core.PriorityHigh,
		core.FCMNotification{
			Badge: "",
			Body:  "",
			Icon:  "",
			Sound: "",
			Title: "",
		},
		ids,
		&struct {
			Type     string      `json:"type"`
			Id       string      `json:"id"`
			Image    string      `json:"image"`
			Message  string      `json:"message"`
			Action   string      `json:"action"`
			From     interface{} `json:"actor"`
			Resource interface{} `json:"resource"`
		}{m.Type, m.GetID().Hex(), "", "", "", "", ""},
	).Send(); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}
