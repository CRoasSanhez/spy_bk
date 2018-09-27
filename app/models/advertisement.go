package models

import "github.com/Reti500/mgomap"

// Advertisement ...
type Advertisement struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Attachment Attachment `json:"-" bson:"attachment"`
	Sponsor    string     `josn:"-" bson:"sponsor"`
	SponsorID  string     `json:"sponsor_id" bson:"sponsor_id"`
	Views      []Views    `json:"views" bson:"views"`
	Tags       []string   `json:"tags" bson:"tags"`
	Name       string     `json:"name" bson:"name"`
	AdType     string     `json:"type" bson:"type"`

	// Delay in seconds
	// Time before skip add
	Delay int `json:"delay" bson:"delay"`
}

// GetDocumentName ...
func (m *Advertisement) GetDocumentName() string {
	return "advertisement"
}

// FindUserInViews ...
func (m *Advertisement) FindUserInViews(userID string) Views {
	for i := 0; i < len(m.Views); i++ {
		if m.Views[i].UserID == userID {
			return m.Views[i]
		}
	}
	return Views{}
}
