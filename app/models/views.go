package models

// Views model
type Views struct {
	UserID    string `json:"user_id" bson:"user_id"`
	CloseTime int    `json:"close_time" bson:"close_time"`
	Completed bool   `json:"completed" bson:"completed"`
	Through   string `json:"through" bson:"through"`
	Action    string `json:"action" bson:"action"`
	Total     int    `json:"total" bson:"total"`
}
