package models

import (
	"spyc_backend/app"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Reti500/mgomap"
)

// Section model
type Section struct {
	Name       string        `json:"name" bson:"name"`
	CanAccess  bool          `json:"can_access" bson:"can_access"`
	TimeLess   time.Duration `json:"time_less" bson:"time_less"`
	LastAccess time.Time     `json:"last_access" bson:"last_access"`

	// AccessTime is in seconds
	AccessTime time.Duration `json:"access_time" bson:"access_time"`
}

// CommonSection document create a section for reply to all users
type CommonSection struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Name       string `json:"name" bson:"name"`
	AccessTime int    `json:"access_time" bson:"access_time"`
}

// GetDocumentName ...
func (m *CommonSection) GetDocumentName() string {
	return "common_sections"
}

// Replicate update all users with new section
func (m CommonSection) Replicate() error {
	if User, ok := app.Mapper.GetModel(&User{}); ok {
		var selector = bson.M{"sections.name": bson.M{"$in": []string{m.Name}}}
		var query = bson.M{"$pull": bson.M{"sections": bson.M{"name": m.Name}}}

		User.UpdateQuery(selector, query, true)

		var section = Section{
			CanAccess:  false,
			LastAccess: time.Now(),
			Name:       m.Name,
			TimeLess:   0,
			AccessTime: time.Duration(m.AccessTime) * time.Second,
		}

		selector = bson.M{"sections.name": bson.M{"$nin": []string{m.Name}}}
		query = bson.M{"$push": bson.M{"sections": section}}

		return User.UpdateQuery(selector, query, true)
	}

	return nil
}

// SetSections ....
func SetSections(u bson.ObjectId) error {
	if Common, ok := app.Mapper.GetModel(&CommonSection{}); ok {
		var sections []CommonSection

		if err := Common.All().Exec(&sections); err != nil {
			return err
		}

		for _, v := range sections {
			var section = Section{
				CanAccess:  false,
				LastAccess: time.Now(),
				Name:       v.Name,
				TimeLess:   0,
				AccessTime: time.Duration(v.AccessTime) * time.Second,
			}

			var selector = bson.M{"_id": u}
			var query = bson.M{"$push": bson.M{"sections": section}}

			if User, ok := app.Mapper.GetModel(&User{}); ok {
				User.UpdateQuery(selector, query, false)
			}
		}

	}

	return nil
}
