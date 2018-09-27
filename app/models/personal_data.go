package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

// PersonalData ...
type PersonalData struct {
	FirstName string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Gender    string     `json:"gender,omitempty" bson:"gender,omitempty"`
	Lang      string     `json:"lang,omitempty" bson:"lang,omitempty"`
	BirthDate CustomTime `json:"birth_date,omitempty" bson:"birth_date,omitempty"`
	Address   Address    `json:"address,omitempty" bson:"address,omitempty"`
}

// Validate ...
func (m PersonalData) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Match(regexp.MustCompile("^[a-zA-Z0-9,. ]*$"))),
		validation.Field(&m.LastName, validation.Match(regexp.MustCompile("^[a-zA-Z0-9,. ]*$"))),
		validation.Field(&m.Lang, validation.Match(regexp.MustCompile("^[a-zA-Z,. ]*$"))),
		validation.Field(&m.Address),
	)
}
