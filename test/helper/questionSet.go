package helper

import (
	"base/app/model"

	faker "github.com/brianvoe/gofakeit/v5"
)

func seedQuestionSetInput(isNew bool) model.QuestionSet {
	input := model.QuestionSet{
		Name: faker.Sentence(2),
	}

	if !isNew {
		input.Questions = []model.Question{
			model.Question{
				Question: faker.Sentence(2),
				Answers: []string{
					faker.Sentence(3),
					faker.Sentence(3),
				},
			},
		}
	}

	return input
}
