package models

import (
	"math/rand"
	"spyc_backend/app/core"
	"strings"
	"time"
)

// Saas ...
type Saas struct {
	Name  string `json:"name,omitempty" bson:"name" validate:"regexp=^[a-zA-Z0-9 ]*$"`
	Token string `json:"token,omitempty" bson:"token" validate:"regexp=^[a-zA-Z0-9.-_:¿?=!¡]*$"`
}

// Init ...
func (s *Saas) Init(email, number string) {
	rand.Seed(time.Now().Unix())
	r1 := rand.Intn(5) + 1
	pref := strings.Split(email, "@")[0] + "_" + strings.Split(email, "@")[1][:1]
	suf := number[r1:10]
	s.Name = pref + "_" + suf[:5]
	s.Token = core.GenerateToken(core.LetterRunes, 32)
}
