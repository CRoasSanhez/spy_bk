package models

import (
	"errors"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"time"

	"github.com/revel/revel"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
)

// Target structfor DB
type Target struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name         Internationalized `json:"name" bson:"name"`
	Description  Internationalized `json:"description" bson:"description"`
	StartDate    time.Time         `json:"start_date" bson:"start_date"`
	EndDate      time.Time         `json:"end_date" bson:"end_date"`
	Type         string            `json:"type" bson:"type" validate:"nonzero,max=50,regexp=^[a-zA-Z0-9 ]*$"`
	Order        int               `json:"order" bson:"order" validate:"min=0,max=10"`
	NextTargetID string            `json:"next_step_id" bson:"next_target_id"`
	Geolocation  *Geo              `json:"geolocation" bson:"geolocation"`
	WebURL       string            `json:"web_url" bson:"web_url" validate:"nonzero,min=6,regexp=(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`
	Missions     bson.ObjectId     `json:"games" bson:"missions" validate:"nonzero"`
	Questions    []string          `json:"-" bson:"questions"`
	Attachment   Attachment        `json:"-" bson:"attachement"`

	// Attribute send to client
	Question Question `json:"question" bson:"-"`

	// Attribute for target type TEXT/NIP
	Score string `json:"score" bson:"score"`

	Subscriptions []TargetUser `json:"subscriptions" bson:"subscriptions"`
	Rewards       []Reward     `json:"-" bson:"rewards"`
	Langs         []string     `json:"-" bson:"langs"`
}

// GetDocumentName needed function for Mongo storage with mgomap
func (m *Target) GetDocumentName() string {
	return "targets"
}

// SetStatus sets Target status according to constants
func (m *Target) SetStatus(status string) {
	m.Status.Code = core.SubscriptionStatus[status]
	m.Status.Name = status
}

// GetRandomQuestion ...
func (m *Target) GetRandomQuestion() (question Question) {

	lenQuestions := len(m.Questions)
	if lenQuestions == 0 {
		return //question
	}
	// We obtain the random question
	id := core.CustomRandomInt(0, lenQuestions)

	Question, _ := app.Mapper.GetModel(&question)

	err := Question.Find(m.Questions[id]).Exec(&question)
	if err != nil {
		revel.ERROR.Print("ERROR FIND Question for Target: " + m.GetID().Hex())
		return //question
	}
	return //question
}

// GetByOrderAndMission ...
func (m *Target) GetByOrderAndMission(order int, missionID string, withRewards bool) (err error) {

	Target, _ := app.Mapper.GetModel(&Target{})
	if withRewards {
		match := bson.M{
			"status.name":         core.StatusActive,
			"missions":            bson.ObjectIdHex(missionID),
			"rewards.status.name": core.StatusActive,
			"order":               order,
		}
		pipe := mgomap.Aggregate{}.LookUp("rewards", "_id", "resource_id", "rewards").Match(match)

		if err = Target.Pipe(pipe, m); err != nil {
			revel.ERROR.Printf("ERROR FIND Target in mission: %s --- %s", missionID, err.Error())
			return err
		}
	} else {
		err = Target.Query(bson.M{
			"order":       order,
			"missions":    bson.ObjectIdHex(missionID),
			"status.name": core.StatusActive,
		}).Exec(m)
		if err != nil {
			revel.ERROR.Printf("ERROR FIND Target in mission: %s --- %s", missionID, err.Error())
			return err
		}
	}
	return nil
}

// AddQuestion adds an ID to the questions slice for the calling Target
func (m *Target) AddQuestion(questionID string) (err error) {

	m.Questions = append(m.Questions, questionID)

	Target, _ := app.Mapper.GetModel(&Target{})
	err = Target.Update(m)
	if err != nil {
		return err
	}
	return nil
}

// FindTargetAndSubscriptions ...
func (m *Target) FindTargetAndSubscriptions(targetID, targetStatus string, userID bson.ObjectId) (target Target, err error) {
	if Target, ok := app.Mapper.GetModel(&target); ok {
		var match1 = bson.M{"$and": []bson.M{
			bson.M{"_id": bson.ObjectIdHex(targetID)},
			bson.M{"status.name": core.StatusActive},
		}}
		var unwind = bson.M{"$unwind": "$subscription"}
		var match2 = bson.M{"$and": []bson.M{
			bson.M{"subscription.user_id": userID},
			bson.M{"subscription.status.name": core.StatusInit},
		}}
		var project = bson.M{"$project": bson.M{"_id": 1, "next_target_id": 1, "geolocation": 1, "name": 1, "order": 1,
			"type": 1, "description": 1, "status": 1, "score": 1, "subscriptions": []string{"$subscription"},
		}}
		var pipe = mgomap.Aggregate{}.Match(match1).LookUp("targets_users", "_id", "target_id", "subscription").Add(unwind).Match(match2).Add(project)

		if err = Target.Pipe(pipe, &target); err != nil {
			revel.ERROR.Printf("ERROR PIPE Target: %s --- %s", targetID, err.Error())
			return target, err
		}
		return target, nil
	}
	return target, errors.New("Error mapping Target")
}
