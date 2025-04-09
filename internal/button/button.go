package button

import tele "gopkg.in/telebot.v4"

type Button interface {
	Run(tele.Context, []string) error
	NeedRightsLevel() int
	Endpoint() string
}
