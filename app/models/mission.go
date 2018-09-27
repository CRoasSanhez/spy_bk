package models

import (
	"errors"
	"time"

	"github.com/revel/revel"

	"spyc_backend/app"
	"spyc_backend/app/core"

	"github.com/Reti500/mgomap"
	"gopkg.in/mgo.v2/bson"
)

// Mission struct for DB
type Mission struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	// Used for internationalization and regresion compatibility
	Title       Internationalized `json:"title" bson:"title"`
	Description Internationalized `json:"description" bson:"description"`

	Type          string        `json:"type" bson:"type" validate:"nonzero,regexp=^[a-zA-Z0-9 ]*$"`
	Advertisement string        `json:"advertisement" bson:"advertisement"`
	Attachment    Attachment    `json:"-" bson:"attachment"`
	Geolocation   *Geo          `json:"geolocation" bson:"geolocation"`
	StartDate     time.Time     `json:"start_date" bson:"start_date" validate:"nonzero"`
	EndDate       time.Time     `json:"end_date" bson:"end_date" validate:"nonzero"`
	Priority      int           `json:"priority" bson:"priority" validate:"nonzero,regexp=^[a-z]*$"`
	Countries     []string      `json:"countries" bson:"countries"`
	Campaign      bson.ObjectId `json:"campaign" bson:"campaign" validate:"nonzero"`
	AdViews       []Views       `json:"-" bson:"ad_views"`

	// Frequency indicates how often a Mission will be activated in a day
	Frequency int `json:"frequency" bson:"frequency"`

	// Shows to user if its subscribed or not
	Subscribed bool `json:"isSubscribed" bson:"-"`

	// Shows to user if its completed or not
	Completed bool `json:"is Completed" bson:"-"`

	// Show the currentt target
	Targets []Target `json:"steps" bson:"targets"`
	Periods []Period `json:"periods" bson:"periods"`

	// For agregation pipeline -- Used only for retreiving data
	// Empty in collection
	Subscriptions []MissionUser `json:"-" bson:"subscriptions"`
	Users         []User        `json:"-" bson:"users"`
	DateDiff      int64         `json:"-" bson:"date_diff"`

	// Used for update Mission Language
	Langs []string `json:"-" bson:"langs"`
}

// Missions struct for GetSubscribed Missions
type Missions struct {
	Missions []Mission `json:"games"`
}

// CashHuntResponse is the response when users selects
// Cash hunt
type CashHuntResponse struct {
	Missions interface{} `json:"games"`
	Rewards  interface{} `json:"rewards"`
}

// CashHuntRequest used for Dashboard
type CashHuntRequest struct {
	Mission Mission  `json:"mission"`
	Targets []Target `json:"targets"`
}

// GetDocumentName needed function for Mongo storage with mgomap
func (m *Mission) GetDocumentName() string {
	return "missions"
}

// SetStatus change Mission status according to constants
func (m *Mission) SetStatus(status string) {
	m.Status.Code = core.SubscriptionStatus[status]
	m.Status.Name = status
}

// FindTargetsByMission returns active targets according
// to the given MissionID
func (m *Mission) FindTargetsByMission(missionID string) (targets []Target, err error) {

	Target, _ := app.Mapper.GetModel(&Target{})

	err = Target.Query(bson.M{
		"status.name": core.StatusActive,
		"missions":    bson.ObjectIdHex(missionID),
	}).Exec(&targets)

	if err != nil {
		return targets, err
	}
	return targets, err
}

// GetMissionTarget gets current Mission active Target
func (m *Mission) GetMissionTarget() (targets []Target) {

	Target, _ := app.Mapper.GetModel(&Target{})

	_ = Target.Query(bson.M{
		"status.name":  core.StatusActive,
		"missions.$id": m.GetID().Hex(),
	}).Exec(&targets)

	return
}

