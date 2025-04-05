package button

import (
	"benzo/internal/lang"
	"benzo/internal/quiz"
	"benzo/internal/repository"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strconv"
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

	var response []string

	response = append(response, i18n.Translatef(lang.QuizInfoTitle, languageCode, q.Name))
	response = append(response, i18n.Translatef(lang.QuizInfoInvitationLinkLine, languageCode, encoded))
	response = append(response, i18n.Translatef(lang.QuizInfoQuestionsCountLine, languageCode, len(questions)))
	response = append(response, "\n")
	
	return nil
}

func (b *GetInfoQuiz) NeedRightsLevel() int {
	return role.RightsLevel(role.Teacher{})
}

func (b *GetInfoQuiz) Endpoint() string {
	return "get_info_quiz"
}
