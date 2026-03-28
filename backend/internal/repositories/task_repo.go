package repositories

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"ai-inference-gateway/internal/models"
)

// TaskRepository зберігає задачі (промпти), які користувачі відправляють на обробку
type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*models.PromptTask
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{tasks: make(map[string]*models.PromptTask)}
}

func (r *TaskRepository) GetByID(id string) (*models.PromptTask, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	t, ok := r.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	cp := *t
	return &cp, nil
}

// GetByUserID повертає всі задачі конкретного користувача
func (r *TaskRepository) GetByUserID(userID string) []*models.PromptTask {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var out []*models.PromptTask
	for _, t := range r.tasks {
		if t.UserID == userID {
			cp := *t
			out = append(out, &cp)
		}
	}
	
	// Сортуємо задачі від найновіших до найстаріших (по даті створення)
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

// Create додає нову задачу і автоматично проставляє їй час створення
func (r *TaskRepository) Create(task *models.PromptTask) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	task.CreatedAt = time.Now()
	r.tasks[task.ID] = task
}

// Complete викликається воркером, коли Ollama згенерувала відповідь
// Змінює статус задачі на Completed і записує результат
func (r *TaskRepository) Complete(id, result string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	t, ok := r.tasks[id]
	if !ok {
		return fmt.Errorf("task not found: %s", id)
	}
	t.Status = models.StatusCompleted
	t.Result = result
	return nil
}

// GetNextQueued - найхитріша функція. Вона атомарно (за одне блокування) шукає найстарішу задачу в статусі Queued для моделей, які підтримує воркер
func (r *TaskRepository) GetNextQueued(supportedModels []string) *models.PromptTask {
	r.mu.Lock() // Використовуємо Lock, бо ми одразу змінимо статус знайденої задачі
	defer r.mu.Unlock()

	// Перетворюємо масив підтримуваних моделей на мапу для швидкого пошуку
	supported := make(map[string]bool, len(supportedModels))
	for _, m := range supportedModels {
		supported[m] = true
	}

	var oldest *models.PromptTask
	
	// Перебираємо всі задачі
	for _, t := range r.tasks {
		// Якщо задача не в черзі АБО модель не підтримується цим воркером - пропускаємо
		if t.Status != models.StatusQueued || !supported[t.ModelID] {
			continue
		}
		// Шукаємо задачу з найменшим часом (найстарішу)
		if oldest == nil || t.CreatedAt.Before(oldest.CreatedAt) {
			oldest = t
		}
	}
	
	// Якщо знайшли задачу, одразу переводимо її в статус "В обробці" (Processing)
	if oldest != nil {
		oldest.Status = models.StatusProcessing
		cp := *oldest
		return &cp // Повертаємо воркеру копію задачі для виконання
	}
	return nil
}