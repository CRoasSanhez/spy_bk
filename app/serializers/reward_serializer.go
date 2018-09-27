package serializers

import (
	"spyc_backend/app/core"
	"spyc_backend/app/models"

	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

// RewardSerializer ...
type RewardSerializer struct {
	ID           string      `json:"id"`
	Name         interface{} `json:"name"`
	Description  interface{} `json:"description"`
	Campaign     interface{} `json:"campaign"`
	CoverPicture interface{} `json:"cover_picture"`
	Coins        int         `json:"coins"`
	Subscriptors int         `json:"subscriptors"`
	Files        interface{} `json:"files"`
	Lang         string      `json:"lang"`
}

// Cast ...
func (s RewardSerializer) Cast(data interface{}) Serializer {
	serializer := new(RewardSerializer)

	if reward, ok := data.(models.Reward); ok {
		serializer.ID = reward.GetID().Hex()
		serializer.Name = InternationalizeSerializer{}.Get(reward.Name, s.Lang)
		serializer.Description = InternationalizeSerializer{}.Get(reward.Description, s.Lang)
		serializer.Coins = reward.Coins
		serializer.Subscriptors = reward.Subscriptors

		var idCampaign, ok = reward.CampaignID.(bson.ObjectId)
		if ok {
			campaign := serializer.FindCampaign(idCampaign.Hex())
			serializer.Campaign = campaign.Name
		} else {
			serializer.Campaign = ""
		}

		if reward.Type == core.ModelTypeChallenge {
			serializer.CoverPicture = nil
			serializer.Files = Serialize(reward.Files, DoubleFileSerializer{Resource: reward.GetID()})
		} else {
			// Validate if reward attachment has expired
			if reward.Attachment.HasExpired() {
				reward.Attachment.UpdateURL()
			}
			serializer.CoverPicture = Serialize(reward.Attachment, AttachmentSerializer{
				Parent: reward.GetDocumentName(), ParentID: reward.GetID().Hex(), Field: "attachment", VerifyURL: true,
			})
		}
	}
	return serializer
}

// RewardDSerializer ...
type RewardDSerializer struct {
	ID          string      `json:"id"`
	Name        interface{} `json:"name"`
	Description interface{} `json:"description"`
	Status      interface{} `json:"status"`
	Lang        string      `json:"lang"`
}

// Cast ...
func (s RewardDSerializer) Cast(data interface{}) Serializer {
	serializer := new(RewardDSerializer)

	if reward, ok := data.(models.Reward); ok {
		serializer.ID = reward.GetID().Hex()
		serializer.Name = InternationalizeSerializer{}.Get(reward.Name, s.Lang)
		serializer.Description = InternationalizeSerializer{}.Get(reward.Description, s.Lang)
		serializer.Status = reward.Status.Name
	}
	return serializer
}

// FindCampaign returns the campaign based on the given id
func (s RewardSerializer) FindCampaign(id string) models.Campaign {
	var campaign models.Campaign
	if err := models.GetDocument(id, &campaign); err != nil {
		revel.ERROR.Printf("ERROR FIND Campaign: %s --- %s", id, err.Error())
	}
	return campaign
}
