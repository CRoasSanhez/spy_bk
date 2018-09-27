package serializers

import "spyc_backend/app/models"

// SessionsSerializer ...
type SessionsSerializer struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

// Cast ...
func (s SessionsSerializer) Cast(data interface{}) Serializer {
	serializer := new(SessionsSerializer)

	if model, ok := data.(models.User); ok {
		serializer.User = UserSerializer{}.Cast(model)
		serializer.Token = s.Token
	}

	return serializer
}
