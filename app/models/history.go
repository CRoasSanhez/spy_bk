package models

import (
	"errors"
	"mime/multipart"
	"spyc_backend/app"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
)

// History model
type History struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Attachment Attachment `bson:"attachment"`
	Views      []Views    `bson:"views"`
	Anim       string     `bson:"anim"`
	Color      string     `bson:"color"`
	Owner      string     `bson:"owner"`
	Text       string     `bson:"text"`
	Type       string     `bson:"type"`
	Reacts     []string   `bson:"reacts"`
	Users      []string   `bson:"users"`
	Public     bool       `bson:"public"`
	Seen       bool       `bson:"seen"`
}

type GHistories struct {
	UserID    string    `bson:"_id"`
	Histories []History `bson:"histories"`
	TotalSeen int       `bson:"total_seen"`
	TotalDocs int       `bson:"total_docs"`
	AVGSeen   float32   `bson:"avgSeen"`
}

// GetDocumentName ...
func (u *History) GetDocumentName() string {
	return "histories"
}

// Save ...
func (u *History) Save() error {
	if History, ok := app.Mapper.GetModel(&History{}); ok {
		if err := History.Create(u); err != nil {
			return err
		}

		return nil
	}

	return errors.New("Error to create History instance")
}

// FindHistories returns histories of users list
func (u History) FindHistories(users []string, userID string) ([]GHistories, error) {
	var result []GHistories

	if History, ok := app.Mapper.GetModel(&u); ok {
		var match = bson.M{"owner": bson.M{"$in": users}}
		var addFields = bson.M{"has_on_users": bson.M{
			"$cond": bson.M{"if": bson.M{"$in": []string{userID, "$users"}}, "then": true, "else": false}}}
		var addFields2 = bson.M{"seen": bson.M{
			"$cond": bson.M{"if": bson.M{"$in": []string{userID, "$reacts"}}, "then": true, "else": false}}}
		var match2 = bson.M{"$or": []bson.M{
			{"has_on_users": bson.M{"$eq": true}},
			{"public": bson.M{"$eq": true}},
		}}
		var group = bson.M{"$group": bson.M{
			"_id":        "$owner",
			"histories":  bson.M{"$push": "$$ROOT"},
			"total_seen": bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$seen", true}}, 1, 0}}},
			"total_docs": bson.M{"$sum": 1}}}

		// {$addFields: {"GG": {$avg: {$multiply: [{$divide: [100, "$total_docs"]}, "$total_seen"]}}} },
		var avgField = bson.M{"$addFields": bson.M{"avgSeen": bson.M{"$avg": bson.M{"$multiply": []interface{}{
			bson.M{"$divide": []interface{}{100, "$total_docs"}}, "$total_seen"}}}}}
		var sort = bson.M{"$sort": bson.M{"avgSeen": 1}}

		var pipe = mgomap.Aggregate{}.Match(match).Add(bson.M{"$addFields": addFields})
		pipe = pipe.Add(bson.M{"$addFields": addFields2}).Match(match2).Add(group).Add(avgField).Add(sort)

		err := History.Pipe(pipe, &result)

		return result, err
	}

	return result, errors.New("Error to create History instance")
}

// CountView ...
func (u History) CountView(completed bool, closeTime int, historyID, userID string) error {
	if History, ok := app.Mapper.GetModel(&u); ok {
		var view = Views{
			Action:    "view",
			CloseTime: closeTime,
			Completed: completed,
			Through:   "History",
			UserID:    userID,
			Total:     1,
		}

		var selector = bson.M{
			"_id":   bson.ObjectIdHex(historyID),
			"views": bson.M{"$elemMatch": bson.M{"action": "view", "user_id": userID}}}

		total, err := History.Query(selector).Count()
		if err != nil {
			return err
		}

		if total > 0 {
			// React exists then increase total
			var query = bson.M{"$inc": bson.M{"views.$.total": 1}}

			return History.UpdateQuery(selector, query, false)
		}

		selector = bson.M{"_id": bson.ObjectIdHex(historyID)}
		var query = bson.M{"$push": bson.M{"views": view}, "$addToSet": bson.M{"reacts": userID}}

		return History.UpdateQuery(selector, query, false)
	}

	return errors.New("Error to create Pin ref")
}

// UpdateAttachment ...
func (u *History) UpdateAttachment(user User, part *multipart.FileHeader) error {
	if part == nil {
		return errors.New("File not found")
	}

	if u.Attachment.PATH != "" {
		u.Attachment.Remove()
	}

	// owner := &mgo.DBRef{
	// 	Id:         user.GetID(),
	// 	Collection: user.GetDocumentName(),
	// 	Database:   app.Mapper.DatabaseName,
	// }

	if err := u.Attachment.Init(AsDocumentBase(&user), part); err != nil {
		return err
	}

	if err := u.Attachment.Upload(); err != nil {
		return err
	}

	if History, ok := app.Mapper.GetModel(u); ok {
		var selector = bson.M{"_id": u.GetID()}
		var query = bson.M{"$set": bson.M{"attachment": u.Attachment}}

		if err := History.UpdateQuery(selector, query, false); err != nil {
			return err
		}
	}

	return nil
}
