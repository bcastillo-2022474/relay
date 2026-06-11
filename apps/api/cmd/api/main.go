package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	charmlog "github.com/charmbracelet/log"
	appcommand "github.com/bcastillo-2022474/relay/internal/domain/application/command"
	endpointcommand "github.com/bcastillo-2022474/relay/internal/domain/endpoint/command"
	etcommand "github.com/bcastillo-2022474/relay/internal/domain/event_type/command"
	msgcommand "github.com/bcastillo-2022474/relay/internal/domain/message/command"
	"github.com/bcastillo-2022474/relay/internal/fakes"
	relayhttp "github.com/bcastillo-2022474/relay/internal/http"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	charmLogger := charmlog.NewWithOptions(os.Stdout, charmlog.Options{
		Level:           charmlog.DebugLevel,
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	})
	appLog  := slog.New(charmLogger.WithPrefix("app"))
	httpLog := slog.New(charmLogger.WithPrefix("http"))

	authz := fakes.AllowAllAuthorization{}

	appRepo      := fakes.NewInMemoryApplicationRepo()
	eventTypeRepo := fakes.NewInMemoryEventTypeRepo()
	endpointRepo  := fakes.NewInMemoryEndpointRepo()
	msgRepo       := fakes.NewInMemoryMessageRepository()

	createApp      := appcommand.NewCreateCommand(appRepo, authz, appLog)
	createEventType := etcommand.NewCreateTypeCommand(eventTypeRepo, appRepo, authz, appLog)
	createEndpoint  := endpointcommand.NewCreateCommand(endpointRepo, appRepo, authz, appLog)
	publishMsg      := msgcommand.NewPublishCommand(eventTypeRepo, appRepo, msgRepo, authz, appLog)

	router := chi.NewMux()
	router.Use(relayhttp.AuthMiddleware)
	router.Use(chimiddleware.RequestLogger(&chimiddleware.DefaultLogFormatter{
		Logger:  &chiSlogLogger{log: httpLog},
		NoColor: false,
	}))

	api := humachi.New(router, huma.DefaultConfig("Relay API", "0.1.0"))

	relayhttp.RegisterApplicationRoutes(api, createApp)
	relayhttp.RegisterEventTypeRoutes(api, createEventType)
	relayhttp.RegisterEndpointRoutes(api, createEndpoint)
	relayhttp.RegisterMessageRoutes(api, publishMsg)

	httpLog.Info("api listening", "addr", ":8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		httpLog.Error("server error", "err", err)
		os.Exit(1)
	}
}

type chiSlogLogger struct{ log *slog.Logger }

func (l *chiSlogLogger) Print(v ...any)                 { l.log.Info(fmt.Sprint(v...)) }
func (l *chiSlogLogger) Printf(format string, v ...any) { l.log.Info(fmt.Sprintf(format, v...)) }
