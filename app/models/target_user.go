package models

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// TargetUser is the user subscription to target
type TargetUser struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	TargetID  bson.ObjectId `json:"step_id" bson:"target_id" validate:"nonzero"`
	UserID    bson.ObjectId `json:"user_id" bson:"user_id" validate:"nonzero"`
	RewardID  interface{}   `json:"-" bson:"reward_id"`
	StartDate time.Time     `json:"start_date" bson:"start_date"`
	EndDate   time.Time     `json:"end_date" bson:"end_date"`

	// Shows if an user completed the target
	Completed  bool   `json:"-" bson:"completed" default:"false"`
	Intents    int    `json:"-" bson:"tries" default:"0"`
	Token      string `json:"-" bson:"token"`
	QuestionID string `json:"question" bson:"question"`
	Score      string `json:"score" bson:"score"`

	// Used when a target is type selfie/photo and is for validation
	ResponseImage Attachment `json:"response_image"`
}

// GetDocumentName ...
func (m *TargetUser) GetDocumentName() string {
	return "targets_users"
}

// SetStatus sets TargetUser status according to given string
func (m *TargetUser) SetStatus(status string) {
	m.Status.Code = core.SubscriptionStatus[status]
	m.Status.Name = status
}

// CompleteSubscription for the subscription
func (m *TargetUser) CompleteSubscription(status string) error {

	TargetUser, _ := app.Mapper.GetModel(&TargetUser{})
	m.SetStatus(status)
	if status == core.StatusCompleted {
		m.Completed = true
	}

	endDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
	m.EndDate = endDate

	err := TargetUser.Update(m)
	if err != nil {
		return err
	}
	return nil
}

// GetSubscription  based on the current user and targetID
// Validates if its complete
func (m *TargetUser) GetSubscription(targetID string, userID bson.ObjectId, status string) error {

	TargetUser, _ := app.Mapper.GetModel(&TargetUser{})
	err := TargetUser.Query(bson.M{
		"$and": []bson.M{
			bson.M{"target_id": bson.ObjectIdHex(targetID)},
			bson.M{"user_id": userID},
			bson.M{"status.name": status},
		},
	}).Exec(m)
	if err != nil {
		return err
	}
	return nil
}

// AddIntent increments by one the intents of subscription
func (m *TargetUser) AddIntent() error {
	TargetUser, _ := app.Mapper.GetModel(&TargetUser{})
	m.Intents++
	err := TargetUser.Update(m)
	if err != nil {
		return err
	}
	return nil
}

// FindSubscriptionsByUser returns target subscriptions
// in status: init, active or pending validaation
func (m *TargetUser) FindSubscriptionsByUser(userID bson.ObjectId) (targetSubscriptions []TargetUser, err error) {

	TargetUser, _ := app.Mapper.GetModel(&TargetUser{})
	// Find active User target subscriptions
	err = TargetUser.Query(bson.M{
		"user_id":   userID,
		"completed": false,
		"status.name": bson.M{
			"$in": []string{core.StatusInit, core.StatusActive, core.StatusPendingValidation},
		},
	}).Exec(&targetSubscriptions)
	if err != nil {
		revel.ERROR.Print(err)
		return targetSubscriptions, err
	}
	return targetSubscriptions, nil
}

// CreateSubscription ...
func (m *TargetUser) CreateSubscription() (err error) {
	TargetUser, _ := app.Mapper.GetModel(&TargetUser{})
	err = TargetUser.Create(m)
	if err != nil {
		revel.ERROR.Printf("ERROR CREATE Target Subscription User: %s __ Target: %s --- %s", m.UserID.Hex(), m.TargetID.Hex(), err.Error())
		return
	}
	return
}

// FindTarget returns the target from the currrent TargetUser
func (m *TargetUser) FindTarget() (target Target) {
	Target, _ := app.Mapper.GetModel(&target)
	Target.FindBy("_id", m.TargetID).Exec(&target)
	return target
}

// FindUser returns the user from the currrent TargetUser
func (m *TargetUser) FindUser() (user User) {
	User, _ := app.Mapper.GetModel(&user)
	User.FindBy("_id", m.UserID).Exec(&user)
	return user
}
