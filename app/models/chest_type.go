package models

import "github.com/Reti500/mgomap"

// ChestType ...
type ChestType struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Description string `json:"description" bson:"description"`
	ImageURL    string `json:"image_url" bson:"image_url"`
	Price       int64  `json:"price" bson:"price"`
	Size        int64  `json:"size" bson:"size"`
}
