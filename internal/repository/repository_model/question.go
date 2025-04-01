package repository_model

type Question struct {
	QuizID   int64  `db:"quiz_id"`
	Question string `db:"question"`
	Answers  []string
}

func NewQuestion(quizID int64, question string, correctAnswer string, otherAnswers ...string) Question {
	return Question{
		QuizID:   quizID,
		Question: question,
		Answers:  append([]string{correctAnswer}, otherAnswers...),
	}
}
