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
	_, err := db.Exec(`PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL,
    expiry_date INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    color TEXT
);

CREATE TABLE IF NOT EXISTS article_topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_id INTEGER NOT NULL,
    topic_id INTEGER NOT NULL,
    foreign key (article_id) REFERENCES articles(id),
    foreign key (topic_id) REFERENCES topics(id)
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

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_id INTEGER NOT NULL,
    author_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    resolved INTEGER DEFAULT 0
)

`)

	if err != nil {
		return err
	}

	return nil
}
