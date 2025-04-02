package button

import (
	"benzo/internal/callback"
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/repository/repository_model"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strings"
)

type CreateQuiz struct {
	log *slog.Logger
}

const (
	minQuizNameLength int = 5
	maxQuizNameLength int = 20

	minQuestionLength int = 5
	maxQuestionLength int = 50

	minQuestionAnswerLength int = 5
	maxQuestionAnswerLength int = 20

	minQuestionAnswersCount int = 2
	maxQuestionAnswersCount int = 5
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

	quizID, err := repository.Repo().SaveNewQuiz(repository_model.NewQuiz(quizName, id))

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

		quiz, err := repository.Repo().GetQuizByID(quizID)

		if err != nil {
			_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))
			return true
		}

		if quiz.ID != id {
			_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode, "quiz does not exist"))
			return true
		}

		if len(question) < minQuestionLength || len(question) > maxQuestionLength {
			_ = ctx.Reply(i18n.Translatef(lang.QuizQuestionNameLengthInvalid, languageCode))
			return false
		}

		questionID, err := repository.Repo().SaveNewQuestion(repository_model.NewQuestion(quizID, question, "", ""))

		if err != nil {
			_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))

			return true
		}

		callback.Subscribe(id, b.addAnswersToQuestionCallback(quizID, questionID))

		_ = ctx.Reply(i18n.Translatef(lang.QuizQuestionCreatedMessage, languageCode))

		return true
	}
}

func (b *CreateQuiz) addAnswersToQuestionCallback(quizID, questionID int64) callback.CallbackFunc {
	return func(_ *tele.Bot, ctx tele.Context) bool {
		id := ctx.Message().Sender.ID
		languageCode := ctx.Message().Sender.LanguageCode

		answers := strings.Split(ctx.Message().Text, ";")

		if len(answers) < minQuestionAnswersCount || len(answers) > maxQuestionAnswersCount {
			_ = ctx.Reply(i18n.Translatef(lang.QuizQuestionAnswersCountInvalid, languageCode))

			return false
		}

		for _, answer := range answers {
			if len(answer) < minQuestionAnswerLength || len(answer) > maxQuestionAnswerLength {
				_ = ctx.Reply(i18n.Translatef(lang.QuizQuestionAnswerInvalidLength, languageCode, answer))

				return false
			}
		}

		questions, err := repository.Repo().GetQuestionsByQuizID(quizID)

		if err != nil {
			_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))

			return false
		}

		if err := repository.Repo().UpdateQuestionAnswers(questionID, answers); err != nil {
			_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))

			return false
		}

		selector := &tele.ReplyMarkup{}

		buttons := []tele.Btn{selector.Data(i18n.Translatef(lang.QuizAddNewQuestionBtn, languageCode), fmt.Sprintf("quiz_add_new_question-%d", quizID))}

		if len(questions) > 0 {
			buttons = append(buttons, selector.Data(i18n.Translatef(lang.QuizAddingQuestionsStopBtn, languageCode), fmt.Sprintf("quiz_adding_buttons_stop-%d", quizID)))
		}

	}
}

func (b *CreateQuiz) NeedRightsLevel() int {
	return role.RightsLevel(role.Teacher{})
}

func (b *CreateQuiz) Endpoint() string {
	return "quiz_create"
}
