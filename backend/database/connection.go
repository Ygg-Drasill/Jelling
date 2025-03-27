package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func Open() *sql.DB {
	db, err := sql.Open("sqlite", "../data/jelling.sqlite")
	if err != nil {
		panic(err)
	}

	return db
}

func Setup(db *sql.DB) error {
	_, err := db.Exec(`
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	hash TEXT NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL,
    expiry_date INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author_id INTEGER,
    last_updated INTEGER,
    title TEXT NOT NULL,
    summary TEXT,
    source TEXT,
    content_type INTEGER,
    foreign key (author_id) REFERENCES users(id) ON DELETE CASCADE
);

		`)
	if err != nil {
		return err
	}

	return nil
}
