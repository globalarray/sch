package button

import (
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	"errors"
	tele "gopkg.in/telebot.v4"
	"log/slog"
)

type RemoveInvitationKey struct {
	log *slog.Logger
}

var ErrInvalidKey = errors.New("invalid key")

func (b *RemoveInvitationKey) Run(_ *tele.Bot, ctx tele.Context, args []string) error {
	if len(args) != 1 {
		return ErrInvalidUsage
	}

	languageCode := ctx.Callback().Sender.LanguageCode

	sec, err := repository.Repo().GetSecretByKey(args[0])

	if err != nil {
		return err
	}

	if sec.Key != args[0] {
		return ErrInvalidKey
	}

	if err := repository.Repo().RemoveSecretByKey(sec.Key); err != nil {
		return err
	}

	return ctx.Send(i18n.Translatef(lang.InvitationKeyRemoved, languageCode, sec.Key))
}

func (b *RemoveInvitationKey) NeedRightsLevel() int {
	return role.RightsLevel(role.Admin{})
}

func (b *RemoveInvitationKey) Endpoint() string {
	return "adm_inv_remove"
}
