package button

import (
	"benzo/internal/repository"
	"benzo/internal/repository/repository_model"
	"benzo/internal/service"
	"errors"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strconv"
)

type AnswerQuestion struct {
	log *slog.Logger
	StudentButton
}

var (
	ErrInvalidQuestion = errors.New("invalid question")
)

func (b *AnswerQuestion) Run(ctx tele.Context, args []string) error {
	if len(args) != 2 {
		return ErrInvalidUsage
	}

	id := ctx.Callback().Sender.ID
	languageCode := ctx.Callback().Sender.LanguageCode

	questionID, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	answerIdx, err := strconv.Atoi(args[1])

	if err != nil {
		return err
	}

	q, err := repository.Repo().GetQuestionByID(questionID)

	if err != nil {
		return err
	}

	if q.ID != questionID {
		return ErrInvalidQuestion
	}

	if err := repository.Repo().SaveNewQuizProgress(repository_model.NewQuizProgress(id, q.QuizID, questionID, answerIdx, answerIdx == 0)); err != nil {
		return err
	}

	return service.Quiz().ProcessQuiz(ctx, q.QuizID, id, languageCode)
}

func (b *AnswerQuestion) Endpoint() string {
	return "question_answer"
}
