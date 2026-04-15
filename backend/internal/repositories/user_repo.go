package repositories

import (
	"database/sql"
	"errors"
	"fmt"

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
		SELECT id, username, token_balance
		FROM users
		ORDER BY username, id
	`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.TokenBalance); err != nil {
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
	user := &models.User{}

	err := r.db.QueryRow(`
		SELECT id, username, token_balance
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Username, &user.TokenBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("get user by id %s: %w", id, err)
	}

	return user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (id, username, token_balance)
		VALUES ($1, $2, $3)
	`, user.ID, user.Username, user.TokenBalance)
	if err != nil {
		return fmt.Errorf("create user %s: %w", user.ID, err)
	}

	return nil
}

func (r *UserRepository) Update(user *models.User) error {
	result, err := r.db.Exec(`
		UPDATE users
		SET username = $2,
		    token_balance = $3
		WHERE id = $1
	`, user.ID, user.Username, user.TokenBalance)
	if err != nil {
		return fmt.Errorf("update user %s: %w", user.ID, err)
	}

	if err := ensureRowsAffected(result, fmt.Sprintf("user not found: %s", user.ID)); err != nil {
		return err
	}

	return nil
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
	result, err := r.db.Exec(`
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
