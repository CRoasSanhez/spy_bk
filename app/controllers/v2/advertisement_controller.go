package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"gopkg.in/mgo.v2/bson"

	"github.com/revel/revel"
)

// AdvertisementController controller
type AdvertisementController struct {
	BaseController
}

// Create add new Advertisement
func (c AdvertisementController) Create() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var advertisement models.Advertisement

	if err := c.PermitParams(&advertisement, false, "sponsor", "tags", "name", "type", "delay"); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if err := models.CreateDocument(&advertisement); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(advertisement, core.StatusSuccess, 200, serializers.AdvertisementSerializer{})
}

// UploadAttachment upload an attachment for advertisement
func (c AdvertisementController) UploadAttachment(id string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if len(c.Params.Files["attachment"]) <= 0 {
		return c.ErrorResponse(nil, "File not found", 400)
	}

	var advertisement models.Advertisement

	if err := models.GetDocument(id, &advertisement); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if advertisement.Attachment.PATH != "" {
		advertisement.Attachment.Remove()
	}

	if err := advertisement.Attachment.Init(models.AsDocumentBase(&c.CurrentUser), c.Params.Files["attachment"][0]); err != nil {
		return c.ErrorResponse(err, err.Error(), 402)
	}

	if err := advertisement.Attachment.Upload(); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	var fields = bson.M{"attachment": advertisement.Attachment}
	if err := models.SetDocument(advertisement.GetDocumentName(), advertisement.GetID().Hex(), fields); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(advertisement, core.StatusSuccess, 200, serializers.AdvertisementSerializer{})
}

// GetAdvertisementsBySponsor ...
func (c AdvertisementController) GetAdvertisementsBySponsor(idSponsor string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if !bson.IsObjectIdHex(idSponsor) {
		return c.ErrorResponse(nil, c.Message("error.invalid", ""), core.ModelStatus[core.StatusInvalidID])
	}

	var advertisements []models.Advertisement
	var query = []bson.M{bson.M{"sponsor_id": idSponsor}}

	if Ad, ok := app.Mapper.GetModel(&models.Advertisement{}); ok {
		if err := Ad.FindWithOperator("$and", query).Exec(&advertisements); err != nil {
			return c.ErrorResponse(err, err.Error(), 400)
		}
		return c.SuccessResponse(advertisements, "success", 200, serializers.AdvertisementSerializer{})
	}
	return c.ServerErrorResponse()
}
