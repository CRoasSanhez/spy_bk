package models

import (
	"strings"
	"time"
)

// TimeFormats
const (
	timeFormat string = "02/01/2006 15:04"
)

// CustomTime ...
type CustomTime struct {
	time.Time
}

// Period ...
type Period struct {
	StartDate time.Time `json:"start_date" bson:"start_date"`
	EndDate   time.Time `json:"end_date" bson:"end_date"`
}

// UnmarshalJSON ...
func (m *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	if s == "null" {
		m.Time = time.Time{}
		return
	}

	m.Time, err = time.Parse(timeFormat, s)
	return
}
