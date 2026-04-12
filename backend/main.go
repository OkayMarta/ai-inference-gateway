package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ai-inference-gateway/internal/handlers"
	"ai-inference-gateway/internal/models"
	"ai-inference-gateway/internal/repositories"
	"ai-inference-gateway/internal/services"
)

func main() {
	// Repositories
	userRepo := repositories.NewUserRepository()
	modelRepo := repositories.NewModelRepository()
	taskRepo := repositories.NewTaskRepository()
	txRepo := repositories.NewTransactionRepository()
	workerRepo := repositories.NewWorkerRepository()
	ollama := services.NewOllamaClient("http://localhost:11434")

	// Seed
	seedUsers(userRepo)
	loadModels(ollama, modelRepo)
	seedWorkers(workerRepo, modelRepo)

	// Services
	userSvc := services.NewUserService(userRepo)
	modelSvc := services.NewModelService(modelRepo)
	inferenceSvc := services.NewInferenceService(userRepo, modelRepo, taskRepo, txRepo)
	workerSvc := services.NewWorkerService(workerRepo, taskRepo, modelRepo, ollama)
	workerSvc.Start()

	// Handlers
	userH := handlers.NewUserHandler(userSvc)
	modelH := handlers.NewModelHandler(modelSvc)
	taskH := handlers.NewTaskHandler(inferenceSvc)

	// Router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(handlers.RecoveryMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/users", userH.GetAll)
		r.Get("/users/{id}", userH.GetByID)

		r.Get("/models", modelH.GetAll)

		r.Post("/tasks", taskH.Submit)
		r.Get("/tasks/{id}", taskH.GetByID)
		r.Get("/tasks", taskH.GetByUserID)
	})

	// Server start
	log.Println("AI Inference Gateway запущено на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Помилка запуску сервера: %v", err)
	}
}

// Seed helpers

func seedUsers(repo *repositories.UserRepository) {
	repo.Create(&models.User{ID: "user-1", Username: "alice", TokenBalance: 100})
	repo.Create(&models.User{ID: "user-2", Username: "bob", TokenBalance: 5})
	repo.Create(&models.User{ID: "user-3", Username: "charlie", TokenBalance: 200})
	log.Println("Seed: Додано 3 користувачі")
}

func loadModels(ollama *services.OllamaClient, repo *repositories.ModelRepository) {
	ollamaModels, err := ollama.ListModels()
	if err != nil {
		log.Printf("Ollama недоступна: %v", err)
		log.Println("Завантажуємо базові (симульовані) моделі...")
		seedDefaultModels(repo)
		return
	}

	if len(ollamaModels) == 0 {
		log.Println("В Ollama немає завантажених моделей. Використовуємо базові...")
		seedDefaultModels(repo)
		return
	}

	for _, m := range ollamaModels {
		id := sanitizeID(m.Name)
		cost := costBySize(m.Size)
		repo.Create(&models.AIModel{
			ID:          id,
			Name:        m.Name,
			Description: fmt.Sprintf("Ollama модель · %s", formatSize(m.Size)),
			TokenCost:   cost,
		})
	}

	log.Printf("Завантажено %d моделей з Ollama", len(ollamaModels))
}

func seedDefaultModels(repo *repositories.ModelRepository) {
	repo.Create(&models.AIModel{
		ID:          "model-1",
		Name:        "Llama-3",
		Description: "Велика мовна модель для генерації тексту",
		TokenCost:   5,
	})
	repo.Create(&models.AIModel{
		ID:          "model-2",
		Name:        "Stable-Diffusion",
		Description: "Модель для генерації зображень з тексту",
		TokenCost:   10,
	})
	repo.Create(&models.AIModel{
		ID:          "model-3",
		Name:        "Whisper",
		Description: "Модель розпізнавання мовлення (Speech-to-text)",
		TokenCost:   3,
	})
	repo.Create(&models.AIModel{
		ID:          "model-4",
		Name:        "GPT-4o",
		Description: "Просунута мультимодальна ШІ-модель",
		TokenCost:   15,
	})
	log.Println("Seed: Додано 4 базові (симульовані) моделі")
}

func seedWorkers(workerRepo *repositories.WorkerRepository, modelRepo *repositories.ModelRepository) {
	allModels := modelRepo.GetAll()

	ids := make([]string, len(allModels))
	for i, m := range allModels {
		ids[i] = m.ID
	}

	workerRepo.Create(&models.WorkerNode{ID: "worker-1", SupportedModels: ids, Status: models.WorkerIdle})
	workerRepo.Create(&models.WorkerNode{ID: "worker-2", SupportedModels: ids, Status: models.WorkerIdle})
	workerRepo.Create(&models.WorkerNode{ID: "worker-3", SupportedModels: ids, Status: models.WorkerIdle})
	log.Println("Seed: Додано 3 фонові воркери")
}

// Utility helpers

func sanitizeID(name string) string {
	s := strings.ReplaceAll(name, ":", "-")
	s = strings.ReplaceAll(s, "/", "-")
	return s
}

func costBySize(bytes int64) float64 {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	switch {
	case gb < 2:
		return 3
	case gb < 5:
		return 5
	case gb < 15:
		return 10
	default:
		return 15
	}
}

func formatSize(bytes int64) string {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	if gb >= 1 {
		return fmt.Sprintf("%.1f GB", gb)
	}

	mb := float64(bytes) / (1024 * 1024)
	return fmt.Sprintf("%.0f MB", mb)
}
