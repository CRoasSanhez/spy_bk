package models

import "github.com/Reti500/mgomap"

// Country model
type Country struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name             string                 `json:"name" bson:"name"`
	PrefISO          string                 `json:"pref_ISO" bson:"pref_ISO"`
	PrefUN           string                 `json:"pref_UN" bson:"pref_UN"`
	NumUN            string                 `json:"num_UN" bson:"num_UN"`
	DialCode         string                 `json:"dialing_code" bson:"dialing_code"`
	MissionLanguages []string               `json:"languages" bson:"languages"`
	ExtraData        map[string]interface{} `json:"extra_data" bson:"extra_data"`
}

// GetDocumentName needed function for Mongo storage with mgomap
func (m *Country) GetDocumentName() string {
	return "countries"
}
