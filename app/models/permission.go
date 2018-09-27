package models

// Permission embedded model
type Permission struct {
	Action     string `bson:"action"`
	Name       string `bson:"name"`
	Section    string `bson:"section"`
	TargetType string `bson:"target_type"`
	TargetID   string `bson:"target_id"`
}
