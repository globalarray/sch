package button

import (
	"benzo/internal/callback"
	"benzo/internal/lang"
	"benzo/internal/quiz"
	"benzo/internal/repository"
	"benzo/internal/ui"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strconv"
	"strings"
)

type QuizResultDetail struct {
	log *slog.Logger
}

func (b *QuizResultDetail) Run(ctx tele.Context, args []string) error {
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

	if q.ID != quizID {
		return quiz.ErrQuizNotFound
	}

	id := ctx.Callback().Sender.ID
	languageCode := ctx.Callback().Sender.LanguageCode

	callback.Subscribe(id, b.detailResultCallback(quizID))

	return ctx.Send(i18n.Translatef(lang.QuizDetailResultDataNeededMessage, languageCode))
}

func (*QuizResultDetail) detailResultCallback(quizID int64) callback.CallbackFunc {
	return func(ctx tele.Context) bool {
		id := ctx.Message().Sender.ID
		languageCode := ctx.Message().Sender.LanguageCode
		targetData := strings.Split(ctx.Message().Text, " ")

		if len(targetData) < 1 {
			_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode, "empty request"))
			return true
		}

		results, err := repository.Repo().GetQuizResultsByQuizID(quizID)

		if err != nil {
			_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode))
			return true
		}

		if len(results) < 1 {
			_ = ctx.Reply(i18n.Translatef(lang.QuizWithoutResultsError, languageCode))
			return true
		}

		if len(targetData) < 3 {
			for i := 0; i <= (3 - len(targetData)); i++ {
				targetData = append(targetData, "")
			}
		}

		questionMapIdToIdx := map[int64]int{}

		questions, err := repository.Repo().GetQuestionsByQuizID(quizID)

		if err != nil {
			_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode))
			return true
		}

		for idx, q := range questions {
			questionMapIdToIdx[q.ID] = idx
		}

		users, err := repository.Repo().GetUsersByPersonalData(targetData[0], targetData[1], targetData[2])

		if err != nil {
			_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))
			return true
		}

		if len(users) == 0 {
			_ = ctx.Reply(i18n.Translatef(lang.QuizDetailResultNotFound, languageCode))

			return true
		}

		for _, u := range users {
			response := []string{i18n.Translatef(lang.QuizDetailResultTitle, languageCode, u.FullName())}
			response = append(response, "")

			progresses, err := repository.Repo().GetQuizProgressByUserID(quizID, id)

			if err != nil {
				_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))
				return true
			}

			if len(progresses) < len(questions) {
				_ = ctx.Reply(i18n.Translatef(lang.QuizDetailResultNotFound, languageCode))
				continue
			}

			for pIdx, p := range progresses {
				idx, ok := questionMapIdToIdx[p.QuestionID]

				if !ok {
					continue
				}

				q := questions[idx]

				answers := strings.Split(q.Answers, ";")

				symbol := ui.CorrectAnswerSymbol

				if !p.Correct {
					symbol = ui.IncorrectAnswerSymbol
				}

				response = append(response, i18n.Translatef(lang.QuizDetailResultQuestionLine, languageCode, pIdx+1, q.Question))

				if p.Correct {
					response = append(response, i18n.Translatef(lang.QuizDetailResultCorrectAnswerLine, languageCode, symbol, answers[p.Answer]))
					response = append(response, "")
					continue
				}

				response = append(response, i18n.Translatef(lang.QuizDetailResultIncorrectAnswerLine, languageCode, symbol, answers[p.Answer], answers[0]))
				response = append(response, "")
			}

			_ = ctx.Reply(strings.Join(response, "\n"))
		}
		return true
	}
}

func (*QuizResultDetail) NeedRightsLevel() int {
	return role.RightsLevel(role.Teacher{})
}

func (*QuizResultDetail) Endpoint() string {
	return "quiz_result_detail"
}
