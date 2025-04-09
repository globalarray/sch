package callback

import (
	tele "gopkg.in/telebot.v4"
	"sync"
)

type CallbackFunc func(ctx tele.Context) bool

var (
	callbacksMu sync.Mutex
	callbacks   = map[int64]CallbackFunc{}
)

func Subscribe(id int64, fn CallbackFunc) {
	callbacksMu.Lock()
	defer callbacksMu.Unlock()

	delete(callbacks, id)

	callbacks[id] = fn
}

func Exists(id int64) bool {
	callbacksMu.Lock()
	defer callbacksMu.Unlock()
	return callbacks[id] != nil
}

func Call(ctx tele.Context) {
	id := ctx.Message().Sender.ID

	if fn, ok := callbacks[id]; ok {
		if ok := fn(ctx); ok {
			delete(callbacks, id)
		}
	}
}
