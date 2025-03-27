package handlers

import "database/sql"

type Context struct {
	Db *sql.DB
}

func NewContext(db *sql.DB) *Context {
	return &Context{Db: db}
}
