package handlers

import (
	"database/sql"
	"log/slog"
)

type Context struct {
	Db     *sql.DB
	Logger *slog.Logger
}

func NewContext(db *sql.DB) *Context {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	return &Context{
		Db:     db,
		Logger: slog.Default(),
	}
}
