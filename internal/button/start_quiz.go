package button

import (
	"benzo/internal/service"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"strconv"
)

type StartQuiz struct {
	log *slog.Logger
	StudentButton
}

func (b *StartQuiz) Run(ctx tele.Context, args []string) error {
	if len(args) != 1 {
		return ErrInvalidUsage
	}

	id := ctx.Callback().Sender.ID
	languageCode := ctx.Callback().Sender.LanguageCode

	quizID, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	return service.Quiz().ProcessQuiz(ctx, quizID, id, languageCode)
}

func (*StartQuiz) Endpoint() string {
	return "start_quiz"
}
