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
	"github.com/prayogopangestu/crm-system/backend/internal/modules/analytics"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/contact"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/deal"
	integrationmodule "github.com/prayogopangestu/crm-system/backend/internal/modules/integration"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/notification"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/pipeline"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/search"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/task"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/user"
	"github.com/prayogopangestu/crm-system/backend/internal/platform/postgresx"
	"github.com/prayogopangestu/crm-system/backend/internal/platform/redisx"
	"github.com/prayogopangestu/crm-system/backend/internal/server/grpcserver"
	"github.com/prayogopangestu/crm-system/backend/internal/server/httpserver"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
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
	location := mustLocation(cfg.App.Timezone)
	store := postgresx.New(pool, location)
	defer store.Close()

	var cache *redisx.Cache
	cache, err = redisx.New(cfg.Redis.URL)
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
	telegramClient := integrationmodule.NewTelegramClient()
	cacheHelper := shared.CacheHelper{Cache: cache, Logger: log}

	userService := user.NewService(user.NewRepository(store), cacheHelper, tokenManager, cfg.App.BaseURL, cfg.Auth.BcryptCost)
	contactService := contact.NewService(contact.NewRepository(store), cacheHelper)
	dealService := deal.NewService(deal.NewRepository(store), cacheHelper)
	taskService := task.NewService(task.NewRepository(store), cacheHelper, location)
	analyticsService := analytics.NewService(analytics.NewRepository(store), cacheHelper, location)
	pipelineService := pipeline.NewService(pipeline.NewRepository(store), cacheHelper)
	integrationRepository := integrationmodule.NewRepository(store)
	integrationService := integrationmodule.NewService(integrationRepository, cipher, telegramClient)
	notificationService := notification.NewService(notification.NewRepository(store))
	searchService := search.NewService(search.NewRepository(store), cacheHelper)

	worker := integrationmodule.NewWorker(
		integrationRepository, cipher, telegramClient, log,
		cfg.Telegram.WorkerInterval, cfg.Telegram.WorkerBatchSize,
	)
	go worker.Run(ctx)

	modules := httpserver.Modules{
		User:         user.NewHandler(userService, log),
		Contact:      contact.NewHandler(contactService, log),
		Deal:         deal.NewHandler(dealService, log),
		Task:         task.NewHandler(taskService, log),
		Analytics:    analytics.NewHandler(analyticsService, log),
		Pipeline:     pipeline.NewHandler(pipelineService, log),
		Integration:  integrationmodule.NewHandler(integrationService, log),
		Notification: notification.NewHandler(notificationService, log),
		Search:       search.NewHandler(searchService, log),
	}
	httpServer := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           httpserver.Router(modules, tokenManager, cache, log, cfg.HTTP.AllowedOrigins, func() error { return store.Ping(context.Background()) }),
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
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcserver.AuthInterceptor(tokenManager)))
	userv1.RegisterUserServiceServer(grpcServer, user.NewGRPCServer(userService))
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
