package repositories

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"ai-inference-gateway/internal/models"
)

// TaskRepository stores prompt tasks submitted for processing.
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

func (r *TaskRepository) GetAll() []*models.PromptTask {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*models.PromptTask, 0, len(r.tasks))
	for _, t := range r.tasks {
		cp := *t
		out = append(out, &cp)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})

	return out
}

// GetByUserID returns all tasks for a specific user.
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

	// Sort from newest to oldest by creation time.
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

// Create adds a new task and stamps its creation time.
func (r *TaskRepository) Create(task *models.PromptTask) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.CreatedAt = time.Now()
	r.tasks[task.ID] = task
}

// Complete is called by a worker when processing finishes successfully.
// It updates the task status to Completed and stores the result.
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

func (r *TaskRepository) Fail(id, result string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tasks[id]
	if !ok {
		return fmt.Errorf("task not found: %s", id)
	}
	t.Status = models.StatusFailed
	t.Result = result
	return nil
}

// GetNextQueued atomically finds the oldest queued task supported by the worker.
func (r *TaskRepository) GetNextQueued(supportedModels []string) *models.PromptTask {
	r.mu.Lock() // Use Lock because the status is updated before returning.
	defer r.mu.Unlock()

	// Convert supported model IDs to a set for fast lookups.
	supported := make(map[string]bool, len(supportedModels))
	for _, m := range supportedModels {
		supported[m] = true
	}

	var oldest *models.PromptTask

	// Scan all tasks.
	for _, t := range r.tasks {
		// Skip tasks that are not queued or not supported by this worker.
		if t.Status != models.StatusQueued || !supported[t.ModelID] {
			continue
		}
		// Track the oldest matching task.
		if oldest == nil || t.CreatedAt.Before(oldest.CreatedAt) {
			oldest = t
		}
	}

	// Move the selected task to Processing before returning it.
	if oldest != nil {
		oldest.Status = models.StatusProcessing
		cp := *oldest
		return &cp // Return a copy for worker processing.
	}
	return nil
}
