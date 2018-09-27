package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Address model
type Address struct {
	Street  string `json:"street,omitempty" bson:"street" validate:"regexp=^[a-zA-Z0-9,. ]*$"`
	City    string `json:"city,omitempty" bson:"city" validate:"regexp=^[a-zA-Z0-9., ]*$"`
	State   string `json:"state,omitempty" bson:"state" validate:"regexp=^[a-zA-Z0-9., ]*$"`
	Country string `json:"country,omitempty" bson:"country" validate:"regexp=^[a-zA-Z0. ]*$"`
	ZipCode string `json:"zip_code,omitempty" bson:"zip_code" validate:"regexp=^[0-9]{1,25}$"`
}

// Validate ...
func (m Address) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Street, validation.Match(regexp.MustCompile("^[a-zA-Z0-9,. ]*$"))),
		validation.Field(&m.City, validation.Match(regexp.MustCompile("^[a-zA-Z0-9., ]*$"))),
		validation.Field(&m.State, validation.Match(regexp.MustCompile("^[a-zA-Z0-9., ]*$"))),
		validation.Field(&m.Country, validation.Match(regexp.MustCompile("^[a-zA-Z0-9., ]*$"))),
		validation.Field(&m.ZipCode, validation.Match(regexp.MustCompile("^[0-9]{1,10}$"))),
	)
}
