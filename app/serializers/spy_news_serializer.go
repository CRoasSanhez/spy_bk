package serializers

import (
	"spyc_backend/app"
	"spyc_backend/app/models"

	"github.com/Reti500/mgomap"
	"gopkg.in/mgo.v2/bson"
)

// SpyNewsSerializer serializer
type SpyNewsSerializer struct {
	ID           string        `json:"id"`
	Title        interface{}   `json:"title"`
	Description  interface{}   `json:"description"`
	Attachment   interface{}   `json:"attachment"`
	Reactions    interface{}   `json:"reactions"`
	Comments     []interface{} `json:"comments"`
	Lang         string        `json:"lang"`
	MaxComments  int           `json:"-"`
	CurrentReact interface{}   `json:"current_react"`
}

// Cast bind serializer with spynews data model
func (s SpyNewsSerializer) Cast(data interface{}) Serializer {
	var serializer = new(SpyNewsSerializer)

	if model, ok := data.(models.SpyNews); ok {
		if s.MaxComments <= 0 {
			s.MaxComments = 5
		}

		serializer.ID = model.GetID().Hex()
		serializer.Title = InternationalizeSerializer{}.Get(model.Title, s.Lang)
		serializer.Description = InternationalizeSerializer{}.Get(model.Description, s.Lang)
		serializer.Attachment = InternationalizeSerializer{}.Get(model.Attachment, s.Lang)
		serializer.Lang = s.Lang
		serializer.CurrentReact = model.Extra["current_interaction"]

		if News, ok := app.Mapper.GetModel(&model); ok {
			var interactions []bson.M
			var pipe = mgomap.Aggregate{}.Match(bson.M{"_id": model.GetID()}).Add(
				bson.M{"$unwind": "$interactions"},
			).Add(bson.M{"$group": bson.M{"_id": "$interactions.name", "total": bson.M{"$sum": 1}}})
			if err := News.Pipe(pipe, &interactions); err == nil {
				serializer.Reactions = interactions
			}
		}
	}

	return serializer
}
