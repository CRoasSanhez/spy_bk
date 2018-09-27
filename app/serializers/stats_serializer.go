package serializers

import "spyc_backend/app/models"

// StatsSerializer serializer
type StatsSerializer struct {
	Missions   interface{} `json:"missions" bson:"missions"`
	Challenges interface{} `json:"challenges" bson:"challenges"`
	Account    interface{} `json:"account" json:"account"`
}

// Cast bind instareface to stats model
func (s StatsSerializer) Cast(data interface{}) Serializer {
	var serializer = new(StatsSerializer)

	if model, ok := data.(models.Stats); ok {
		serializer.Missions = model.Missions
		serializer.Challenges = model.Challenges
		serializer.Account = model.Account
	}

	return serializer
}
