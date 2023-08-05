package app

import (
	"database/sql"
	"github.com/dimfeld/httptreemux/v5"
	"go.uber.org/zap"
)

type App struct {
	*httptreemux.ContextMux
	logger *zap.SugaredLogger
	db     *sql.DB
}

func NewApp(logger *zap.SugaredLogger, db *sql.DB) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		logger:     logger,
		db:         db,
	}
}
