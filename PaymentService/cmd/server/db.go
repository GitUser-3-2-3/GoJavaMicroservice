package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(cfg.Timeout)
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
