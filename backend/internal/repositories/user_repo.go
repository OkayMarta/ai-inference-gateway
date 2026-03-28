package repositories

import (
	"fmt"
	"sync"

	"ai-inference-gateway/internal/models"
)

// UserRepository зберігає інформацію про користувачів системи в оперативній пам'яті.
// Забезпечує потокобезпечний доступ до даних (щоб баланс не списався двічі одночасно)
type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

// NewUserRepository - конструктор для ініціалізації мапи користувачів
func NewUserRepository() *UserRepository {
	return &UserRepository{users: make(map[string]*models.User)}
}

// GetByID шукає користувача за ID (наприклад, "user-1")
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	r.mu.RLock() // Блокування тільки для читання
	defer r.mu.RUnlock()
	
	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	
	// Повертаємо копію, щоб сервіси випадково не змінили баланс 
	// прямо в кеші без виклику спеціальної функції UpdateBalance
	cp := *u
	return &cp, nil
}

// GetAll повертає список усіх зареєстрованих користувачів (для адмін-панелі чи фронтенду)
func (r *UserRepository) GetAll() []*models.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	out := make([]*models.User, 0, len(r.users))
	for _, u := range r.users {
		cp := *u
		out = append(out, &cp)
	}
	return out
}

// Create додає нового користувача в систему
func (r *UserRepository) Create(u *models.User) {
	r.mu.Lock() // Повне блокування для запису
	defer r.mu.Unlock()
	
	r.users[u.ID] = u
}

// UpdateBalance оновлює кількість токенів на рахунку користувача
// Це ключова функція для білінгу: вона викликається, коли ми списуємо кошти за промпт
func (r *UserRepository) UpdateBalance(id string, balance float64) error {
	r.mu.Lock() // Обов'язково Lock, адже ми змінюємо фінансові дані
	defer r.mu.Unlock()
	
	u, ok := r.users[id]
	if !ok {
		return fmt.Errorf("user not found: %s", id)
	}
	u.TokenBalance = balance
	return nil
}