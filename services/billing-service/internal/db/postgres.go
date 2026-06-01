package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"billing-service/internal/config"

	_ "github.com/lib/pq"
)

const (
	postgresConnectTimeout = 30 * time.Second
	postgresRetryDelay     = 2 * time.Second
)

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

func InitDB(cfg config.DBConfig) (*sql.DB, error) {
	deadline := time.Now().Add(postgresConnectTimeout)
	attempt := 1

	for {
		db, err := sql.Open("postgres", cfg.ConnectionString())
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
			db.Close()
			err = fmt.Errorf("ping postgres connection: %w", err)
		} else {
			err = fmt.Errorf("open postgres connection: %w", err)
		}

		if time.Now().Add(postgresRetryDelay).After(deadline) {
			return nil, fmt.Errorf("initialize postgres connection after %s: %w", postgresConnectTimeout, err)
		}

		log.Printf("PostgreSQL connection attempt %d failed: %v; retrying in %s", attempt, err, postgresRetryDelay)
		time.Sleep(postgresRetryDelay)
		attempt++
	}
}
