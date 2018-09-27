package models

import "github.com/Reti500/mgomap"

// SpyNews model
type SpyNews struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Title        Internationalized      `json:"titles" bson:"titles"`
	Description  Internationalized      `json:"descriptions" bson:"descriptions"`
	Attachment   Internationalized      `json:"-" bson:"attachments"`
	Langs        []string               `json:"langs" bson:"languages"`
	Interactions []React                `json:"-" bson:"interactions"`
	Extra        map[string]interface{} `json:"-" bson:"extra"`
}

// GetDocumentName set a name for document on MongoDB
func (m SpyNews) GetDocumentName() string {
	return "spy_news"
}

// SetTitle set a title field by lang
func (m *SpyNews) SetTitle(lang, title string) {
	m.Title[lang] = title
}

// SetDescription set a description field by lang
func (m *SpyNews) SetDescription(lang, desc string) {
	m.Description[lang] = desc
}

// SetAttachment set a attachment field by lang
func (m *SpyNews) SetAttachment(lang string, attachment Attachment) {
	m.Attachment[lang] = attachment
}
