// Package main — точка входа в приложение.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/bublik008/threat-exchange/internal/api"      
	"github.com/bublik008/threat-exchange/internal/api/gen"   
	"github.com/bublik008/threat-exchange/internal/config"
	"github.com/bublik008/threat-exchange/internal/repository/postgres"
	"github.com/bublik008/threat-exchange/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}


	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to create postgres pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping postgres", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to postgres successfully")


	repo := postgres.NewRepository(pool)
	svc := service.NewThreatService(repo)
	server := api.NewServer(svc)

	mux := http.NewServeMux()
	
	handler := gen.HandlerFromMux(server, mux)

	handler = api.LoggingMiddleware(handler)


	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: handler,
		// Для защиты от медленных клиентов
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("starting HTTP server", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server exited properly")
}