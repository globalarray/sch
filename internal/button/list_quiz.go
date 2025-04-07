package button

import (
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log/slog"
)

type ListQuiz struct {
	log *slog.Logger
}

func (b *ListQuiz) Run(_ *tele.Bot, ctx tele.Context, _ []string) error {
	id := ctx.Callback().Sender.ID
	languageCode := ctx.Callback().Sender.LanguageCode

	quizzes, err := repository.Repo().GetQuizzesCreatedByUserID(id)

	if err != nil {
		return err
	}

	if len(quizzes) == 0 {
		return ctx.Send(i18n.Translatef(lang.QuizEmptyListTitle, languageCode))
	}

	selector := &tele.ReplyMarkup{}

	var rows []tele.Row

	for _, quiz := range quizzes {
		rows = append(rows, selector.Row(selector.Data(quiz.Name, fmt.Sprintf("get_info_quiz-%d", quiz.ID))))
	}

	selector.Inline(rows...)

	return ctx.Send(i18n.Translatef(lang.QuizListTitle, languageCode), selector)
}

func (b *ListQuiz) NeedRightsLevel() int {
	return role.RightsLevel(role.Teacher{})
}

func (b *ListQuiz) Endpoint() string {
	return "quiz_list"
}
