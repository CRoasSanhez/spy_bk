package models

import (
	"errors"
	"spyc_backend/app"

	"github.com/Reti500/mgomap"
)

// Color model
type Color struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name        string     `bson:"color"`
	Hexadecimal string     `bson:"hex"`
	Attachment  Attachment `bson:"attachment"`
}

// GetDocumentName ...
func (u *Color) GetDocumentName() string {
	return "colors"
}

// GetAll ...
func (u Color) GetAll() ([]Color, error) {
	var colors []Color

	if Model, ok := app.Mapper.GetModel(&u); ok {
		if err := Model.All().Exec(&colors); err != nil {
			return colors, err
		}

		return colors, nil
	}

	return colors, errors.New("Error to create color instance")
}
