package models

import (
	"errors"
	"spyc_backend/app"
	"time"

	"github.com/revel/revel"

	"github.com/Reti500/mgomap"
	"gopkg.in/mgo.v2/bson"
)

// StatsSection embedded model to section stats
type StatsSection struct {
	Created          int            `json:"created" bson:"created"`
	Subscribed       int            `json:"subscribed" bson:"subscribed"`
	Total            int            `json:"total" bson:"total"`
	Completed        int            `json:"completed" bson:"completed"`
	Uncompleted      int            `json:"uncompleted" bson:"uncompleted"`
	LastResult       string         `json:"last_result" bson:"last_result"`
	FirstAccessDate  time.Time      `json:"first_access_date" bson:"first_access_date"`
	LastAccess       time.Time      `json:"last_access" bson:"last_access"`
	LastCompleted    time.Time      `json:"last_completed" bson:"last_completed"`
	LastUncompleted  time.Time      `json:"last_uncompleted" bson:"last_uncompleted"`
	CreatedTypes     map[string]int `json:"created_types" bson:"created_types"`
	SubscribedTypes  map[string]int `json:"subscribed_types" bson:"subscribed_types"`
	CreateScores     map[string]int `json:"create_scores" bson:"create_scores"`
	SubscribedScores map[string]int `json:"subscribed_scores" bson:"scosubscribed_scoresres"`
}

// StatsAccount embedded model to account stats
type StatsAccount struct {
	IsValidated       bool      `json:"is_validated" bson:"is_validated"`
	FriendRequests    int       `json:"friend_requests" bson:"friend_requests"`
	SMSSent           int       `json:"sms_sent" bson:"sms_sent"`
	CratedAt          time.Time `json:"created_at" bson:"created_at"`
	LastFriendRequest time.Time `json:"last_friend_request" bson:"last_friend_request"`
	LastProfileUpdate time.Time `json:"last_profile_update" bson:"last_profile_update"`
	LastSMSSent       time.Time `json:"first_sms_sent" bson:"first_sms_sent"`
	LastUpdated       time.Time `json:"last_updated" bson:"last_updated"`
}

// Stats model
type Stats struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Owner      bson.ObjectId `json:"owner" bson:"owner"`
	Missions   StatsSection  `json:"missions" bson:"missions"`
	Challenges StatsSection  `json:"challenges" bson:"challenges"`
	Account    StatsAccount  `json:"account" json:"account"`
}

// GetDocumentName ...
func (m *Stats) GetDocumentName() string {
	return "stats"
}

// CreateStats create a stats for user
func CreateStats(owner bson.ObjectId) error {
	var account = StatsAccount{CratedAt: time.Now()}
	var stats = Stats{Owner: owner, Account: account}

	if Stats, ok := app.Mapper.GetModel(&stats); ok {
		if err := Stats.Create(&stats); err != nil {
			revel.ERROR.Print(err)
			return err
		}
	}

	return nil
}

// StatsIncFields update fields
// Field have to be {field: inc}
// {"total": 1, "completed": 1}
func StatsIncFields(owner bson.ObjectId, fields interface{}) error {
	if Stats, ok := app.Mapper.GetModel(&Stats{}); ok {
		var selector = bson.M{"owner": owner}
		var query = bson.M{"$inc": fields}

		return Stats.UpdateQuery(selector, query, false)
	}

	return errors.New("Error to parse Stats document")
}

// StatsNowFields update fields with time.now()
// fields have to be {"first", "last"}
func StatsNowFields(owner bson.ObjectId, fields []string) error {
	if len(fields) <= 0 {
		return errors.New("There are no fields to update")
	}

	var f = bson.M{}
	var now = time.Now()
	for _, v := range fields {
		f[v] = now
	}

	if Stats, ok := app.Mapper.GetModel(&Stats{}); ok {
		var selector = bson.M{"owner": owner}
		var query = bson.M{"$set": f}

		return Stats.UpdateQuery(selector, query, false)
	}

	return errors.New("Error to parse Stats document")
}

// StatsSetField update last result field
// Field have to be {field: inc}
// {"section.property": "value"}
func StatsSetField(owner bson.ObjectId, field interface{}) error {
	if Stats, ok := app.Mapper.GetModel(&Stats{}); ok {
		var selector = bson.M{"owner": owner}
		var query = bson.M{"$set": field}

		return Stats.UpdateQuery(selector, query, false)
	}

	return errors.New("Error to parse Stats document")
}

// StatsUpdateFields ...
// incFields have to be {"section.property": 1, "challenges.completed": 1}
// nowFields have to be {"challenges.last_access", "challenges.last_completed"}
func StatsUpdateFields(incFields interface{}, nowFields interface{}, owner bson.ObjectId) (err error) {
	if incFields != nil {
		err = StatsIncFields(owner, incFields)
	}
	if nowFields != nil {
		if timeFields, ok := nowFields.([]string); ok {
			err = StatsNowFields(owner, timeFields)
		} else {
			err = errors.New("Failed parse nowFields to []string")
		}
	}
	if err != nil {
		revel.ERROR.Printf("ERROR UPDATING Stats field for user: %s --- %s", owner.Hex(), err.Error())
		return err
	}
	return nil
}

// ReplicateStats delete stats and create new for all users
func ReplicateStats() error {
	if U, ok := app.Mapper.GetModel(&User{}); ok {
		var users []User
		var lookup = bson.M{"from": "stats", "localField": "_id", "foreignField": "owner", "as": "stats"}
		var project = bson.M{"count": bson.M{"$size": "$stats"}}
		var match = bson.M{"count": bson.M{"$lte": 0}}

		var pipe = mgomap.Aggregate{}.Add(bson.M{"$lookup": lookup}).Add(bson.M{"$project": project}).Add(bson.M{"$match": match})

		if err := U.Pipe(pipe, &users); err != nil {
			return err
		}

		for _, v := range users {
			CreateStats(v.GetID())
		}

	}

	return errors.New("Err to create User instance")
}