// FindMissionsByCoords retreives a list of Missions based on the
// given coordinates, status active and within the MaxDistanceForPin constant
func (m *Mission) FindMissionsByCoords(long, lat float64) (missions []Mission, err error) {

	Mission, _ := app.Mapper.GetModel(&Mission{})
	err = Mission.Query(
		bson.M{
			"status.name": core.StatusActive,
			"geolocation": bson.M{
				"$nearSphere": bson.M{
					"$geometry": bson.M{
						"type":        "Point",
						"coordinates": []float64{long, lat},
					},
					"$minDistance": 0,
					"$maxDistance": core.MaxDistanceForPin,
				},
			},
		}).Exec(&missions)
	if err != nil {
		return missions, err
	}
	return missions, nil
}

// GetAllTargetsByMission returns all active Targets of the given Mission
func (m *Mission) GetAllTargetsByMission(missionID string) (targets []Target, err error) {
	Target, _ := app.Mapper.GetModel(&Target{})

	var query = []bson.M{
		bson.M{"status.name": core.StatusActive},
		bson.M{"missions": bson.ObjectIdHex(missionID)},
	}

	if err = Target.FindWithOperator("$and", query).Exec(&targets); err != nil {
		return targets, err
	}
	return
}

// Complete changes Mission status to complete
func (m *Mission) Complete() bool {
	Mission, _ := app.Mapper.GetModel(&Mission{})
	var status = core.StatusCompleted

	var selector = bson.M{
		"_id":         m.GetID(),
		"status.name": core.StatusActive,
	}
	endDate, _ := time.Parse(core.MXTimeFormat, time.Now().Format(core.MXTimeFormat))
	var query = bson.M{"$set": bson.M{
		"completed":   true,
		"status.name": status,
		"status.code": core.SubscriptionStatus[status],
		"end_date":    endDate,
	}}

	if err := Mission.UpdateQuery(selector, query, false); err != nil {
		revel.ERROR.Printf("ERROR UPDATE Mission %s --- %s", m.GetID().Hex(), err.Error())
		return false
	}

	core.Notify(core.InfoEntry, "Cash Hunt was Deactivated", "A Cash Hunt was Dectivated", m.Title.GetString(m.Langs[0]),
		core.GetDashboardPath()+"missions/"+m.GetID().Hex(), "", m.Description.GetString(m.Langs[0]), []interface{}{
			bson.M{"title": m.Type, "value": core.ConcatArray(m.Countries), "short": true},
		})

	return true
}

// FindTargetByID find in missions targets array the element by the given Target ObjectID
func (m *Mission) FindTargetByID(id bson.ObjectId) (target Target) {
	for i := 0; i < len(m.Targets); i++ {
		if m.Targets[i].GetID() == id {
			return m.Targets[i]
		}
	}
	return target
}

// FindMissionAndSubscriptions ...
func (m *Mission) FindMissionAndSubscriptions(missionID, missionStatus string, userID bson.ObjectId) (mission Mission, err error) {
	if Mission, ok := app.Mapper.GetModel(&mission); ok {
		var match1 = bson.M{"$and": []bson.M{
			bson.M{"_id": bson.ObjectIdHex(missionID)},
			bson.M{"status.name": missionStatus},
		}}
		var unwind = bson.M{"$unwind": "$subscription"}
		var match2 = bson.M{"$and": []bson.M{
			bson.M{"subscription.user_id": userID},
			bson.M{"subscription.status.name": core.StatusInit},
		}}
		var project = bson.M{"$project": bson.M{"_id": 1, "status": 1, "title": 1, "desccription": 1, "subscriptions": []string{"$subscription"}}}
		var pipe = mgomap.Aggregate{}.Match(match1).LookUp("missions_users", "_id", "mission_id", "subscription").Add(unwind).Match(match2).Add(project)

		if err = Mission.Pipe(pipe, &mission); err != nil {
			revel.ERROR.Printf("ERROR PIPE Mission: %s --- %s", missionID, err.Error())
			return mission, err
		}
		return mission, nil
	}
	return mission, errors.New("Error mapping Mission")
}

// FindUserInViews ...
func (m *Mission) FindUserInViews(userID string) Views {
	for i := 0; i < len(m.AdViews); i++ {
		if m.AdViews[i].UserID == userID {
			return m.AdViews[i]
		}
	}
	return Views{}
}
