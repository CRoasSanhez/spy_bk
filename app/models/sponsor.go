package models

import (
	"log"
	"regexp"
	"spyc_backend/app"
	"spyc_backend/app/core"

	"github.com/Reti500/mgomap"
	validation "github.com/go-ozzo/ozzo-validation"
)

// Sponsor ...
type Sponsor struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name       string     `json:"name" bson:"name" validate:"nonzero,max=60,regexp=^[a-zA-Z0-9 ]*$"`
	Attachment Attachment `json:"-" bson:"attachment"`

	// Field used to show campaigns from sponsor in Dashboard
	Campaigns []Campaign `json:"campaigns" bson:"-"`
}

// Validate ...
func (m Sponsor) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Length(1, 60), validation.Match(regexp.MustCompile("^[a-zA-Z0-9 ]*$"))),
	)
}

// GetDocumentName needed function for Mongo storage with mgomap
func (m *Sponsor) GetDocumentName() string {
	return "sponsors"
}

// SetStatus change Reward status according to constants
func (m *Sponsor) SetStatus(status string) {
	m.Status.Code = core.SubscriptionStatus[status]
	m.Status.Name = status
}

// GetSponsor gets the sponsor for the curret Campaign
func (m *Campaign) GetSponsor() (sponsor Sponsor) {

	Sponsor, _ := app.Mapper.GetModel(&Sponsor{})
	err := Sponsor.FindBy("_id", m.Sponsor).Exec(&sponsor)
	if err != nil {
		log.Println(err.Error())
	}
	return
}
