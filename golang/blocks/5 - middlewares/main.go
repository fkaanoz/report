package main

import (
	"fmt"
	"github.com/dimfeld/httptreemux/v5"
	app2 "github.com/fkaanoz/middlewares/internal/app"
	usergroup2 "github.com/fkaanoz/middlewares/internal/handlers/usergroup"
	mids2 "github.com/fkaanoz/middlewares/internal/mids"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
)

var service = "middleware"

func main() {

	zapLogger, err := initLogger(service)
	if err != nil {
		// if logger could not be initialized, continuing to serve is a little bit meaningless. log.Fatal calls os.Exit(1)
		log.Fatal("log creation error : ", err)
	}

	a := app2.App{
		ContextMux:  httptreemux.NewContextMux(),
		Middlewares: []app2.Middleware{mids2.Logger(zapLogger), mids2.Panic(), mids2.Error(zapLogger)}, // app wised middlewares.
	}

	a.Handle(http.MethodGet, "/users", usergroup2.List)
	a.Handle(http.MethodPost, "/users/create", usergroup2.Create)

	log.Fatal(http.ListenAndServe(":4000", a))
}

// initLogger provides SugaredLogger. This logger is used for structured logging purposes. Key-Value type logging mechanism is important.
func initLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.InitialFields = map[string]interface{}{
		"service": service,
	}
	config.OutputPaths = []string{"stdout"}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("creating logger err : %w", err)
	}

	return logger.Sugar(), nil
}
