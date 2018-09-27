package models

import (
	"spyc_backend/app/core"

	"github.com/Reti500/mgomap"
)

// AppPath ...
type AppPath struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name       string `json:"name" bson:"name"`
	Section    string `json:"section" bson:"section"`
	URL        string `json:"url" bson:"url"`
	Enviroment string `json:"enviroment" bson:"enviroment"`
}

// GetDocumentName return name collection in DB
func (m *AppPath) GetDocumentName() string {
	return "app_paths"
}

// SetStatus change Mission status according to constants
func (m *AppPath) SetStatus(status string) {
	m.Status.Code = core.SubscriptionStatus[status]
	m.Status.Name = status
}
