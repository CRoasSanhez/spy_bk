package models

import (
	"errors"
	"spyc_backend/app"

	"github.com/Reti500/mgomap"
)

// Anim model
type Anim struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name       string     `bson:"name"`
	Attachment Attachment `bson:"attachment"`
}

// GetDocumentName ...
func (u *Anim) GetDocumentName() string {
	return "animations"
}

// GetAll ...
func (u Anim) GetAll() ([]Anim, error) {
	var anims []Anim

	if Model, ok := app.Mapper.GetModel(&u); ok {
		if err := Model.All().Exec(&anims); err != nil {
			return anims, err
		}

		return anims, nil
	}

	return anims, errors.New("Error to create color instance")
}
