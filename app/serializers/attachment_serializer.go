package serializers

import "spyc_backend/app/models"

// AttachmentSerializer serializer
type AttachmentSerializer struct {
	CurrentName string `json:"current_name"`
	Format      string `json:"format"`
	Size        int64  `json:"size"`
	URL         string `json:"file_url"`

	// Update and signing URL params
	Parent    string `json:"-"`
	ParentID  string `json:"-"`
	Field     string `json:"-"`
	VerifyURL bool   `json:"-"`
}

// Cast ...
func (s AttachmentSerializer) Cast(data interface{}) Serializer {
	serializer := new(AttachmentSerializer)

	if model, ok := data.(models.Attachment); ok {
		if serializer.VerifyURL {
			if model.HasExpired() {
				model.UpdateURLParentField(serializer.Parent, serializer.ParentID, serializer.Field)
			}
		} else {
			model.UpdateURL()
		}

		serializer.CurrentName = model.CurrentName
		serializer.Format = model.Format
		serializer.Size = model.Size
		serializer.URL = model.URL
	}

	return serializer
}
