package button

import (
	"benzo/internal/callback"
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	"errors"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strings"
)

const (
	RequiredDataLen int = 3
)

var (
	ErrInvalidUsage = errors.New("invalid usage")
)

type FillPersonalData struct {
	log *slog.Logger
}

func (b *FillPersonalData) Run(_ *tele.Bot, ctx tele.Context, args []string) error {
	if len(args) < 1 {
		return ErrInvalidUsage
	}

	id := ctx.Callback().Sender.ID
	languageCode := ctx.Callback().Sender.LanguageCode

	callback.Subscribe(id, fillPersonalDataKeyCallback(args[0]))

	return ctx.Send(i18n.Translatef(lang.InvitationKeyNameRefillMessage, languageCode))
}

func (b *FillPersonalData) NeedRightsLevel() int {
	return role.RightsLevel(role.Admin{})
}

func (b *FillPersonalData) Endpoint() string {
	return "adm_inv_refill"
}

func fillPersonalDataKeyCallback(key string) callback.CallbackFunc {
	return func(bot *tele.Bot, ctx tele.Context) bool {
		languageCode := ctx.Message().Sender.LanguageCode
		data := strings.Split(ctx.Message().Text, " ")

		if len(data) != RequiredDataLen {
			_ = ctx.Send(i18n.Translatef(lang.InvitationKeyNameIncorrect, languageCode))

			return false
		}

		if err := repository.Repo().UpdateSecretPersonalData(key, data[0], data[1], data[2]); err != nil {
			_ = ctx.Send(i18n.Translatef(lang.RuntimeError, languageCode, err.Error()))

			return false
		}

		selector := &tele.ReplyMarkup{}

		refillBtn := selector.Data(i18n.Translatef(lang.No, languageCode), fmt.Sprintf("adm_inv_refill-%s", key))
		nextStageBtn := selector.Data(i18n.Translatef(lang.Yes, languageCode), fmt.Sprintf("adm_inv_set_role-%s", key))

		selector.Inline(selector.Row(nextStageBtn, refillBtn))

		_, _ = bot.Reply(ctx.Message(), i18n.Translatef(lang.InvitationKeyNameSaved, languageCode, data[0], data[1], data[2]), selector)

		return true
	}
}
