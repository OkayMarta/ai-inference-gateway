package repositories

import (
	"fmt"
	"sync"

	"ai-inference-gateway/internal/models"
)

// ModelRepository зберігає ІШ-моделі в оперативній пам'яті (in-memory)
// Використовує м'ютекс для безпечного доступу з різних потоків (горутин)
type ModelRepository struct {
	mu     sync.RWMutex // RWMutex дозволяє багатьом читати одночасно, але записувати - тільки по одному
	models map[string]*models.AIModel
}

// NewModelRepository - конструктор, який ініціалізує пусту мапу для моделей
func NewModelRepository() *ModelRepository {
	return &ModelRepository{models: make(map[string]*models.AIModel)}
}

// GetByID шукає модель за її унікальним ідентифікатором
func (r *ModelRepository) GetByID(id string) (*models.AIModel, error) {
	r.mu.RLock() // Блокуємо тільки для читання (Read Lock)
	defer r.mu.RUnlock() // Гарантовано розблокуємо при виході з функції

	m, ok := r.models[id]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", id)
	}
	
	// Важливий момент: ми робимо копію об'єкта (*m) і повертаємо вказівник на копію (&cp)
	// Це робиться для того, щоб хтось ззовні випадково не змінив дані прямо в нашій мапі
	cp := *m
	return &cp, nil
}

// GetAll повертає список усіх доступних моделей
func (r *ModelRepository) GetAll() []*models.AIModel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	out := make([]*models.AIModel, 0, len(r.models))
	for _, m := range r.models {
		cp := *m
		out = append(out, &cp)
	}
	return out
}

// Create додає нову модель у сховище.
func (r *ModelRepository) Create(m *models.AIModel) {
	r.mu.Lock() // Повне блокування (Write Lock), бо ми змінюємо дані
	defer r.mu.Unlock()
	
	r.models[m.ID] = m
}