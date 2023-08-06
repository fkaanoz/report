package mids

import (
	"context"
	"errors"
	app2 "github.com/fkaanoz/middlewares/internal/app"
	"github.com/fkaanoz/middlewares/internal/handlers"
	"go.uber.org/zap"
	"net/http"
)

// Error will return a middleware which deal with errors returned from handlers.
func Error(logger *zap.SugaredLogger) app2.Middleware {

	m := func(handler app2.Handler) app2.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			if err := handler(ctx, w, r); err != nil {
				logger.Infow("REQUEST", "status", "UNSUCCESSFUL", "error", err.Error(), "traceID", app2.GetTraceID(ctx))

				// Default case any other errors returned from handlers.
				switch {
				case errors.Is(err, app2.ErrorRequest):
					handlers.Respond(w, http.StatusBadRequest, app2.ErrorRequest.Error())
				case errors.Is(err, app2.ErrorServer):
					handlers.Respond(w, http.StatusInternalServerError, app2.ErrorServer.Error())
				default:
					handlers.Respond(w, http.StatusInternalServerError, "unknown error")
				}
			}
			return nil
		}
		return h
	}
	return m
}
