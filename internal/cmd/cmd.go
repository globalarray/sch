package cmd

import tele "gopkg.in/telebot.v4"

type Command interface {
	Run(tele.Context, []string) error
	Endpoint() string
}
