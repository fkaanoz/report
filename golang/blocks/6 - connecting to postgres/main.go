package main

import (
	"errors"
	"github.com/ardanlabs/conf"
	"github.com/fkaanoz/connect/internal/app"
	"github.com/fkaanoz/connect/internal/database"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"time"
)

var service = "connect"

func main() {

	logger, err := initLogger(service)
	if err != nil {
		log.Fatal("init logger err:", err)
	}

	if err := run(logger); err != nil {
		logger.Errorw("RUN", "status", "service stopped with error", "ERROR", err)
	}
}

func run(logger *zap.SugaredLogger) error {

	// init config
	cfg := struct {
		conf.Version
		Web struct {
			ApiHost         string        `conf:"default:0.0.0.0:3000"`
			ReadTimeout     time.Duration `conf:"default:10s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:20s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
		}
		DB struct {
			Host     string `conf:"default:localhost"`
			Port     string `conf:"default:5432"`
			Database string `conf:"default:connect"`
			Username string `conf:"default:postgres"`
			Password string `conf:"default:fkaanoz"`
			SSLMode  string `conf:"default:disable"`
		}
	}{
		Version: conf.Version{
			SVN:  service,
			Desc: "postgres connection project",
		},
	}

	prefix := "CONNECT_API"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			log.Fatal(help)
		}
		return err
	}

	// database init
	db, err := database.Open(database.DBConfig{
		Host:        cfg.DB.Host,
		Port:        cfg.DB.Port,
		Username:    cfg.DB.Username,
		Password:    cfg.DB.Password,
		Database:    cfg.DB.Database,
		SSLMode:     cfg.DB.SSLMode,
		MaxIdleConn: 25,
		MaxOpenConn: 25,
	})
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:         cfg.Web.ApiHost,
		Handler:      app.NewApp(logger, db),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	log.Fatal(server.ListenAndServe())

	return nil
}

func initLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": service,
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil

}
