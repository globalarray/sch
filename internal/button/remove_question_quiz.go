package button

import (
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	"errors"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strconv"
)

type RemoveQuestionQuiz struct {
	log *slog.Logger
}

var errInvalidQuestionID = errors.New("invalid question id")

func (b *RemoveQuestionQuiz) Run(ctx tele.Context, args []string) error {
	if len(args) != 1 {
		return ErrInvalidUsage
	}

	questionID, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return ErrInvalidUsage
	}

	q, err := repository.Repo().GetQuestionByID(questionID)

	if err != nil {
		return err
	}

	languageCode := ctx.Callback().Sender.LanguageCode

	if q.ID != questionID {
		return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, errInvalidQuestionID.Error()))
	}

	questions, err := repository.Repo().GetQuestionsByQuizID(q.QuizID)

	if err != nil {
		return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))
	}

	selector := &tele.ReplyMarkup{}

	buttons := []tele.Btn{selector.Data(i18n.Translatef(lang.QuizAddNewQuestionBtn, languageCode), fmt.Sprintf("add_new_question_quiz-%d", q.QuizID))}

	if len(questions) > 1 {
		buttons = append(buttons, selector.Data(i18n.Translatef(lang.QuizAddingQuestionsStopBtn, languageCode), fmt.Sprintf("get_info_quiz-%d", q.QuizID)))
	}

	selector.Inline(selector.Row(buttons...))

	if err := repository.Repo().RemoveQuestionByID(questionID); err != nil {
		return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))
	}

	return ctx.Send(i18n.Translatef(lang.QuizQuestionDeletedMessage, languageCode), selector)
}

func (b *RemoveQuestionQuiz) NeedRightsLevel() int {
	return role.RightsLevel(role.Teacher{})
}

func (b *RemoveQuestionQuiz) Endpoint() string {
	return "remove_question_quiz"
}
