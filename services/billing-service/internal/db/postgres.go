package db

import (
	"database/sql"
	"fmt"

	"billing-service/internal/config"

	_ "github.com/lib/pq"
)

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

func InitDB(cfg config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres connection: %w", err)
	}

	return db, nil
}
