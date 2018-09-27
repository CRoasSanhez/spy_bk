package serializers

import "spyc_backend/app/models"

// GroupHistoriesSerializer ...
type GroupHistoriesSerializer struct {
	User      interface{}   `json:"user"`
	Histories []interface{} `json:"histories"`
	// TotalSeen int           `json:"total_seen"`
	// TotalDocs int           `json:"total_histories"`
}

// HistoriesSerializer ...
type HistoriesSerializer struct {
	ID         string      `json:"id"`
	Text       string      `json:"text"`
	Type       string      `json:"type"`
	Anim       string      `json:"anim"`
	Color      string      `json:"color"`
	Seen       bool        `json:"seen"`
	Attachment interface{} `json:"attachment"`
}

// Cast ...
func (s HistoriesSerializer) Cast(data interface{}) Serializer {
	var serializer = new(HistoriesSerializer)

	if model, ok := data.(models.History); ok {
		serializer.ID = model.GetID().Hex()
		serializer.Text = model.Text
		serializer.Type = model.Type
		serializer.Anim = model.Anim
		serializer.Color = model.Color
		serializer.Seen = model.Seen
		serializer.Attachment = Serialize(model.Attachment, AttachmentSerializer{
			Parent: model.GetDocumentName(), ParentID: model.GetID().Hex(), Field: "attachment", VerifyURL: true,
		})
	}

	return serializer
}

// Cast ...
func (s GroupHistoriesSerializer) Cast(data interface{}) Serializer {
	var serializer = new(GroupHistoriesSerializer)

	if model, ok := data.(models.GHistories); ok {
		var user models.User
		user.PublicProfile(model.UserID, "user_name", "personal_data")

		serializer.User = Serialize(user, UserSerializer{})

		for _, v := range model.Histories {
			serializer.Histories = append(serializer.Histories, Serialize(v, HistoriesSerializer{}))
		}
	}

	return serializer
}
