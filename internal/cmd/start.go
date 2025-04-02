package cmd

import (
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/repository/repository_model"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"time"
)

type Start struct {
	log *slog.Logger
}

func (s *Start) Run(b *tele.Bot, ctx tele.Context, args []string) error {
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

		if err := repository.Repo().SaveNewUser(repository_model.NewUser(id, sec.Name, sec.Surname, sec.Patronymic, sec.Role)); err != nil {
			s.log.Error("error saving new user", slog.Any("err", err))

			_ = repository.Repo().SaveNewSecret(sec)

			return ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, "error saving new user"))
		}

		r, _ := role.FromName(sec.Role)

		return ctx.Send(i18n.Translatef(lang.InvitationKeyApplied, languageCode, sec.Name, sec.Patronymic, i18n.Translatef(r.Translation(), languageCode)))
	}

	selector := &tele.ReplyMarkup{}

	quizCreateBtn := selector.Data(i18n.Translatef(lang.QuizCreateBtn, languageCode), "quiz_create")

	if u.Role == (role.Admin{}).Name() {
		admInvitationCreateBtn := selector.Data(i18n.Translatef(lang.InvitationKeyCreateBtn, languageCode), "adm_inv_create")

		selector.Inline(selector.Row(admInvitationCreateBtn, quizCreateBtn))

		if _, err := b.Reply(ctx.Message(), i18n.Translatef(lang.AdminPanelTitle, languageCode, u.Name, u.Patronymic), selector); err != nil {
			return err
		}

		return nil
	}

	//sendMenu

	return ctx.Send("Hello World!")
}

func (s *Start) Endpoint() string {
	return "start"
}
