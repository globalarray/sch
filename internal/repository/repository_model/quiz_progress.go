package repository_model

type QuizProgress struct {
	UserID     int64 `db:"user_id"`
	QuizID     int64 `db:"quiz_id"`
	QuestionID int64 `db:"question_id"`
	Answer     int   `db:"answer"`
	Correct    bool  `db:"correct"`
}

func NewQuizProgress(userID, quizID, questionID int64, answer int, correct bool) QuizProgress {
	return QuizProgress{
		UserID:     userID,
		QuizID:     quizID,
		QuestionID: questionID,
		Answer:     answer,
		Correct:    correct,
	}
}
