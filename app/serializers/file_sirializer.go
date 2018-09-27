package serializers

import (
	"spyc_backend/app"
	"spyc_backend/app/models"

	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

// FileSerializer ...
type FileSerializer struct {
	Action      string      `json:"action"`
	CurrentName string      `json:"current_name"`
	Format      string      `json:"format"`
	Size        int64       `json:"size"`
	URL         string      `json:"url"`
	Resource    interface{} `json:"-"`
}

// Cast ...
func (s FileSerializer) Cast(data interface{}) Serializer {
	serializer := new(FileSerializer)

	if model, ok := data.(models.File); ok {
		serializer.Action = model.Action
		serializer.CurrentName = model.CurrentName
		serializer.Format = model.Format
		serializer.Size = model.Size
		serializer.URL = model.URL
	} else {
		revel.INFO.Print("Not serialize file :(")
	}

	return serializer
}

// DoubleFileSerializer ...
type DoubleFileSerializer struct {
	File      interface{} `json:"file"`
	Thumbnail interface{} `json:"thumbnail"`
	Resource  interface{} `json:"-"`
}

// Cast ...
func (s DoubleFileSerializer) Cast(data interface{}) Serializer {
	serializer := new(DoubleFileSerializer)

	if model, ok := data.(models.ChallengeFiles); ok {
		if model.Attachment.HasExpired() {
			//model.Attachment.UpdateURL()
			model.Attachment.UpdateURLParentQuery("rewards", "files", "attachement", s.Resource)
		}
		if model.Thumbnail.HasExpired() {
			//model.Thumbnail.UpdateURL()
			model.Attachment.UpdateURLParentQuery("rewards", "files", "thumbnail", s.Resource)
		}
		serializer.File = Serialize(model.Attachment, AttachmentSerializer{})
		serializer.Thumbnail = Serialize(model.Thumbnail, AttachmentSerializer{})
	} else {
		revel.INFO.Print("Not serialize file :(")
	}

	return serializer
}

// FindFile ...
func (s DoubleFileSerializer) FindFile(id string) interface{} {
	file := models.File{}
	File, _ := app.Mapper.GetModel(&file)

	if !bson.IsObjectIdHex(id) {
		return file
	}

	File.Find(id).Exec(&file)

	return file
}
