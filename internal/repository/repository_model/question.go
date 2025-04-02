package repository_model

import "strings"

type Question struct {
	ID       int64  `db:"id"`
	QuizID   int64  `db:"quiz_id"`
	Question string `db:"question"`
	Answers  string `db:"answers"`
}

func NewQuestion(quizID int64, question string, correctAnswer string, otherAnswers ...string) Question {
	return Question{
		QuizID:   quizID,
		Question: question,
		Answers:  strings.Join(append([]string{correctAnswer}, otherAnswers...), ";"),
	}
}
