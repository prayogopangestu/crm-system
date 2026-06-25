package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	userv1 "github.com/prayogopangestu/crm-system/backend/api/protobuf/gen"
	"github.com/prayogopangestu/crm-system/backend/internal/config"
	grpcdelivery "github.com/prayogopangestu/crm-system/backend/internal/delivery/grpc"
	httpdelivery "github.com/prayogopangestu/crm-system/backend/internal/delivery/http"
	"github.com/prayogopangestu/crm-system/backend/internal/integration/telegram"
	postgresrepo "github.com/prayogopangestu/crm-system/backend/internal/repository/postgres"
	redisrepo "github.com/prayogopangestu/crm-system/backend/internal/repository/redis"
	"github.com/prayogopangestu/crm-system/backend/internal/usecase"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
	"github.com/prayogopangestu/crm-system/backend/pkg/cryptoutil"
	"github.com/prayogopangestu/crm-system/backend/pkg/database"
	"github.com/prayogopangestu/crm-system/backend/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}
	log := logger.New(cfg.App.LogLevel)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := database.OpenPostgres(ctx, cfg.Database.URL, cfg.Database.MinConns, cfg.Database.MaxConns)
	if err != nil {
		log.Error("postgres startup failed", "error", err)
		os.Exit(1)
	}
	repository := postgresrepo.New(pool, mustLocation(cfg.App.Timezone))
	defer repository.Close()

	var cache *redisrepo.Cache
	cache, err = redisrepo.New(cfg.Redis.URL)
	if err != nil {
		log.Warn("redis configuration invalid; cache disabled", "error", err)
		cache = nil
	} else if err := cache.Ping(ctx); err != nil {
		log.Warn("redis unavailable; cache will fail open", "error", err)
	}
	if cache != nil {
		defer cache.Close()
	}

	tokenManager := auth.New(cfg.Auth.JWTSecret, cfg.Auth.JWTTTL)
	cipher, err := cryptoutil.New(cfg.Security.EncryptionKey)
	if err != nil {
		log.Error("encryption setup failed", "error", err)
		os.Exit(1)
	}
	telegramClient := telegram.New()
	service := usecase.New(
		usecase.Repositories{
			Users: repository, Contacts: repository, Deals: repository, Tasks: repository,
			Analytics: repository, Pipeline: repository, Integrations: repository,
			Search: repository, Notifications: repository,
		},
		cache, tokenManager, cipher, telegramClient, log,
		mustLocation(cfg.App.Timezone), cfg.App.BaseURL, cfg.Auth.BcryptCost,
	)
	worker := usecase.NewOutboxWorker(
		repository, repository, cipher, telegramClient, log,
		cfg.Telegram.WorkerInterval, cfg.Telegram.WorkerBatchSize,
	)
	go worker.Run(ctx)

	handler := httpdelivery.NewHandler(httpdelivery.Services{
		Users: service, Contacts: service, Deals: service, Tasks: service,
		Analytics: service, Settings: service, Search: service, Notifications: service,
	}, log)
	httpServer := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           httpdelivery.Router(handler, tokenManager, cache, log, cfg.HTTP.AllowedOrigins, func() error { return repository.Ping(context.Background()) }),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	grpcListener, err := net.Listen("tcp", cfg.GRPC.Addr)
	if err != nil {
		log.Error("grpc listen failed", "error", err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcdelivery.AuthInterceptor(tokenManager)))
	userv1.RegisterUserServiceServer(grpcServer, grpcdelivery.NewUserServer(service))
	if cfg.App.Env != "production" {
		reflection.Register(grpcServer)
	}

	go func() {
		log.Info("HTTP server started", "addr", cfg.HTTP.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server failed", "error", err)
			stop()
		}
	}()
	go func() {
		log.Info("gRPC server started", "addr", cfg.GRPC.Addr)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Error("gRPC server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Info("shutdown started")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(shutdownCtx)
	grpcServer.GracefulStop()
	log.Info("shutdown complete")
}

func mustLocation(value string) *time.Location {
	location, err := time.LoadLocation(value)
	if err != nil {
		panic(err)
	}
	return location
}
