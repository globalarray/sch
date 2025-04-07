package button

import (
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strconv"
)

type RemoveQuiz struct {
	log *slog.Logger
}

func (b *RemoveQuiz) Run(_ *tele.Bot, ctx tele.Context, args []string) error {
	if len(args) != 1 {
		return ErrInvalidUsage
	}

	id := ctx.Callback().Sender.ID
	languageCode := ctx.Callback().Sender.LanguageCode

	quizID, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	q, err := repository.Repo().GetQuizByID(quizID)

	if err != nil {
		return err
	}

	if q.CreatedBy != id {
		u, err := repository.Repo().GetUserByTelegramID(id)

		if err != nil {
			return err
		}

		r, err := role.FromName(u.Role)

		if err != nil {
			return err
		}

		if role.RightsLevel(r) != role.RightsLevel(role.Admin{}) {
			return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, "you have not rights for this quiz"))
		}
	}

	if err := repository.Repo().RemoveQuizByID(quizID); err != nil {
		return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, "failed to remove quiz"))
	}

	if err := repository.Repo().RemoveQuestionsByQuizID(quizID); err != nil {
		return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, "failed to remove questions"))
	}

	if err := repository.Repo().RemoveResultsByQuizID(quizID); err != nil {
		return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, "failed to remove results"))
	}

	if err := repository.Repo().RemoveProgressesByQuizID(quizID); err != nil {
		return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, "failed to remove progresses"))
	}

	return ctx.Send(i18n.Translatef(lang.QuizRemovedMessage, languageCode, q.Name))
}

func (b *RemoveQuiz) NeedRightsLevel() int {
	return role.RightsLevel(role.Teacher{})
}

func (b *RemoveQuiz) Endpoint() string {
	return "quiz_remove"
}
