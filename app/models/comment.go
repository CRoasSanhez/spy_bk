package models

import (
	"github.com/Reti500/mgomap"
	"gopkg.in/mgo.v2/bson"
)

// React model
type React struct {
	User bson.ObjectId `bson:"user"`
	Name string        `json:"name" bson:"name"`
}

// LittleComment embedded model
type LittleComment struct {
	Owner      bson.ObjectId `json:"-" bson:"owner"`
	Message    string        `json:"message" bson:"message"`
	Attachment Attachment    `json:"-" bson:"attachment"`
	User       interface{}   `json:"owner" bson:"user"`
	Reactions  []React       `json:"-" bson:"reactions"`
}

// Comment model
type Comment struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Owner      bson.ObjectId   `json:"-" bson:"owner"`
	Message    string          `json:"message" bson:"message"`
	Target     string          `bson:"target"`
	TargetID   string          `json:"target_id"`
	Attachment Attachment      `json:"-" bson:"attachment"`
	Comments   []LittleComment `json:"-" bson:"comments"`
	User       interface{}     `json:"owner" bson:"user"`
}
