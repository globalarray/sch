package repository_model

import "time"

type QuizResult struct {
	QuizID      int64     `db:"quiz_id"`
	UserID      int64     `db:"user_id"`
	Score       int       `db:"score"`
	CompletedIn time.Time `db:"completed_in"`
}

func NewQuizResult(quizID int64, userID int64, score int) QuizResult {
	return QuizResult{
		QuizID: quizID,
		UserID: userID,
		Score:  score,
	}
}
