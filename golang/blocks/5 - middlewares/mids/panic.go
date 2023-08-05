package mids

import (
	"context"
	"fmt"
	"github.com/fkaanoz/middlewares/app"
	"net/http"
	"runtime/debug"
)

// Panic will handle any panic occurred in handlers. In recovery() part, this panic will be transformed into a known error. (IMPORTANT) h function has named return type.
func Panic() app.Middleware {

	m := func(handler app.Handler) app.Handler {

		// look at return part!
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()

					// err & trace is used for logging purposes.
					err = fmt.Errorf("PANIC [%v] TRACE[%s] : %w", rec, string(trace), app.ErrorServer)
				}
			}()

			return handler(ctx, w, r)
		}
		return h
	}
	return m
}
