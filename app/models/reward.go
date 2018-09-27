package models

import (
	"math"
	"spyc_backend/app"
	"spyc_backend/app/core"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// Reward ...
type Reward struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	ResourceID   bson.ObjectId     `json:"resource_id" bson:"resource_id"`
	ResourceType string            `json:"resource_type" bson:"resource_type"`
	Name         Internationalized `json:"name" bson:"name"`
	Description  Internationalized `json:"description" bson:"description"`
	CampaignID   interface{}       `json:"-" bson:"campaign_id"`
	UserID       string            `json:"-" bson:"user_id"`
	Users        []string          `json:"-" bson:"users"`
	Type         string            `json:"-" bson:"type"`
	Coins        int               `json:"coins" bson:"coins"`
	Files        []ChallengeFiles  `json:"files" bson:"files"`
	Views        int               `json:"-" bson:"views"`
	IsVisible    bool              `json:"-" bson:"is_visible" default:"true"`
	CanDownload  bool              `json:"-" bson:"can_download"`

	// Max winners if > 0 then allows a max number of winners
	// else allows no limit
	MaxWinners int `json:"max_winners" bson:"max_winners"`

	// Multi is the flag used to seach wether user_id field of users field
	// Used for regresion compatibility
	Multi bool `json:"multi" bson:"multi"`

	// Returns the number of participants in mission
	Subscriptors int        `json:"subscriptors" bson:"-"`
	Attachment   Attachment `json:"-" bson:"attachment"`
	Langs        []string   `json:"-" bson:"langs"`
}

// GetDocumentName ...
func (m *Reward) GetDocumentName() string {
	return "rewards"
}

/*
// GetRewradByID returns the reward that match the given ID
func (m *Rewards) GetRewradByID(id string)error{
	Reward,_:= app.Mapper.GetModel(&Reward)
	err := Reward.Find(id).Excec(m)
	if err != nil{
		revel.ERROR.Print(err)
		return err
	}
	return nil
}
*/

// GetRewardsByResource returns rewards by the given resource
// The resource could be Mission or step
func (m *Reward) GetRewardsByResource(resourceID string) (rewards []Reward) {
	Reward, _ := app.Mapper.GetModel(&Reward{})
	err := Reward.Query(
		bson.M{
			"$and": []bson.M{
				bson.M{"status.name": core.StatusActive},
				bson.M{"resource_id": bson.ObjectIdHex(resourceID)},
			},
		}).Exec(&rewards)
	if err != nil {
		revel.ERROR.Print(err)
	}
	return rewards
}

// GetRewardsByUserAndResourceType returns  the rewards based on
// the given resourceTypes (cash_hunt, chat_games)
func (m *Reward) GetRewardsByUserAndResourceType(userID string, resourceTypes []string, visivle bool, limit int) (rewards []Reward) {

	if Reward, ok := app.Mapper.GetModel(&Reward{}); ok {
		var params = []bson.M{}

		var userFinder = bson.M{"$or": []bson.M{
			bson.M{"user_id": userID},
			bson.M{"users": bson.M{"$elemMatch": bson.M{"$eq": userID}}},
		}}

		if len(resourceTypes) == 0 {
			params = []bson.M{
				userFinder,
				//bson.M{"status.name": core.StatusCompleted},
				bson.M{"is_visible": visivle},
			}
		} else {
			params = []bson.M{
				userFinder,
				//bson.M{"status.name": core.StatusCompleted},
				bson.M{"resource_type": bson.M{"$in": resourceTypes}},
				bson.M{"is_visible": visivle},
			}
		}

		var query = Reward.FindWithOperator("$and", params)

		if limit > 0 {
			query = query.Limit(limit)
		}

		// Get pending Rewards for the user
		if err := query.Sort([]string{"-updated_at"}).Exec(&rewards); err != nil {
			revel.ERROR.Printf("ERROR FIND Rewards for user: %s --- %s", userID, err.Error())
			return rewards
		}
		return rewards
	}
	revel.ERROR.Print("ERROR MAPPING Reward")
	return rewards
}

// CompleteReward completes a reward and saves in DB
func (m *Reward) CompleteReward(userID string) bool {

	Reward, _ := app.Mapper.GetModel(&Reward{})

	//if !m.Multi {
	//	m.UserID = userID
	//	m.Status.Name = core.StatusCompleted
	//} else {
	// If there're no items left for this reward
	if m.MaxWinners <= 0 || len(m.Users) < m.MaxWinners {
		m.Users = append(m.Users, userID)
	}
	if m.MaxWinners == len(m.Users) {
		m.Status.Name = core.StatusCompleted
	}
	//}

	if err := Reward.Update(m); err != nil {
		revel.ERROR.Printf("ERROR COMPLETE REWARD %s --- %s", m.GetID(), err.Error())
	}

	return true
}

// SetStatus change Reward status according to constants
func (m *Reward) SetStatus(status string) {
	m.Status.Code = core.SubscriptionStatus[status]
	m.Status.Name = status
}

// GetCoinsByMission ...
func (m *Reward) GetCoinsByMission(missionID string) {

	gu := MissionUser{}
	subs := gu.FindSubscriptionsByMission(missionID)
	m.Subscriptors = len(subs)

	// Return total subscriptors by 0.5 and round to floor
	m.Coins = int(math.Ceil(float64(m.Subscriptors/2)) + 1)
}
