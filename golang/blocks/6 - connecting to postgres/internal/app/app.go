package app

import (
	"database/sql"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/fkaanoz/connect/internal/controllers"
	"go.uber.org/zap"
	"net/http"
)

type App struct {
	*httptreemux.ContextMux
	logger *zap.SugaredLogger
	db     *sql.DB
}

func NewApp(logger *zap.SugaredLogger, db *sql.DB) *App {
	a := &App{
		ContextMux: httptreemux.NewContextMux(),
		logger:     logger,
		db:         db,
	}
	v1(a)
	return a
}

func v1(app *App) {
	app.ContextMux.Handle(http.MethodGet, "/", controllers.Test)
}
