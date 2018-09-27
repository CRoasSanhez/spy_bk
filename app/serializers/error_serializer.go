package serializers

// ErrorSerializer ...
type ErrorSerializer struct {
	SerializerModel
}

// Cast ...
func (s ErrorSerializer) Cast(data interface{}) Serializer {
	serializer := new(ErrorSerializer)

	return serializer
}
