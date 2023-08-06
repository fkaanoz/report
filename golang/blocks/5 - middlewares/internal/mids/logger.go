package mids

import (
	"context"
	app2 "github.com/fkaanoz/middlewares/internal/app"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Logger will write some request related data into stdout. It gets TraceID from context.
func Logger(logger *zap.SugaredLogger) app2.Middleware {

	m := func(handler app2.Handler) app2.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			t := time.Now()
			logger.Infow("REQUEST", "status", "START", "traceID", app2.GetTraceID(ctx))

			if err := handler(ctx, w, r); err != nil {
				// when an error occurred, error middleware will log status
				return err
			}

			logger.Infow("REQUEST", "status", "SUCCESS", "traceID", app2.GetTraceID(ctx), "elapsed time", time.Since(t))

			return nil
		}
		return h
	}
	return m
}
