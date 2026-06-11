package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bcastillo-2022474/relay/internal/config"
	appcommand "github.com/bcastillo-2022474/relay/internal/domain/application/command"
	endpointcommand "github.com/bcastillo-2022474/relay/internal/domain/endpoint/command"
	etcommand "github.com/bcastillo-2022474/relay/internal/domain/event_type/command"
	msgcommand "github.com/bcastillo-2022474/relay/internal/domain/message/command"
	"github.com/bcastillo-2022474/relay/internal/fakes"
	relayhttp "github.com/bcastillo-2022474/relay/internal/http"
	"github.com/bcastillo-2022474/relay/internal/postgres"
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

	cfg, err := config.Load()
	if err != nil {
		appLog.Error("loading config", "err", err)
		os.Exit(1)
	}

	authz := fakes.AllowAllAuthorization{}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		appLog.Error("creating db pool", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	// Fail at startup, not on request N: wiring errors must not hide.
	if err := pool.Ping(context.Background()); err != nil {
		appLog.Error("db unreachable", "err", err)
		os.Exit(1)
	}

	appRepo := postgres.NewApplicationRepo(pool)
	eventTypeRepo := postgres.NewEventTypeRepo(pool)
	endpointRepo := postgres.NewEndpointRepo(pool)
	msgRepo := postgres.NewMessageRepo(pool)

	createApp := appcommand.NewCreateCommand(appRepo, authz, appLog)
	createEventType := etcommand.NewCreateTypeCommand(eventTypeRepo, appRepo, authz, appLog)
	createEndpoint := endpointcommand.NewCreateCommand(endpointRepo, appRepo, authz, appLog)
	publishMsg := msgcommand.NewPublishCommand(eventTypeRepo, appRepo, msgRepo, authz, appLog)

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

	httpLog.Info("api listening", "addr", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		httpLog.Error("server error", "err", err)
		os.Exit(1)
	}
}

type chiSlogLogger struct{ log *slog.Logger }

func (l *chiSlogLogger) Print(v ...any)                 { l.log.Info(fmt.Sprint(v...)) }
func (l *chiSlogLogger) Printf(format string, v ...any) { l.log.Info(fmt.Sprintf(format, v...)) }
