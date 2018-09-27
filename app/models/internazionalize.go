package models

// Internationalized field for models
type Internationalized map[string]interface{}

// Get returns key value or first element
func (m Internationalized) Get(key string) interface{} {
	if m[key] != nil {
		return m[key]
	}

	for _, v := range m {
		return v
	}

	return nil
}

// GetString returns the value as string
func (m Internationalized) GetString(key string) string {
	if m[key] != nil {
		if str, ok := m[key].(string); ok {
			return str
		}
	}
	return ""
}
