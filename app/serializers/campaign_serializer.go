package serializers

import (
	"spyc_backend/app/core"
	"spyc_backend/app/models"
)

// CampaignSerializer ...
type CampaignSerializer struct {
	ID        string      `json:"id"`
	Name      interface{} `json:"name"`
	StartDate string      `json:"start_date"`
	EndDate   string      `json:"end_date"`
	Budget    interface{} `json:"budget"`
}

// Cast ...
func (s CampaignSerializer) Cast(data interface{}) Serializer {
	serializer := new(CampaignSerializer)

	if model, ok := data.(models.Campaign); ok {
		serializer.ID = model.GetID().Hex()
		serializer.Name = model.Name
		serializer.StartDate = model.StartDate.Format(core.MXTimeFormat)
		serializer.EndDate = model.EndDate.Format(core.MXTimeFormat)
		serializer.Budget = model.Budget
	}

	return serializer
}
