package main

import (
	"context"
	"errors"
	"github.com/ardanlabs/conf"
	"github.com/fkaanoz/gracefully-shutdown/internal/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var version string
var build = "DEVELOPMENT"

func main() {
	// init conf
	cfg := struct {
		conf.Version
		Web struct {
			ApiHost         string        `conf:"default:0.0.0.0:3000"`
			ReadTimeout     time.Duration `conf:"default:10s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:20s"`
			ShutdownTimeout time.Duration `conf:"default:30s"`
		}
		Debug struct {
			DebugHost string `conf:"default:0.0.0.0:4000"`
		}
	}{
		Version: conf.Version{
			SVN:  build,
			Desc: "gracefully shutdown example",
		},
	}

	prefix := "GRACEFULLY"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			log.Fatal(help)
		}
		log.Fatal("parsing conf error:", err)
	}

	// create channels
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	serverError := make(chan error, 1)

	// start debug server
	go func() {
		if err := http.ListenAndServe(cfg.Debug.DebugHost, app.DebugServer()); err != nil {
			log.Println("debug server cannot serve! : ", err)
		}
	}()

	// init & start app
	api := http.Server{
		Addr:         cfg.Web.ApiHost,
		Handler:      nil,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	go func() {
		serverError <- api.ListenAndServe()
	}()

	// listen channels
	select {
	case err = <-serverError:
		log.Println("server error occured. : ", err)

	case <-shutdown:
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		log.Println("gracefully shutdown started.")

		if err := api.Shutdown(ctx); err != nil {
			log.Println("gracefully shutdown is not possible... CLOSING!")
			api.Close()
			return
		}
		log.Println("gracefully shutdown completed.")
	}
}
