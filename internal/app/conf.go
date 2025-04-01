package app

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
)

const (
	EnvTelegramToken = "TG_TOKEN"
	EnvMySQLLogin    = "MYSQL_LOGIN"
	EnvMySQLHostname = "MYSQL_HOSTNAME"
	EnvMySQLPort     = "MYSQL_PORT"
	EnvMySQLDatabase = "MYSQL_DATABASE"
	EnvMySQLPassword = "MYSQL_PWD"
)

const (
	LoggerLevelDebug = "debug"
	LoggerLevelInfo  = "info"
	LoggerLevelWarn  = "warn"
	LoggerLevelError = "error"
)

type (
	DataBase struct {
		Login    string
		Hostname string
		Port     uint16
		Database string
		Password string
	}
	Config struct {
		Logger struct {
			Level      string
			TimeFormat string
		}
	}
)

func (c Config) New() *App {
	var logLevel slog.Level

	switch c.Logger.Level {
	case LoggerLevelDebug:
		logLevel = slog.LevelDebug
	case LoggerLevelWarn:
		logLevel = slog.LevelWarn
	case LoggerLevelError:
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: c.Logger.TimeFormat,
		Level:      logLevel,
	}))

	slog.SetDefault(logger)

	return &App{cfg: c, log: logger}
}
