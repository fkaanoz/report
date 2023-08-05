package app

import (
	"context"
	"fmt"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
	"net/http"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is wrapper around ContextMux. We want to create context for each request, and this is possible with App struct.
type App struct {
	*httptreemux.ContextMux
	Middlewares []Middleware // app wise middlewares.
}

func (a *App) Handle(method string, path string, handler Handler) {

	handler = Wrap(handler, a.Middlewares...)

	h := func(w http.ResponseWriter, r *http.Request) {

		values := Values{
			traceID: uuid.New().String(),
		}

		ctx := context.WithValue(context.Background(), Key, values)

		if err := handler(ctx, w, r); err != nil {
			fmt.Printf("app wise error!")
		}
	}

	// call the underlying handle function.
	a.ContextMux.Handle(method, path, h)
}
