package models

import (
	"log"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"time"

	"github.com/Reti500/mgomap"
	"gopkg.in/mgo.v2/bson"
)

// Campaign model
type Campaign struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name      string        `json:"name" bson:"name" validation:"nonzero,min=6,max=120,regexp=^[a-zA-Z0-9 ]*$"`
	Budget    int           `json:"budget" bson:"budget" validate:"regexp=^[a-zA-Z0-9 ]*$"`
	StartDate time.Time     `json:"start_date" bson:"start_date" validate:"nonzero"`
	EndDate   time.Time     `json:"end_date" bson:"end_date" validate:"nonzero"`
	Sponsor   bson.ObjectId `json:"sponsor" bson:"sponsor" validate:"nonzero"`
}

// GetDocumentName needed function for Mongo storage with mgomap
func (m *Campaign) GetDocumentName() string {
	return "campaigns"
}

// SetStatus change Reward status according to constants
func (m *Campaign) SetStatus(status string) {
	m.Status.Code = core.AccountStatus[status]
	m.Status.Name = status
}

// GetAllCampaignsBySponsor returns all the campaigns of the sponsor
func (m *Sponsor) GetAllCampaignsBySponsor() (campaigns []Campaign, err error) {
	Campaign, _ := app.Mapper.GetModel(&Campaign{})

	err = Campaign.Query(bson.M{
		"sponsor": bson.ObjectId(m.GetID()),
	}).Exec(&campaigns)

	if err != nil {
		return campaigns, err
	}
	return
}

// GetCampaignByID returns a campaign by thegiven campaignID type bson.ObjectId
// Used to get Sponsor in Sponsor model
func (m *Campaign) GetCampaignByID(campaignID bson.ObjectId) (campaign *Campaign) {
	camp := Campaign{}
	Campaign, _ := app.Mapper.GetModel(&camp)

	err := Campaign.FindBy("_id", campaignID).Exec(&camp)
	if err != nil {
		log.Println(err.Error())
	}
	return &camp

}
