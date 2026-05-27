package db

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

// DBTX описує мінімальний контракт для виконання SQL-запитів як через *sql.DB,
// так і через *sql.Tx. Це дозволяє сервісному шару керувати транзакцією, а
// репозиторіям лишатись тонким шаром доступу до даних.
type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// InitDB ініціалізує підключення до PostgreSQL на основі змінних середовища.
func InitDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := url.URL{
		Scheme:  "postgres",
		User:    url.UserPassword(user, password),
		Host:    net.JoinHostPort(host, port),
		Path:    "/" + name,
		RawPath: "/" + url.PathEscape(name),
	}
	query := dsn.Query()
	query.Set("sslmode", sslMode)
	query.Set("TimeZone", "UTC")
	dsn.RawQuery = query.Encode()
	connStr := dsn.String()

	// Використовуємо database/sql як стандартний абстрактний шар для роботи з БД: він дає уніфікований API, пул з'єднань і дозволяє підключати драйвер pq без жорсткої прив'язки бізнес-логіки до конкретної реалізації драйвера.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("не вдалося відкрити конфігурацію підключення до PostgreSQL: %v", err)
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	// sql.Open лише готує об'єкт підключення і не гарантує, що БД вже доступна. Викликаємо Ping одразу, щоб перевірити реальну мережеву доступність, облікові дані та коректність конфігурації на етапі ініціалізації.
	if err := db.Ping(); err != nil {
		log.Printf("не вдалося підключитися до PostgreSQL: %v", err)
		db.Close()
		return nil, fmt.Errorf("ping postgres connection: %w", err)
	}

	log.Println("PostgreSQL connected successfully")

	return db, nil
}
