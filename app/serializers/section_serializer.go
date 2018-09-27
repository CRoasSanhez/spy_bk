package serializers

import (
	"spyc_backend/app"
	"spyc_backend/app/models"

	"gopkg.in/mgo.v2/bson"
)

// SectionSerializer serializer
type SectionSerializer struct {
	Name          string      `json:"name"`
	CanAccess     bool        `json:"can_access"`
	TimeLess      float64     `json:"time_less"`
	LastAccess    int64       `json:"last_access"`
	Advertisement interface{} `json:"advertisement"`
	AccessTime    float64     `json:"access_time"`
}

// Cast ...
func (s SectionSerializer) Cast(data interface{}) Serializer {
	var serializer = new(SectionSerializer)

	if model, ok := data.(models.Section); ok {
		serializer.Name = model.Name
		serializer.CanAccess = model.CanAccess
		serializer.LastAccess = model.LastAccess.Unix()
		serializer.TimeLess = model.TimeLess.Seconds()
		serializer.AccessTime = model.AccessTime.Seconds()

		var advertisement models.Advertisement
		if Advertisement, ok := app.Mapper.GetModel(&advertisement); ok {
			if err := Advertisement.Query(
				bson.M{"tags": bson.M{"$in": []string{"section_" + model.Name}}}).Exec(&advertisement); err == nil {
				serializer.Advertisement = Serialize(advertisement, AdvertisementSerializer{})
			}
		}
	}

	return serializer
}
