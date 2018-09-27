package serializers

import "spyc_backend/app/models"

// InternationalizeSerializer return data on selected lang
type InternationalizeSerializer struct {
	Data interface{}
	Lang string
}

// GetData returns data serialized
func (s InternationalizeSerializer) GetData() interface{} {
	return s.Data
}

// Get returns lang selected
func (s InternationalizeSerializer) Get(data interface{}, lang string) interface{} {
	if model, ok := data.(models.Internationalized); ok {
		if model[lang] != nil {
			return model[lang]
		}

		for _, v := range model {
			return v
		}
	}

	return nil
}

// SetLang set a current lang
func (s InternationalizeSerializer) SetLang(lang string) InternationalizeSerializer {
	s.Lang = lang
	return s
}

// Cast ....
func (s InternationalizeSerializer) Cast(data interface{}) Serializer {
	var serializer = new(InternationalizeSerializer)

	if model, ok := data.(models.Internationalized); ok {
		serializer.Data = model[s.Lang]
	}

	return serializer
}
