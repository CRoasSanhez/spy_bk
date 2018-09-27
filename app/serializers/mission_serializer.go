package serializers

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
)

// MissionSerializer ...
type MissionSerializer struct {
	ID            string          `json:"id"`
	Title         interface{}     `json:"title"`
	Description   interface{}     `json:"description"`
	ImageURL      string          `json:"image_url"`
	Type          string          `json:"type"`
	CoverPicture  interface{}     `json:"cover_picture"`
	StartDate     string          `json:"start_date"`
	EndDate       string          `json:"end_date"`
	Geolocation   *models.Geo     `json:"geolocation"`
	Status        mgomap.Status   `json:"status"`
	Targets       interface{}     `json:"steps"`
	Subscribed    bool            `json:"subscribed"`
	Periods       []models.Period `json:"periods"`
	Advertisement interface{}     `json:"advertisement"`

	// Used for internatiolaization
	Lang string `json:"lang"`
}

// MissionsSerializer ...
type MissionsSerializer struct {
	Missions interface{} `json:"games"`

	// Used for internatiolaization
	Lang string `json:"-"`
}

// Cast ...
func (s MissionSerializer) Cast(data interface{}) Serializer {
	serializer := new(MissionSerializer)

	if model, ok := data.(models.Mission); ok {
		serializer.ID = model.GetID().Hex()
		serializer.Title = InternationalizeSerializer{}.Get(model.Title, s.Lang)
		serializer.Description = InternationalizeSerializer{}.Get(model.Description, s.Lang)
		serializer.Type = model.Type
		serializer.StartDate = model.StartDate.Format(core.MXTimeFormat)
		serializer.EndDate = model.EndDate.Format(core.MXTimeFormat)
		serializer.Geolocation = model.Geolocation
		serializer.Status = model.Status
		serializer.Subscribed = model.Subscribed
		serializer.Targets = Serialize(model.Targets, TargetSerializer{Lang: s.Lang})
		serializer.Periods = model.Periods
		serializer.CoverPicture = Serialize(model.Attachment, AttachmentSerializer{
			Parent: model.GetDocumentName(), ParentID: model.GetID().Hex(), Field: "attachment", VerifyURL: true,
		})

		if model.Advertisement != "" {
			serializer.Advertisement = Serialize(s.FindAd(model.Advertisement), AdvertisementSerializer{})
		} else {
			serializer.Advertisement = nil
		}
	}

	return serializer
}

// Cast ...
func (s MissionsSerializer) Cast(data interface{}) Serializer {
	serializer := new(MissionsSerializer)

	if model, ok := data.(models.Missions); ok {
		serializer.Missions = Serialize(model.Missions, MissionSerializer{})
	}
	return serializer
}

// FindAd ...
func (s MissionSerializer) FindAd(id string) models.Advertisement {
	var advert models.Advertisement
	if Ad, ok := app.Mapper.GetModel(&advert); ok {
		if err := Ad.Find(id).Exec(&advert); err != nil {
			revel.ERROR.Printf("ERROR FIND Advertisement %s --- %s", id, err.Error())
		}
		return advert
	}
	return advert
}
