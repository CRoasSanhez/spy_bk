package models

import (
	"github.com/Reti500/mgomap"
)

// EjabberdUser model
type EjabberdUser struct {
	JID   string `json:"jid" bson:"jid"`
	Nick  string `json:"nick" bson:"nick"`
	Admin bool   `json:"admin" bson:"admin"`
}

// EjabberdGroup model
type EjabberdGroup struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name       string          `bson:"name"`
	JID        string          `bson:"jid"`
	Owner      string          `bson:"owner"`
	Attachment Attachment      `bson:"attachment"`
	Members    []*EjabberdUser `bson:"members"`
}

// GetDocumentName ...
func (u *EjabberdGroup) GetDocumentName() string {
	return "ejabberdgroups"
}

// AddMember add user to group
func (u *EjabberdGroup) AddMember(jid string, nick string, admin bool) {
	if len(u.Members) <= 0 {
		u.Members = make([]*EjabberdUser, 1)
	}

	u.Members = append(u.Members, &EjabberdUser{jid, nick, admin})
}

// RemoveMember remove user to group
func (u *EjabberdGroup) RemoveMember(jid string) {
	if len(u.Members) <= 0 || len(jid) <= 0 {
		return
	}

	for k, v := range u.Members {
		if v.JID == jid {
			u.Members = append(u.Members[:k], u.Members[k+1:]...)
			return
		}
	}
}
