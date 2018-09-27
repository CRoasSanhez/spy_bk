package serializers

import "spyc_backend/app/models"

// CommentSerializer serializer
type CommentSerializer struct {
	Message string      `json:"message"`
	Owner   interface{} `json:"owner"`
}

// Cast bind comment model to serializer
func (s CommentSerializer) Cast(data interface{}) Serializer {
	var serializer = new(CommentSerializer)

	if model, ok := data.(models.Comment); ok {
		serializer.Message = model.Message
		serializer.Owner = Serialize(model.User, CommentOwnerSerializer{})
	}

	return serializer
}
