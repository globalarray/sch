package handlers

import (
	"benzo/internal/button"
	"benzo/internal/callback"
	"benzo/internal/cmd"
	"benzo/internal/repository"
	"benzo/internal/user/role"
	"errors"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strings"
)

type TeleHandler struct {
	b   *tele.Bot
	log *slog.Logger
}

var (
	studentRightsLevel = role.RightsLevel(role.Student{})

	ErrRightsLevel = errors.New("the role rights level does not allow you to perform this action")
)

func NewTeleHandler(b *tele.Bot, log *slog.Logger) *TeleHandler {
	h := &TeleHandler{b: b, log: log.With(slog.String("level", "handlers/teleHandler"))}

	b.Handle(tele.OnText, h.onText)
	b.Handle(tele.OnCallback, h.onCallback)

	return h
}

func (h *TeleHandler) onText(ctx tele.Context) error {
	id := ctx.Message().Sender.ID

	if callback.Exists(id) {
		callback.Call(ctx)

		return nil
	}

	if ctx.Text()[0] == '/' {
		cmdData := strings.Split(ctx.Text(), " ")

		if c, ok := cmd.Mgr().Get(cmdData[0][1:]); ok {
			return c.Run(ctx, cmdData[1:])
		}
	}

	return nil
}

func (h *TeleHandler) onCallback(ctx tele.Context) error {
	args := strings.Split(strings.TrimSpace(ctx.Callback().Data), "-")

	id := ctx.Callback().Sender.ID

	if b, ok := button.Mgr().Get(args[0]); ok {
		if err := h.b.Delete(ctx.Callback().Message); err != nil {
			return err
		}

		if b.NeedRightsLevel() > studentRightsLevel {
			u, err := repository.Repo().GetUserByTelegramID(id)

			if err != nil {
				return err
			}

			r, err := role.FromName(u.Role)

			if err != nil {
				return err
			}

			if role.RightsLevel(r) < b.NeedRightsLevel() {
				return ErrRightsLevel
			}
		}

		return b.Run(ctx, args[1:])
	}

	return nil
}
