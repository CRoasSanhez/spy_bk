package v2

import (
	"spyc_backend/app"
	"spyc_backend/app/models"
	"spyc_backend/app/serializers"

	"github.com/Reti500/mgomap"
	"github.com/revel/revel"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MessagesController ...
type MessagesController struct {
	BaseController
}

// Index return last 20 messages from chat
func (c MessagesController) Index(user string) revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	var messages []models.Message
	var pipe = mgomap.Aggregate{}.Match(bson.M{
		"sender": c.CurrentUser.Saas.Name, "receiver": user}).Sort(bson.M{"created_at": -1})

	if err := models.AggregateQuery(&models.Message{}, pipe, &messages); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(messages, "success", 200, nil)
}

// UploadFile ...
func (c MessagesController) UploadFile() revel.Result {
	if !c.GetCurrentUser() {
		return c.ForbiddenResponse()
	}

	if len(c.Params.Files["file"]) != 1 {
		return c.ErrorResponse(nil, "File not found", 400)
	}

	file := models.File{
		Action:  "messages",
		Current: false,
		Ref: &mgo.DBRef{
			Collection: c.CurrentUser.GetDocumentName(),
			Database:   app.Mapper.DatabaseName,
			Id:         c.CurrentUser.GetID().Hex(),
		},
	}

	_, err := file.Upload(c.Params.Files["file"][0])
	if err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	if err := models.CreateDocument(&file); err != nil {
		return c.ErrorResponse(err, err.Error(), 400)
	}

	return c.SuccessResponse(file, "Created!!", 200, serializers.FileSerializer{})
}

// DeleteFile remove file from messages
func (c MessagesController) DeleteFile() revel.Result {
	return c.SuccessResponse("ok", "ok", 200, nil)
}
