package models

import (
	"time"

	"github.com/Reti500/mgomap"
)

// Map model
type Map struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name     string    `bson:"name"`
	Coords   string    `bson:"coords"`
	URL      string    `bson:"url"`
	ExpireOn time.Time `bson:"expire_on"`
}

// GetDocumentName Required Method ...
func (m *Map) GetDocumentName() string {
	return "maps"
}
