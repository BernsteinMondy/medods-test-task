package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/medods-test-task/src/internal/authorization"
	"github.com/BernsteinMondy/medods-test-task/src/internal/hasher"
	"github.com/BernsteinMondy/medods-test-task/src/internal/httpserver"
	"github.com/BernsteinMondy/medods-test-task/src/internal/repository"
	"github.com/BernsteinMondy/medods-test-task/src/internal/service"
	"github.com/BernsteinMondy/medods-test-task/src/pkg/database"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	err := run()
	if err != nil {
		println("run() returned error: " + err.Error())
	}
}

func run() (err error) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	defer cancel()

	slog.Info("Loading config...")
	cfg, err := loadConfigFromEnv()
	if err != nil {
		return fmt.Errorf("loadConfigFromEnv() returned error: %w", err)
	}
	slog.Info("Config loaded")

	select {
	case <-ctx.Done():
	default:
	}

	slog.Info("Connecting to database")
	db, err := newDB(cfg.DB)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	slog.Info("Database connected")
	defer func() {
		slog.Info("Close connection with database")
		closeErr := db.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close database connection: %w", err))
		}
		slog.Info("Connection closed")
	}()

	repo := repository.New(db)

	hasherService := hasher.NewService()
	tokenService := authorization.NewTokenService(cfg.TokenService.SecretKey)

	srvc := service.NewService(
		repo,
		hasherService,
		tokenService,
	)

	app := httpserver.NewFiber(srvc)

	stopWg := &sync.WaitGroup{}
	stopWg.Add(1)
	go func(ctx context.Context) {
		defer stopWg.Done()
		err = launchFiberServer(ctx, app, cfg.HTTPServer.ListenAddr)
		if err != nil {
			slog.Error("launch fiber error: ", slog.String("error", err.Error()))
		}
	}(ctx)

	<-ctx.Done()
	stopWg.Wait()
	return nil
}

func newDB(dbCfg DB) (*sql.DB, error) {
	c := database.Config{
		Host:     dbCfg.Host,
		Port:     dbCfg.Port,
		User:     dbCfg.User,
		Password: dbCfg.Password,
		DBName:   dbCfg.DBName,
		SSLMode:  dbCfg.SSLMode,
	}

	db, err := database.NewSQL(c)
	if err != nil {
		return nil, fmt.Errorf("database new sql: %w", err)
	}

	return db, nil
}

func launchFiberServer(ctx context.Context, f *fiber.App, listenAddr string) (err error) {
	var fiberShutdownErr error
	defer func() {
		err = errors.Join(err, fiberShutdownErr)
	}()

	shutdownDone := make(chan struct{})
	go func(ctx context.Context) {
		<-ctx.Done()

		slog.Info("Shutting down fiber")
		fiberShutdownErr = f.ShutdownWithTimeout(time.Second * 6)
		slog.Info("Fiber shut down")

		close(shutdownDone)
	}(ctx)

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	slog.Info("Starting to listen on specified address", slog.String("address", listenAddr))
	err = f.Listen(listenAddr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", listenAddr, err)
	}

	<-shutdownDone

	return nil
}
