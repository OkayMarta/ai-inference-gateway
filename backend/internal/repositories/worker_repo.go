package repositories

import (
	"fmt"
	"sync"

	"ai-inference-gateway/internal/models"
)

// WorkerRepository зберігає вузли-обробники (воркери) в оперативній пам'яті.
// Це імітація бази даних для ЛР №2
type WorkerRepository struct {
	// sync.RWMutex потрібен для потокобезпеки. 
	// Оскільки наш HTTP-сервер і фоновий цикл воркера працюють паралельно, 
	// вони можуть спробувати одночасно читати і писати в мапу, що призведе до падіння програми
	mu      sync.RWMutex
	workers map[string]*models.WorkerNode // Ключ - це ID воркера, значення - вказівник на об'єкт
}

// NewWorkerRepository - це конструктор (фабрична функція), який ініціалізує порожню мапу
func NewWorkerRepository() *WorkerRepository {
	return &WorkerRepository{workers: make(map[string]*models.WorkerNode)}
}

// GetAll повертає список усіх воркерів
func (r *WorkerRepository) GetAll() []*models.WorkerNode {
	// RLock (Read Lock) дозволяє багатьом горутинам читати дані одночасно, 
	// але блокує запис, поки читання не завершиться
	r.mu.RLock()
	defer r.mu.RUnlock() // defer гарантує, що блокування зніметься при виході з функції

	out := make([]*models.WorkerNode, 0, len(r.workers))
	for _, w := range r.workers {
		// ВАЖЛИВО: Ми створюємо копію об'єкта (cp := *w).
		// Якби ми повернули оригінальний вказівник (w), хтось інший міг би 
		// випадково змінити статус воркера в обхід м'ютекса і зламати логіку.
		cp := *w
		out = append(out, &cp)
	}
	return out
}

// Create додає нового воркера в "базу".
func (r *WorkerRepository) Create(w *models.WorkerNode) {
	// Lock (Write Lock) повністю блокує мапу і для читання, і для запису іншими потоками.
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.workers[w.ID] = w
}

// UpdateStatus змінює статус воркера (наприклад, з Idle на Busy).
func (r *WorkerRepository) UpdateStatus(id string, status models.WorkerStatus) error {
	r.mu.Lock() // Блокуємо для запису
	defer r.mu.Unlock()
	
	w, ok := r.workers[id]
	if !ok {
		return fmt.Errorf("worker not found: %s", id) // Повертаємо помилку, якщо такого воркера немає
	}
	w.Status = status // Оновлюємо статус оригінального об'єкта в мапі
	return nil
}

// GetIdle повертає лише тих воркерів, які зараз вільні (мають статус WorkerIdle).
// Цей метод постійно викликається фоновим процесом (WorkerService), щоб знайти вільні "руки" для задач.
func (r *WorkerRepository) GetIdle() []*models.WorkerNode {
	r.mu.RLock() // Блокуємо тільки для читання
	defer r.mu.RUnlock()
	
	var out []*models.WorkerNode
	for _, w := range r.workers {
		if w.Status == models.WorkerIdle {
			cp := *w // Знову робимо копію для безпеки
			out = append(out, &cp)
		}
	}
	return out
}