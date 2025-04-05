package button

import (
	"benzo/internal/callback"
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strconv"
)

type AddQuestionQuiz struct {
	log *slog.Logger
}

func (b *AddQuestionQuiz) Run(_ *tele.Bot, ctx tele.Context, args []string) error {
	if len(args) != 1 {
		return ErrInvalidUsage
	}

	id := ctx.Callback().Sender.ID
	languageCode := ctx.Callback().Sender.LanguageCode

	quizID, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	quiz, err := repository.Repo().GetQuizByID(quizID)

	if err != nil {
		return err
	}

	if quiz.ID != quizID {
		return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, "quiz does not exist"))
	}

	callback.Subscribe(id, createQuestionCallback(quizID))

	return ctx.Send(i18n.Translatef(lang.QuizQuestionSetNameMessage, languageCode))
}

func (b *AddQuestionQuiz) NeedRightsLevel() int {
	return role.RightsLevel(role.Teacher{})
}

func (b *AddQuestionQuiz) Endpoint() string {
	return "add_new_question_quiz"
}
