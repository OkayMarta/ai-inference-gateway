package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	appdb "ai-inference-gateway/internal/db"
	"ai-inference-gateway/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	rows, err := r.db.Query(`
		SELECT id, username, email, password_hash, token_balance, role, created_at
		FROM users
		ORDER BY username, id
	`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	return r.getByID(r.db, id)
}

func (r *UserRepository) GetByIDTx(tx appdb.DBTX, id string) (*models.User, error) {
	return r.getByID(tx, id)
}

func (r *UserRepository) getByID(exec appdb.DBTX, id string) (*models.User, error) {
	user := &models.User{}

	err := exec.QueryRow(`
		SELECT id, username, email, password_hash, token_balance, role, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.TokenBalance,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("get user by id %s: %w", id, err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := r.db.QueryRow(`
		SELECT id, username, email, password_hash, token_balance, role, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.TokenBalance,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %s", email)
		}
		return nil, fmt.Errorf("get user by email %s: %w", email, err)
	}

	return user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	if user.Role == "" {
		user.Role = models.RoleUser
	}

	if user.CreatedAt.IsZero() {
		err := r.db.QueryRow(`
			INSERT INTO users (id, username, email, password_hash, token_balance, role)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING created_at
		`, user.ID, user.Username, user.Email, user.PasswordHash, user.TokenBalance, user.Role).Scan(&user.CreatedAt)
		if err != nil {
			return fmt.Errorf("create user %s: %w", user.ID, err)
		}
		user.CreatedAt = user.CreatedAt.UTC()
		return nil
	}

	user.CreatedAt = user.CreatedAt.UTC()

	_, err := r.db.Exec(`
		INSERT INTO users (id, username, email, password_hash, token_balance, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, user.ID, user.Username, user.Email, user.PasswordHash, user.TokenBalance, user.Role, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("create user %s: %w", user.ID, err)
	}

	return nil
}

func (r *UserRepository) Update(user *models.User) error {
	result, err := r.db.Exec(`
		UPDATE users
		SET username = $2,
		    email = $3,
		    password_hash = $4,
		    token_balance = $5,
		    role = $6
		WHERE id = $1
	`, user.ID, user.Username, user.Email, user.PasswordHash, user.TokenBalance, user.Role)
	if err != nil {
		return fmt.Errorf("update user %s: %w", user.ID, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("user not found: %s", user.ID)); err != nil {
		return err
	}

	return nil
}

func scanUser(scanner interface {
	Scan(dest ...any) error
}) (*models.User, error) {
	user := &models.User{}
	if err := scanner.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.TokenBalance,
		&user.Role,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	user.CreatedAt = user.CreatedAt.UTC()
	return user, nil
}

func (r *UserRepository) Delete(id string) error {
	result, err := r.db.Exec(`
		DELETE FROM users
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete user %s: %w", id, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("user not found: %s", id)); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateBalance(id string, balance float64) error {
	return r.updateBalance(r.db, id, balance)
}

func (r *UserRepository) UpdateBalanceTx(tx appdb.DBTX, id string, balance float64) error {
	return r.updateBalance(tx, id, balance)
}

func (r *UserRepository) DeductBalanceTx(tx appdb.DBTX, id string, amount float64) error {
	result, err := tx.Exec(`
		UPDATE users
		SET token_balance = token_balance - $2
		WHERE id = $1
		  AND token_balance >= $2
	`, id, amount)
	if err != nil {
		return fmt.Errorf("deduct balance for user %s: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows for user %s balance deduction: %w", id, err)
	}

	if rowsAffected == 0 {
		return errors.New("insufficient token balance")
	}

	return nil
}

func (r *UserRepository) updateBalance(exec appdb.DBTX, id string, balance float64) error {
	result, err := exec.Exec(`
		UPDATE users
		SET token_balance = $2
		WHERE id = $1
	`, id, balance)
	if err != nil {
		return fmt.Errorf("update balance for user %s: %w", id, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("user not found: %s", id)); err != nil {
		return err
	}

	return nil
}
