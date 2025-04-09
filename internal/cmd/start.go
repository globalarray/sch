package cmd

import (
	"benzo/internal/lang"
	"benzo/internal/quiz"
	"benzo/internal/repository"
	"benzo/internal/repository/repository_model"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"time"
)

type Start struct {
	log *slog.Logger
}

func (s *Start) Run(ctx tele.Context, args []string) error {
	id := ctx.Sender().ID
	languageCode := ctx.Sender().LanguageCode

	u, err := repository.Repo().GetUserByTelegramID(id)

	if err != nil {
		return err
	}

	if u.TelegramID != id {
		if len(args) == 0 {
			return ctx.Send(i18n.Translatef(lang.InvitationKeyRequired, languageCode))
		}

		sec, err := repository.Repo().GetSecretByKey(args[0])

		if err != nil {
			return err
		}

		if sec.Key != args[0] {
			return ctx.Send(i18n.Translatef(lang.InvitationKeyInvalid, languageCode))
		}

		if time.Now().Unix() > sec.Expiration.Unix() {
			_ = repository.Repo().RemoveSecretByKey(sec.Key)

			return ctx.Send(i18n.Translatef(lang.InvitationKeyInvalid, languageCode))
		}

		if err := repository.Repo().RemoveSecretByKey(sec.Key); err != nil {
			s.log.Error("error removing secret", slog.Any("err", err))

			return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode))
		}

		u := repository_model.NewUser(id, sec.Name, sec.Surname, sec.Patronymic, sec.Role)

		if err := repository.Repo().SaveNewUser(u); err != nil {
			s.log.Error("error saving new user", slog.Any("err", err))

			_ = repository.Repo().SaveNewSecret(sec)

			return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, "error saving new user"))
		}

		r, _ := role.FromName(sec.Role)

		return ctx.Send(i18n.Translatef(lang.InvitationKeyApplied, languageCode, u.PrettyName(), i18n.Translatef(r.Translation(), languageCode)))
	}

	selector := &tele.ReplyMarkup{}

	if len(args) == 1 {
		quizID, err := quiz.Decode(args[0])

		if err != nil {
			return err
		}

		q, err := repository.Repo().GetQuizByID(quizID)

		if err != nil {
			return err
		}

		if q.ID != quizID {
			return ctx.Reply(i18n.Translatef(lang.QuizNotFound, languageCode))
		}

		selector.Inline(selector.Row(selector.Data(i18n.Translatef(lang.QuizStartBtn, languageCode), fmt.Sprintf("start_quiz-%d", quizID))))

		return ctx.Reply(i18n.Translatef(lang.QuizStartMessage, languageCode, u.PrettyName(), q.Name), selector)
	}

	quizCreateBtn := selector.Data(i18n.Translatef(lang.QuizCreateBtn, languageCode), "quiz_create")
	quizListBtn := selector.Data(i18n.Translatef(lang.QuizListBtn, languageCode), "quiz_list")

	teacherButtonsRow := selector.Row(quizListBtn, quizCreateBtn)

	if u.Role == (role.Teacher{}).Name() {
		selector.Inline(teacherButtonsRow)

		return ctx.Reply(i18n.Translatef(lang.TeacherPanelTitle, languageCode, u.PrettyName()), selector)
	}

	if u.Role == (role.Admin{}).Name() {
		admInvitationCreateBtn := selector.Data(i18n.Translatef(lang.InvitationKeyCreateBtn, languageCode), "adm_inv_create")

		selector.Inline(teacherButtonsRow, selector.Row(admInvitationCreateBtn))

		return ctx.Reply(i18n.Translatef(lang.AdminPanelTitle, languageCode, u.PrettyName()), selector)
	}

	//sendMenu

	return ctx.Send("TODO")
}

func (s *Start) Endpoint() string {
	return "start"
}
