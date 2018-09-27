package models

import (
	"errors"
	"spyc_backend/app"

	"github.com/Reti500/mgomap"
	"gopkg.in/mgo.v2/bson"
)

// PINAdvertisement model
type PINAdvertisement struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Titles          Internationalized `bson:"titles"`
	Descriptions    Internationalized `bson:"decriptions"`
	Attachments     Internationalized `bson:"attachments"`
	AdvertisementID string            `json:"advertisement_id" bson:"advertisement_id"`
	Link            string            `bson:"link"`
	Scope           string            `bson:"scope"`
	TargetMinAge    int               `json:"target_min_age" bson:"target_min_age"`
	TargetMaxAge    int               `json:"target_max_age" bson:"target_max_age"`
	TargetGender    []string          `json:"target_gender" bson:"target_gender"`
	Tags            []string          `bson:"tags"`
	Geolocation     Geo               `bson:"geolocation"`
	Views           []Views           `bson:"views"`
}

// GetDocumentName ...
func (u *PINAdvertisement) GetDocumentName() string {
	return "pin_advertisement"
}

// GetPinsForUser returns all pins associated to user
func (u PINAdvertisement) GetPinsForUser(user User) ([]PINAdvertisement, error) {
	var pins []PINAdvertisement

	if Pin, ok := app.Mapper.GetModel(&u); ok {

		var match = bson.M{
			"target_min_age": bson.M{"$lte": user.GetAge()},
			"target_max_age": bson.M{"$gte": user.GetAge()},
			// "target_gender":  bson.M{"$in": []string{c.CurrentUser.PersonalData.Gender}},
		}

		var pipe = mgomap.Aggregate{}.Match(match)

		if err := Pin.Pipe(pipe, &pins); err != nil {
			return pins, err
		}

		return pins, nil
	}

	return pins, errors.New("Error to create Pin ref")
}

// React ...
// @react can be ["click" or "view"]
func (u PINAdvertisement) React(react, pinID, userID string) error {
	if Pin, ok := app.Mapper.GetModel(&u); ok {
		var view = Views{
			Action:    react,
			Completed: true,
			Through:   "PIN",
			UserID:    userID,
			Total:     1,
		}

		var selector = bson.M{
			"_id":   bson.ObjectIdHex(pinID),
			"views": bson.M{"$elemMatch": bson.M{"action": react, "user_id": userID}}}

		total, err := Pin.Query(selector).Count()
		if err != nil {
			return err
		}

		if total > 0 {
			// React exists then increase total
			var query = bson.M{"$inc": bson.M{"views.$.total": 1}}

			return Pin.UpdateQuery(selector, query, false)
		}

		selector = bson.M{"_id": bson.ObjectIdHex(pinID)}
		var query = bson.M{"$push": bson.M{"views": view}}

		return Pin.UpdateQuery(selector, query, false)
	}

	return errors.New("Error to create Pin ref")
}
