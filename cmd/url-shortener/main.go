package main

import (
	"log/slog"
	"os"
	"restapi/URL-Shortener/internal/config"
	"restapi/URL-Shortener/internal/lib/logger/sl"
	"restapi/URL-Shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage
	// _, err = storage.SaveURL("www.google.com", "8.8.8.8")

	// if err != nil {
	// 	slog.Error("falied to add alias", sl.Err(err))
	// }

	// url, err := storage.GetURL("8.8.8.8")

	// if err != nil {
	// 	slog.Error("failed to get alias", sl.Err(err))
	// }

	// fmt.Println(url)

	// id, err := storage.DelURLByAlias("8.8.8.8")

	// if err != nil {
	// 	slog.Error("failed to del url by alias", sl.Err(err))
	// }

	// fmt.Println(id)

	// TODO: init router: chi, "chi render"

	// TODO: run server:
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		)
	}

	return log
}
