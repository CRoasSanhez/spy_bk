package models

import (
	"log"
	"spyc_backend/app"

	"spyc_backend/app/core"

	"github.com/Reti500/mgomap"
)

// Question struct
type Question struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	Title   string   `json:"title" bson:"title" validate:"regexp=^[a-zA-Z0-9.():¿?=!¡ ]*$"`
	Options []Option `json:"options" bson:"options"`
}

// Option struct
type Option struct {
	Value    string `json:"value" bson:"value" validate:"regexp=^[a-zA-Z0-9.():¿?=!¡ ]*$"`
	ID       int    `json:"id" bson:"id"`
	IsAnswer bool   `json:"isAnswer" bson:"isAnswer"`
}

// GetDocumentName ...
func (m *Question) GetDocumentName() string {
	return "questions"
}

// SetStatus sets Question status according to constants
func (m *Question) SetStatus(status string) {
	m.Status.Code = core.GameStatus[status]
	m.Status.Name = status
}

// GetQuestionByID ...
func (m *Question) GetQuestionByID(id string) error {
	Quest, _ := app.Mapper.GetModel(&Question{})

	err := Quest.Find(id).Exec(m)
	if err != nil {
		return err
	}
	return nil
}

// FindQuestionsByTarget returns an array of questions
func (m *Question) FindQuestionsByTarget(targetID string) (questions []Question, err error) {
	target := Target{}
	Target, _ := app.Mapper.GetModel(&target)

	err = Target.Find(targetID).Exec(&target)
	if err != nil {
		return questions, err
	}

	question := Question{}

	for i := 0; i < len(target.Questions); i++ {
		err = question.GetQuestionByID(target.Questions[i])
		if err != nil {
			log.Print(err.Error())
			continue
		}
		questions = append(questions, question)
	}
	return questions, nil
}

// DeleteQuestion is a logicall delete of question
func (m *Question) DeleteQuestion() bool {
	Question, _ := app.Mapper.GetModel(&Question{})

	m.SetStatus(core.StatusInactive)
	m.Deleted = true

	err := Question.Update(m)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}
