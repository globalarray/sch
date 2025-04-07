package button

import (
	"benzo/internal/repository"
	"benzo/internal/repository/repository_model"
	"benzo/internal/service"
	"errors"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strconv"
	"strings"
)

type AnswerQuestion struct {
	log *slog.Logger
	StudentButton
}

var (
	ErrInvalidQuestion = errors.New("invalid question")
)

func (b *AnswerQuestion) Run(_ *tele.Bot, ctx tele.Context, args []string) error {
	if len(args) != 2 {
		return ErrInvalidUsage
	}

	id := ctx.Callback().Sender.ID
	languageCode := ctx.Callback().Sender.LanguageCode

	questionID, err := strconv.ParseInt(args[0], 10, 64)

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

	correctAnswer := strings.Split(q.Answers, ";")[0]

	if err := repository.Repo().SaveNewQuizProgress(repository_model.NewQuizProgress(id, q.QuizID, questionID, args[1], args[1] == correctAnswer)); err != nil {
		return err
	}

	return service.Quiz().ProcessQuiz(ctx, q.QuizID, id, languageCode)
}

func (b *AnswerQuestion) Endpoint() string {
	return "question_answer"
}
