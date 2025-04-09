package button

import (
	"benzo/internal/callback"
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/repository/repository_model"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	"github.com/samber/lo"
	tele "gopkg.in/telebot.v4"
	"log/slog"
)

type CreateInvitationKey struct {
	log *slog.Logger
}

func (b *CreateInvitationKey) Run(ctx tele.Context, _ []string) error {
	languageCode := ctx.Callback().Sender.LanguageCode
	id := ctx.Callback().Sender.ID

	key := lo.RandomString(16, lo.LettersCharset)

	sec := repository_model.NewSecretDefault(key)
	sec.CreatedBy = id

	if err := repository.Repo().SaveNewSecret(sec); err != nil {
		return err
	}

	callback.Subscribe(id, fillPersonalDataKeyCallback(key))

	if err := ctx.Send(i18n.Translatef(lang.InvitationKeyCreateTitle, languageCode, key)); err != nil {
		return err
	}

	return ctx.Send(i18n.Translatef(lang.InvitationKeyCreateName, languageCode))
}

func (b *CreateInvitationKey) NeedRightsLevel() int {
	return role.RightsLevel(role.Admin{})
}

func (b *CreateInvitationKey) Endpoint() string {
	return "adm_inv_create"
}
