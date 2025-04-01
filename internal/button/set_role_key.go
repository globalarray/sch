package button

import (
	"benzo/internal/lang"
	"benzo/internal/repository"
	"benzo/internal/user/role"
	"benzo/pkg/i18n"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strings"
	"time"
)

type SetRoleKey struct {
	log *slog.Logger
}

func (b *SetRoleKey) Run(_ *tele.Bot, ctx tele.Context, args []string) error {
	if len(args) < 1 {
		return ErrInvalidUsage
	}

	languageCode := ctx.Callback().Sender.LanguageCode

	sec, err := repository.Repo().GetSecretByKey(args[0])

	if err != nil {
		return err
	}

	selector := &tele.ReplyMarkup{}

	if len(args) > 1 {
		roleName := args[1]

		r, err := role.FromName(roleName)

		if err != nil {
			return err
		}

		if err := repository.Repo().UpdateSecretRole(args[0], r.Name()); err != nil {
			return err
		}

		ownerFullName := strings.Join([]string{sec.Name, sec.Patronymic, sec.Surname}, " ")

		removeBtn := selector.Data(i18n.Translatef(lang.InvitationKeyRemoveBtn, languageCode), fmt.Sprintf("adm_inv_remove-%s", args[0]))

		selector.Inline(selector.Row(removeBtn))

		return ctx.Send(i18n.Translatef(lang.InvitationKeyCreated, languageCode, ownerFullName, i18n.Translatef(r.Translation(), languageCode), sec.Expiration.Format(time.DateTime)), selector)
	}

	uniqFormat := fmt.Sprintf(b.Endpoint()+"-%s", args[0]) + "-%s"

	var buttons []tele.Btn

	for _, r := range []role.Role{role.Student{}, role.Teacher{}, role.Admin{}} {
		buttons = append(buttons, selector.Data(i18n.Translatef(r.Translation(), languageCode), fmt.Sprintf(uniqFormat, r.Name())))
	}

	selector.Inline(selector.Row(buttons...))

	return ctx.Send(i18n.Translatef(lang.InvitationKeyStageRoleMessage, languageCode), selector)
}

func (b *SetRoleKey) NeedRightsLevel() int {
	return role.RightsLevel(role.Admin{})
}

func (b *SetRoleKey) Endpoint() string {
	return "adm_inv_set_role"
}
