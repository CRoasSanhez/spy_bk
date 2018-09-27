package models

import (
	"spyc_backend/app"
	"spyc_backend/app/core"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// Win ...
type Win struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	UserID       bson.ObjectId `json:"user_id" bson:"user_id"`
	ResourceID   bson.ObjectId `json:"resource_id" bson:"resource_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	Score        string        `json:"score" bson:"score"`
}

// GetDocumentName ...
func (m *Win) GetDocumentName() string {
	return "wins"
}

// SetStatus change Reward status according to constants
func (m *Win) SetStatus(status string) {
	m.Status.Code = core.SubscriptionStatus[status]
	m.Status.Name = status
}

// CreateWin inserts a document into Winners collection
// Set this as gorutine
func (m *Win) CreateWin(resourceID, userID bson.ObjectId,
	winType, score string) (win Win, err error) {

	Win, _ := app.Mapper.GetModel(&win)

	win.UserID = userID
	win.ResourceID = resourceID
	win.ResourceType = winType
	win.Score = score
	win.Status.Name = core.StatusPendingCollection

	err = Win.Create(&win)
	if err != nil {
		return win, err
	}
	return win, nil
}

// GetWinsByUser get Win documents for the given user
func (m *Win) GetWinsByUser(userID string) (wins []Win) {

	Win, _ := app.Mapper.GetModel(&Win{})
	Win.Query(
		bson.M{
			"$and": []bson.M{
				bson.M{"status.name": core.StatusInit},
				bson.M{"user_id": bson.ObjectIdHex(userID)},
			},
		}).Exec(&wins)
	return wins
}

// GetByResource get the wint based on the given user and resource
func (m *Win) GetByResource(resourceID, userID string) error {

	Win, _ := app.Mapper.GetModel(&Win{})
	err := Win.Query(
		bson.M{
			"$and": []bson.M{
				bson.M{"resource_id": bson.ObjectIdHex(resourceID)},
				bson.M{"user_id": bson.ObjectIdHex(userID)},
				bson.M{"status.name": core.StatusInit},
			},
		}).Exec(&m)
	if err != nil {
		revel.ERROR.Print(err.Error())
		return err
	}
	return nil
}
