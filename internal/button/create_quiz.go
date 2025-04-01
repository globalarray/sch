package button

import (
	"benzo/internal/callback"
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/repository/repository_model"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	tele "gopkg.in/telebot.v4"
	"log/slog"
)

type CreateQuiz struct {
	log *slog.Logger
}

const (
	minQuizNameLength int = 5
	maxQuizNameLength int = 20

	minQuestionLength int = 5
	maxQuestionLength int = 50
)

func (b *CreateQuiz) Run(bot *tele.Bot, ctx tele.Context, args []string) error {
	id := ctx.Callback().Sender.ID
	languageCode := ctx.Callback().Sender.LanguageCode

	callback.Subscribe(id, b.createQuizCallback)
}

func (b *CreateQuiz) createQuizCallback(_ *tele.Bot, ctx tele.Context) bool {
	languageCode := ctx.Message().Sender.LanguageCode
	id := ctx.Callback().Sender.ID
	quizName := ctx.Message().Text

	if len(quizName) < minQuizNameLength || len(quizName) > maxQuizNameLength {
		_ = ctx.Send(i18n.Translatef(lang.QuizCreateNameLengthInvalid, languageCode))
		return false
	}

	quizId, err := repository.Repo().SaveNewQuiz(repository_model.NewQuiz(quizName, id))

	if err != nil {
		_ = ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))

		return false
	}

	_ = ctx.Send(i18n.Translatef(lang.QuizCreatedMessage, languageCode))

	return true
}

func (b *CreateQuiz) createQuestionCallback(quizID int64) callback.CallbackFunc {
	return func(_ *tele.Bot, ctx tele.Context) bool {
		id := ctx.Message().Sender.ID
		languageCode := ctx.Message().Sender.LanguageCode
		question := ctx.Message().Text

		if len(question) < minQuestionLength || len(question) > maxQuestionLength {
			_ = ctx.Send(i18n.Translatef(lang.QuizQuestionNameLengthInvalid, languageCode))
			return false
		}

		if err := repository.Repo().SaveNewQuestion(repository_model.NewQuestion(quizID, question, "", "")); err != nil {
			//koroche pust teper answeri kidaet potom etu xuetu v cikl

			return false
		}
	}
}

func (b *CreateQuiz) NeedRightsLevel() int {
	return role.RightsLevel(role.Teacher{})
}

func (b *CreateQuiz) Endpoint() string {
	return "quiz_create"
}
