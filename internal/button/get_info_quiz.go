package button

import (
	"benzo/internal/lang"
	"benzo/internal/quiz"
	"benzo/internal/repository"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type GetInfoQuiz struct {
	log *slog.Logger
}

func (b *GetInfoQuiz) Run(_ *tele.Bot, ctx tele.Context, args []string) error {
	if len(args) != 1 {
		return ErrInvalidUsage
	}

	quizID, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	q, err := repository.Repo().GetQuizByID(quizID)

	if err != nil {
		return err
	}

	languageCode := ctx.Callback().Sender.LanguageCode

	if q.ID != quizID {
		return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, "quiz does not exist"))
	}

	encoded, err := quiz.Encode(q)

	if err != nil {
		return err
	}

	questions, err := repository.Repo().GetQuestionsByQuizID(quizID)

	if err != nil {
		return err
	}

	results, err := repository.Repo().SelectQuizResultByQuizID(quizID)

	if err != nil {
		return err
	}

	var response []string

	response = append(response, i18n.Translatef(lang.QuizInfoTitle, languageCode, q.Name))
	response = append(response, "")
	response = append(response, i18n.Translatef(lang.QuizInfoInvitationLinkLine, languageCode, encoded))
	response = append(response, i18n.Translatef(lang.QuizInfoQuestionsCountLine, languageCode, len(questions)))

	if len(results) > 0 {
		response = append(response, "\n")
		response = append(response, i18n.Translatef(lang.QuizInfoCompletedUsersLine, languageCode, len(results)))

		for _, result := range results {
			u, err := repository.Repo().GetUserByTelegramID(result.UserID)

			if err != nil {
				continue
			}

			response = append(response, "\n"+i18n.Translatef(lang.QuizInfoUserResultLine, languageCode, result.CompletedIn.Format(time.DateTime), u.Name, u.Surname, result.Score, len(questions)))
		}
	}

	selector := &tele.ReplyMarkup{}

	selector.Inline(selector.Row(selector.Data(i18n.Translatef(lang.QuizRemoveBtn, languageCode), fmt.Sprintf("quiz_remove-%d", quizID))))

	return ctx.Send(strings.Join(response, "\n"), selector)
}

func (b *GetInfoQuiz) NeedRightsLevel() int {
	return role.RightsLevel(role.Teacher{})
}

func (b *GetInfoQuiz) Endpoint() string {
	return "get_info_quiz"
}
