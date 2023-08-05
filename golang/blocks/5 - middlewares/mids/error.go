package mids

import (
	"context"
	"errors"
	"github.com/fkaanoz/middlewares/app"
	"github.com/fkaanoz/middlewares/handlers"
	"go.uber.org/zap"
	"net/http"
)

// Error will return a middleware which deal with errors returned from handlers.
func Error(logger *zap.SugaredLogger) app.Middleware {

	m := func(handler app.Handler) app.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			if err := handler(ctx, w, r); err != nil {
				logger.Infow("REQUEST", "status", "UNSUCCESSFUL", "error", err.Error(), "traceID", app.GetTraceID(ctx))

				// Default case any other errors returned from handlers.
				switch {
				case errors.Is(err, app.ErrorRequest):
					handlers.Respond(w, http.StatusBadRequest, app.ErrorRequest.Error())
				case errors.Is(err, app.ErrorServer):
					handlers.Respond(w, http.StatusInternalServerError, app.ErrorServer.Error())
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
