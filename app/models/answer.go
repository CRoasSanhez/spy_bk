package models

// AnswerRequest is used for parse request to
// PIN, TEXT, QRdata, CheckIn
type AnswerRequest struct {
	Data        string `json:"data" bson:"data" validate:"regexp=^[a-zA-Z0-9 ]*$"`
	Type        string `json:"-" bson:"type"`
	Geolocation Geo    `json:"geolocation" bson:"geolocation"`
}
