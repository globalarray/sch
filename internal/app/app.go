package app

import (
	"benzo/internal/button"
	"benzo/internal/cmd"
	"benzo/internal/handlers"
	"benzo/internal/repository"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type App struct {
	cfg Config
	log *slog.Logger
	b   *tele.Bot
}

func (a *App) Run() error {
	mysqlData := map[string]string{
		EnvMySQLLogin:    "",
		EnvMySQLHostname: "",
		EnvMySQLPort:     "",
		EnvMySQLDatabase: "",
		EnvMySQLPassword: "",
	}

	for idx := range mysqlData {
		k, ok := os.LookupEnv(idx)

		if !ok {
			return fmt.Errorf("environment variable %s not found", idx)
		}

		mysqlData[idx] = k
	}

	mysqlPort, err := strconv.Atoi(mysqlData[EnvMySQLPort])

	if err != nil {
		return fmt.Errorf("getting incorrect MySQL port")
	}

	if _, err := repository.New(mysqlData[EnvMySQLLogin], mysqlData[EnvMySQLPassword], mysqlData[EnvMySQLHostname], mysqlData[EnvMySQLDatabase], uint16(mysqlPort)); err != nil {
		return fmt.Errorf("error initializing MySQL connection: %s", err.Error())
	}

	token, ok := os.LookupEnv(EnvTelegramToken)

	if !ok {
		return fmt.Errorf("environment variable %s not found", EnvTelegramToken)
	}

	pref := tele.Settings{
		Token: token,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
		},
		ParseMode: tele.ModeHTML,
	}

	b, err := tele.NewBot(pref)

	if err != nil {
		return err
	}

	a.b = b

	_ = cmd.NewManager(a.log)

	_ = button.NewManager(a.log)

	_ = handlers.NewTeleHandler(b, a.log)

	a.log.Info("Bot is now running.  Press CTRL-C to exit.")

	b.Start()

	return nil
}

func (a *App) Log() *slog.Logger {
	return a.log
}

func (a *App) Shutdown() {
	if a.b != nil {
		a.b.Stop()
	}
}
