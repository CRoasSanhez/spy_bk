package models

import "github.com/Reti500/mgomap"

// People ...
type People struct {
	People []string `json:"people"`
}

// Person ...
type Person struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Number    string     `json:"number"`
	Gender    string     `json:"gender"`
	PayPal    string     `json:"paypal"`
	BirthDate CustomTime `json:"birth_date"`
	Coupon    string     `json:"-"`
}

// GetDocumentName ...
func (m Person) GetDocumentName() string {
	return "people"
}
