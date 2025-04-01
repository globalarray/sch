package main

import (
	"benzo/internal/app"
	"benzo/pkg/i18n"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"syscall"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	cfg, err := readConfig()

	if err != nil {
		panic(err)
	}

	a := cfg.New()

	log := a.Log().With(slog.String("level", "main"))

	wd, err := os.Getwd()

	if err != nil {
		log.Error("error getting wd", err)
		return
	}

	if err := i18n.LoadLangs(path.Join(wd, "static", "languages")); err != nil {
		log.Error("error loading languages", err)
		return
	}

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan

		log.Warn("Shutting down...")

		a.Shutdown()
	}()

	if err := a.Run(); err != nil {
		log.Error("error running app:", err)
	}
}

func readConfig() (cfg app.Config, err error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
