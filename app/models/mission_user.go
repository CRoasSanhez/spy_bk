package models

import (
	"errors"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"time"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// MissionUser model
type MissionUser struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	MissionID     bson.ObjectId `json:"game_id" bson:"mission_id" validate:"nonzero"`
	UserID        bson.ObjectId `json:"user_id" bson:"user_id" validate:"nonzero"`
	StartDate     time.Time     `json:"start_date" bson:"start_date"`
	EndDate       time.Time     `json:"end_date" bson:"end_date"`
	CurrentTarget bson.ObjectId `json:"-" bson:"current_target"`

	// Shows if an user completes the mission
	Completed bool `json:"-" bson:"completed" default:"false"`
}

// GetDocumentName ...
func (m *MissionUser) GetDocumentName() string {
	return "missions_users"
}

// GetActiveSubscriptionsByUser ...
func (m *MissionUser) GetActiveSubscriptionsByUser(userID bson.ObjectId) (subscriptions []MissionUser) {

	MissionUser, _ := app.Mapper.GetModel(&MissionUser{})

	err := MissionUser.Query(
		bson.M{
			"$and": []bson.M{
				bson.M{"status.name": core.StatusInit},
				bson.M{"user_id": userID},
			},
		}).Exec(&subscriptions)
	if err != nil {
		revel.ERROR.Print(err)
	}
	return subscriptions
}

// FindSubscriptionsByUserStatus retrieves user subscriptions to missions
// With status incomplete
func (m *MissionUser) FindSubscriptionsByUserStatus(userID bson.ObjectId) (subscriptions []MissionUser) {
	MissionUser, _ := app.Mapper.GetModel(&MissionUser{})

	statuses := []string{core.StatusInit, core.StatusActive, core.StatusPendingValidation, core.StatusCompleted, core.StatusCompletedNotWinner, core.StatusCompletedWinner}

	match := bson.M{"$and": []bson.M{
		bson.M{"user_id": userID},
		bson.M{"status.name": bson.M{"$in": statuses}},
	},
	}

	pipe := mgomap.Aggregate{}.Match(match)

	err := MissionUser.Pipe(pipe, &subscriptions)
	if err != nil {
		revel.ERROR.Print(err)
	}
	return subscriptions
}

// SetStatus ...
func (m *MissionUser) SetStatus(status string) {
	m.Status.Code = core.SubscriptionStatus[status]
	m.Status.Name = status
}

// CompleteSubscription for the subscription
func (m *MissionUser) CompleteSubscription(status string) error {
	MissionUser, _ := app.Mapper.GetModel(&MissionUser{})
	m.SetStatus(status)
	m.Completed = true
	endDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
	m.EndDate = endDate
	err := MissionUser.Update(m)
	if err != nil {
		revel.ERROR.Printf("ERROR UPDATE Mission Subscription %s --- %s", m.GetID().Hex(), err.Error())
		return err
	}
	return nil
}

// CompleteUnsubscription for the subscription
func (m *MissionUser) CompleteUnsubscription() error {
	MissionUser, _ := app.Mapper.GetModel(&MissionUser{})
	m.SetStatus(core.StatusIncomplete)
	endDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
	m.EndDate = endDate
	err := MissionUser.Update(m)
	if err != nil {
		revel.ERROR.Printf("ERROR COMPLETE Mission UnSubscription %s --- %s", m.GetID().Hex(), err.Error())
		return err
	}
	return nil
}

// CompleteSubscriptionByMissionUser change mission Subscription
func (m *MissionUser) CompleteSubscriptionByMissionUser(missionID, userID, status string) error {

	subscription := MissionUser{}
	if MissionUser, ok := app.Mapper.GetModel(&subscription); ok {
		var selector = bson.M{
			"user_id":     bson.ObjectIdHex(userID),
			"mission_id":  bson.ObjectIdHex(missionID),
			"status.name": bson.M{"$in": []string{core.StatusInit, core.StatusPendingValidation}},
		}
		endDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
		var query = bson.M{"$set": bson.M{
			"completed":   true,
			"status.name": status,
			"status.code": core.SubscriptionStatus[status],
			"end_date":    endDate,
		}}

		if err := MissionUser.UpdateQuery(selector, query, false); err != nil {
			revel.ERROR.Printf("ERROR UPDATE MissionUser Mission: %s --- %s", missionID, err.Error())
			return err
		}
		return nil
	}
	return errors.New("Error mapping MissionUser model")
}

// GetSubscriptionByMissionUser returns the user subscription to the mission
func (m *MissionUser) GetSubscriptionByMissionUser(missionID string, userID bson.ObjectId) error {
	MissionUser, _ := app.Mapper.GetModel(&MissionUser{})

	err := MissionUser.Query(
		bson.M{
			"$and": []bson.M{
				bson.M{"status.name": core.StatusInit},
				bson.M{"user_id": userID},
				bson.M{"mission_id": bson.ObjectIdHex(missionID)},
			},
		}).Exec(&m)
	if err != nil {
		return err
	}
	return nil
}

// FindSubscriptionsByMission return all currently subscriptions for the fiven Mission
func (m *MissionUser) FindSubscriptionsByMission(missionID string) (subscriptions []MissionUser) {

	MissionUser, _ := app.Mapper.GetModel(&MissionUser{})
	err := MissionUser.Query(bson.M{
		"$and": []bson.M{
			bson.M{"mission_id": bson.ObjectIdHex(missionID)},
			bson.M{"status.name": core.StatusInit},
		},
	}).Exec(&subscriptions)
	if err != nil {
		revel.ERROR.Println(err.Error())
	}
	return subscriptions
}

// CompleteAllSubscriptions is recomended to run this on gorutine
func (m MissionUser) CompleteAllSubscriptions(missionID string) {

	subscription := MissionUser{}
	if MissionUser, ok := app.Mapper.GetModel(&subscription); ok {
		var selector = bson.M{"$and": []bson.M{
			bson.M{"mission_id": bson.ObjectIdHex(missionID)},
			bson.M{"completed": false},
		}}
		endDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
		var query = bson.M{"$set": bson.M{
			"completed":   true,
			"status.name": core.StatusIncomplete,
			"status.code": core.SubscriptionStatus[core.StatusIncomplete],
			"end_date":    endDate,
		}}

		if err := MissionUser.UpdateQuery(selector, query, true); err != nil {
			revel.ERROR.Printf("ERROR UPDATE MissionUser Mission: %s --- %s", missionID, err.Error())
		}
	}
	/*
		subscriptions := []MissionUser{}
		MissionUser, _ := app.Mapper.GetModel(&MissionUser{})

		MissionUser.Query(bson.M{
			"$and": []bson.M{
				bson.M{"mission_id": bson.ObjectIdHex(missionID)},
				bson.M{"completed": false},
			},
		}).Exec(subscriptions)

		for i := 0; i < len(subscriptions); i++ {
			subscriptions[i].SetStatus(core.StatusIncomplete)
			subscriptions[i].Completed = true
			err := MissionUser.Update(&subscriptions[i])
			if err != nil {
				revel.ERROR.Print(err)
				continue
			}
		}
	*/
}
