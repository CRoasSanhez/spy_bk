package models

import (
	"net/http"

	"github.com/Reti500/mgomap"
	"github.com/mholt/binding"
)

// StoragedPlace ...
type StoragedPlace struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Type    string `json:"type" bson:"type"`
	Hashtag string `json:"hashtag" bson:"hashtag"`
	Files   []File `json:"files" bson:"files"`
}

// GetDocumentName Requeired Method ...
func (p *StoragedPlace) GetDocumentName() string {
	return "storagedplaces"
}

// FieldMap Required Method ...
func (p *StoragedPlace) FieldMap(rep *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&p.Type:    "type",
		&p.Files:   "files",
		&p.Hashtag: "hashtag",
	}
}
