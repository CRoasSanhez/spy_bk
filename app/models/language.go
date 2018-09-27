package models

import "github.com/Reti500/mgomap"

// Language ...
type Language struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name  Internationalized `json:"name" bson:"name"`
	Code  string            `json:"code" bson:"code"`
	Langs []string          `json:"-" bson:"langs"`
}

// GetDocumentName returns the collection name
func (m *Language) GetDocumentName() string {
	return "languages"
}
