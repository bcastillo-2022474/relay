package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bcastillo-2022474/relay/internal/config"
	appcommand "github.com/bcastillo-2022474/relay/internal/domain/application/command"
	endpointcommand "github.com/bcastillo-2022474/relay/internal/domain/endpoint/command"
	etcommand "github.com/bcastillo-2022474/relay/internal/domain/event_type/command"
	msgcommand "github.com/bcastillo-2022474/relay/internal/domain/message/command"
	"github.com/bcastillo-2022474/relay/internal/fakes"
	relayhttp "github.com/bcastillo-2022474/relay/internal/http"
	"github.com/bcastillo-2022474/relay/internal/postgres"
	"github.com/bcastillo-2022474/relay/internal/relayer"
	charmlog "github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	charmLogger := charmlog.NewWithOptions(os.Stdout, charmlog.Options{
		Level:           charmlog.DebugLevel,
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	})
	appLog := slog.New(charmLogger.WithPrefix("app"))
	httpLog := slog.New(charmLogger.WithPrefix("http"))
	relayLog := slog.New(charmLogger.WithPrefix("relayer"))

	cfg, err := config.Load()
	if err != nil {
		appLog.Error("loading config", "err", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		appLog.Error("creating db pool", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		appLog.Error("db unreachable", "err", err)
		os.Exit(1)
	}

	authz := fakes.AllowAllAuthorization{}
	broker := fakes.NewLogBroker(relayLog)

	appRepo := postgres.NewApplicationRepo(pool)
	eventTypeRepo := postgres.NewEventTypeRepo(pool)
	endpointRepo := postgres.NewEndpointRepo(pool)
	msgRepo := postgres.NewMessageRepo(pool)

	createApp := appcommand.NewCreateCommand(appRepo, authz, appLog)
	createEventType := etcommand.NewCreateTypeCommand(eventTypeRepo, appRepo, authz, appLog)
	createEndpoint := endpointcommand.NewCreateCommand(endpointRepo, appRepo, authz, appLog)
	publishMsg := msgcommand.NewPublishCommand(eventTypeRepo, appRepo, msgRepo, authz, appLog)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	relayerDone := make(chan struct{})
	go func() {
		defer close(relayerDone)
		relayer.New(msgRepo, broker, relayLog).Run(ctx)
	}()

	router := chi.NewMux()
	router.Use(relayhttp.AuthMiddleware)
	router.Use(chimiddleware.RequestLogger(&chimiddleware.DefaultLogFormatter{
		Logger:  &chiSlogLogger{log: httpLog},
		NoColor: false,
	}))

	api := relayhttp.NewAPI(router, "Relay API", "0.1.0")

	relayhttp.RegisterApplicationRoutes(api, createApp)
	relayhttp.RegisterEventTypeRoutes(api, createEventType)
	relayhttp.RegisterEndpointRoutes(api, createEndpoint)
	relayhttp.RegisterMessageRoutes(api, publishMsg)

	// ListenAndServe ignores ctx, and NotifyContext swallows further SIGINTs —
	// without an explicit Shutdown the process would hang after Ctrl+C with
	// the relayer stopped but the server still serving.
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: router}
	go func() {
		<-ctx.Done()
		httpLog.Info("shutting down http server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			httpLog.Error("http shutdown", "err", err)
		}
	}()

	httpLog.Info("api listening", "addr", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		httpLog.Error("server error", "err", err)
		os.Exit(1)
	}

	// Let the relayer finish its in-flight batch before the process exits.
	<-relayerDone
	appLog.Info("shutdown complete")
}

type chiSlogLogger struct{ log *slog.Logger }

func (l *chiSlogLogger) Print(v ...any)                 { l.log.Info(fmt.Sprint(v...)) }
func (l *chiSlogLogger) Printf(format string, v ...any) { l.log.Info(fmt.Sprintf(format, v...)) }
