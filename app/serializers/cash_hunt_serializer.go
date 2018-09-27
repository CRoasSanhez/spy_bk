package serializers

import (
	"spyc_backend/app/models"
)

// CashHuntSerializer serializer to games/geo
type CashHuntSerializer struct {
	Missions interface{} `json:"games"`
	Rewards  interface{} `json:"rewards"`
	Lang     string      `json:"lang"`
}

// Cast ...
func (s CashHuntSerializer) Cast(data interface{}) Serializer {
	serializer := new(CashHuntSerializer)

	if model, ok := data.(models.CashHuntResponse); ok {
		serializer.Missions = Serialize(model.Missions, MissionSerializer{Lang: s.Lang})
		serializer.Rewards = Serialize(model.Rewards, RewardSerializer{})
	}
	return serializer
}
