package models

import (
	"net/http"

	"github.com/Reti500/mgomap"
	"github.com/mholt/binding"
)

// StoragedPlacesUser ...
// Visibility represents the access level to the storagedPlace
type StoragedPlacesUser struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	StoragedPlaceID string `json:"storagedPlace_id" bson:"storagedPlace_id"`
	UserID          string `json:"user_id" bson:"user_id"`
	Visibility      string `json:"visibility" bson:"visibility"`
	//Status          string `json:"status" bson:"status"`
}

// GetDocumentName required method...
func (spu *StoragedPlacesUser) GetDocumentName() string {
	return "storagedplaces_users"
}

// FieldMap required method...
func (spu *StoragedPlacesUser) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&spu.StoragedPlaceID: "storagedPlace_id",
		&spu.UserID:          "user_id",
		&spu.Visibility:      "visibility",
		//&spu.Status:          "status",
	}
}
