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
	maxQuizNameLength int = 256

	minQuestionLength int = 5
	maxQuestionLength int = 256

	minQuestionAnswerLength int = 2
	maxQuestionAnswerLength int = 100

	minQuestionAnswersCount int = 2
	maxQuestionAnswersCount int = 5
)

func (b *CreateQuiz) Run(ctx tele.Context, _ []string) error {
	id := ctx.Callback().Sender.ID
	languageCode := ctx.Callback().Sender.LanguageCode

	callback.Subscribe(id, b.createQuizCallback)

	return ctx.Send(i18n.Translatef(lang.QuizCreateStartMessage, languageCode))
}

func (b *CreateQuiz) createQuizCallback(ctx tele.Context) bool {
	languageCode := ctx.Message().Sender.LanguageCode
	id := ctx.Message().Sender.ID
	quizName := ctx.Message().Text

	if len(quizName) < minQuizNameLength || len(quizName) > maxQuizNameLength {
		_ = ctx.Reply(i18n.Translatef(lang.QuizCreateNameLengthInvalid, languageCode))

		return false
	}

	quizID, err := repository.Repo().SaveNewQuiz(repository_model.NewQuiz(quizName, id))

	if err != nil {
		_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))

		return false
	}

	callback.Subscribe(id, createQuestionCallback(quizID))

	_ = ctx.Reply(i18n.Translatef(lang.QuizCreatedMessage, languageCode))

	return false
}

func createQuestionCallback(quizID int64) callback.CallbackFunc {
	return func(ctx tele.Context) bool {
		id := ctx.Message().Sender.ID
		languageCode := ctx.Message().Sender.LanguageCode
		question := ctx.Message().Text

		quiz, err := repository.Repo().GetQuizByID(quizID)

		if err != nil {
			_ = ctx.Reply(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))
			return true
		}

		if quiz.ID != quizID {
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

		callback.Subscribe(id, addAnswersToQuestionCallback(quizID, questionID))

		_ = ctx.Reply(i18n.Translatef(lang.QuizQuestionCreatedMessage, languageCode))

		return false
	}
}

func addAnswersToQuestionCallback(quizID, questionID int64) callback.CallbackFunc {
	return func(ctx tele.Context) bool {
		languageCode := ctx.Message().Sender.LanguageCode

		formattedAns := strings.Replace(ctx.Message().Text, "-", "–", -1)

		answers := strings.Split(formattedAns, ";")

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

		addNewQuestionBtn := selector.Data(i18n.Translatef(lang.QuizAddNewQuestionBtn, languageCode), fmt.Sprintf("add_new_question_quiz-%d", quizID))
		deleteQuestionBtn := selector.Data(i18n.Translatef(lang.QuizQuestionRemoveBtn, languageCode), fmt.Sprintf("remove_question_quiz-%d", questionID))

		rows := []tele.Row{selector.Row(addNewQuestionBtn, deleteQuestionBtn)}

		if len(questions) > 1 {
			rows = append(rows, selector.Row(selector.Data(i18n.Translatef(lang.QuizAddingQuestionsStopBtn, languageCode), fmt.Sprintf("get_info_quiz-%d", quizID))))
		}

		selector.Inline(rows...)

		_ = ctx.Reply(i18n.Translatef(lang.QuizQuestionAnswersAdded, languageCode), selector)

		return true
	}
}

func (b *CreateQuiz) NeedRightsLevel() int {
	return role.RightsLevel(role.Teacher{})
}

func (b *CreateQuiz) Endpoint() string {
	return "quiz_create"
}
