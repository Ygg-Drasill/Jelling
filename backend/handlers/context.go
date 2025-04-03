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
	return &Context{
		Db:     db,
		Logger: slog.Default(),
	}
}
