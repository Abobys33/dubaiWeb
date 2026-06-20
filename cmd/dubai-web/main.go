package main

import (
	"context"
	"dubaiWeb/internal/config"
	"dubaiWeb/internal/http-server/handlers/feedback"
	"dubaiWeb/internal/lib/mwlogger"
	"dubaiWeb/internal/lib/sl"
	"dubaiWeb/internal/mailer"
	"dubaiWeb/internal/storage/postgres"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting event booker", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	storage, err := postgres.InitDB(&cfg.Database)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	mailerSvc := mailer.New(cfg.SMTP)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	fs := http.FileServer(http.Dir("./web/static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.html", http.StatusFound)
	})

	router.Post("/api/feedback", feedback.New(log, storage, mailerSvc))

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Канал ловит сигналы ОС, читает его только main.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Отдельный канал для остановки фоновой горутины,
	// чтобы не было гонки за чтение из stop.
	done := make(chan struct{})

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	select {
	case sign := <-stop:
		log.Info("application stopping", slog.String("signal", sign.String()))
	case err := <-serverErr:
		if err != nil {
			log.Error("failed to start server", sl.Err(err))
		}
	}

	close(done)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", sl.Err(err))
	}

	log.Info("application stopped")

	if err := storage.Close(); err != nil {
		log.Error("failed to close postgres connection", sl.Err(err))
	}

	log.Info("postgres connection closed")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
