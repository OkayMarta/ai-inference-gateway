package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ai-inference-gateway/internal/handlers"
	"ai-inference-gateway/internal/models"
	"ai-inference-gateway/internal/repositories"
	"ai-inference-gateway/internal/services"
)

func main() {
	// 1. --- Repositories (Сховища даних) ---
	userRepo := repositories.NewUserRepository()
	modelRepo := repositories.NewModelRepository()
	taskRepo := repositories.NewTaskRepository()
	txRepo := repositories.NewTransactionRepository()
	workerRepo := repositories.NewWorkerRepository()

	// 2. --- Seed data (Тестові дані) ---
	// Заповнюємо In-Memory сховища тестовими даними при старті сервера
	seedUsers(userRepo)
	seedDefaultModels(modelRepo)
	seedWorkers(workerRepo, modelRepo)

	// 3. --- Services (Бізнес-логіка) ---
	userSvc := services.NewUserService(userRepo)
	modelSvc := services.NewModelService(modelRepo)
	inferenceSvc := services.NewInferenceService(userRepo, modelRepo, taskRepo, txRepo)

	// TODO: Ініціалізація Ollama та Воркера

	// 4. --- Handlers (HTTP Контролери) ---
	userH := handlers.NewUserHandler(userSvc)
	modelH := handlers.NewModelHandler(modelSvc)
	taskH := handlers.NewTaskHandler(inferenceSvc)

	// 5. --- Роутер (Chi) ---
	r := chi.NewRouter()
	
	// Middlewares (проміжні обробники)
	r.Use(chimw.Logger)          
	r.Use(chimw.Recoverer)       
	r.Use(handlers.RecoveryMiddleware) 
	
	// Налаштування CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	}))

	// /healthz (Health-Check) едпоінт показує, що сервер працює
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Реєстрація маршрутів (URL)
	r.Route("/api", func(r chi.Router) {
		r.Get("/users", userH.GetAll)
		r.Get("/users/{id}", userH.GetByID)

		r.Get("/models", modelH.GetAll)

		r.Post("/tasks", taskH.Submit)
		r.Get("/tasks/{id}", taskH.GetByID)
		r.Get("/tasks", taskH.GetByUserID)
	})

	// 6. --- Запуск сервера ---
	log.Println("AI Inference Gateway запущено на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Помилка запуску сервера: %v", err)
	}
}

// --- Функції наповнення бази даних тестовими даними (Seeding) ---

// seedUsers створює трьох користувачів з різним балансом токенів
func seedUsers(repo *repositories.UserRepository) {
	repo.Create(&models.User{ID: "user-1", Username: "alice", TokenBalance: 100})
	repo.Create(&models.User{ID: "user-2", Username: "bob", TokenBalance: 5}) // Бобу ледь вистачає на 1 запит
	repo.Create(&models.User{ID: "user-3", Username: "charlie", TokenBalance: 200})
	log.Println("Seed: Додано 3 користувачі")
}

// seedDefaultModels створює список доступних ШІ-моделей (заглушки) із різною вартістю
func seedDefaultModels(repo *repositories.ModelRepository) {
	repo.Create(&models.AIModel{
		ID: "model-1", Name: "Llama-3",
		Description: "Велика мовна модель для генерації тексту",
		TokenCost:   5, // Вартість виклику цієї моделі - 5 токенів
	})
	repo.Create(&models.AIModel{
		ID: "model-2", Name: "Stable-Diffusion",
		Description: "Модель для генерації зображень з тексту",
		TokenCost:   10,
	})
	repo.Create(&models.AIModel{
		ID: "model-3", Name: "Whisper",
		Description: "Модель розпізнавання мовлення (Speech-to-text)",
		TokenCost:   3,
	})
	repo.Create(&models.AIModel{
		ID: "model-4", Name: "GPT-4o",
		Description: "Просунута мультимодальна ШІ-модель",
		TokenCost:   15,
	})
	log.Println("Seed: Додано 4 базові (симульовані) моделі")
}

// seedWorkers створює три "віртуальні" обчислювальні вузли (Воркери).
// Налаштов. їх так, що кожен з них підтримує ВСІ створені раніше моделі.
// Статус кожного при старті - Idle (Вільний).
func seedWorkers(workerRepo *repositories.WorkerRepository, modelRepo *repositories.ModelRepository) {
	// Спочатку отримуємо всі моделі, щоб дізнатись їхні ID
	allModels := modelRepo.GetAll()
	
	// Збираємо список ID всіх моделей
	ids := make([]string, len(allModels))
	for i, m := range allModels {
		ids[i] = m.ID
	}

	// Створюємо трьох воркерів
	workerRepo.Create(&models.WorkerNode{ID: "worker-1", SupportedModels: ids, Status: models.WorkerIdle})
	workerRepo.Create(&models.WorkerNode{ID: "worker-2", SupportedModels: ids, Status: models.WorkerIdle})
	workerRepo.Create(&models.WorkerNode{ID: "worker-3", SupportedModels: ids, Status: models.WorkerIdle})
	log.Println("Seed: Додано 3 фонові воркери (кожен підтримує всі моделі)")
}