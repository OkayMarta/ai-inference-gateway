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
	// 1. --- Repositories (Сховища даних) ---
	// Ініціалізуємо in-memory сховища. Всередині вони використовують потокобезпечні мапи (з sync.RWMutex), щоб уникнути конфліктів при паралельних запитах
	userRepo := repositories.NewUserRepository()
	modelRepo := repositories.NewModelRepository()
	taskRepo := repositories.NewTaskRepository()
	txRepo := repositories.NewTransactionRepository()
	workerRepo := repositories.NewWorkerRepository()
	ollama := services.NewOllamaClient("http://localhost:11434") // Ініціалізуємо HTTP-клієнт для зв'язку з локальною нейромережею Ollama

	// 2. --- Seed data (Тестові дані) ---
	// Наповнюємо систему початковими даними при старті сервера для зручності тестування
	seedUsers(userRepo)
	loadModels(ollama, modelRepo) // Пробуємо знайти реальні моделі або вантажимо заглушки
	seedWorkers(workerRepo, modelRepo)

	// 3. --- Services (Бізнес-логіка) ---
	userSvc := services.NewUserService(userRepo)
	modelSvc := services.NewModelService(modelRepo)
	inferenceSvc := services.NewInferenceService(userRepo, modelRepo, taskRepo, txRepo) // Гол. сервіс, який відповідає за створ. задач та білінг (списання токенів)
	workerSvc := services.NewWorkerService(workerRepo, taskRepo, modelRepo, ollama) // Сервіс фонової обробки задач. Він відповідає за те, щоб брати задачі з черги (Queued) і відправляти їх в Ollama
	
	workerSvc.Start() // Запускаємо фоновий цикл в окремій горутині (асинхронно)

	// 4. --- Handlers (HTTP Контролери) ---
	// Контролери відповідають виключно за прийом HTTP-запитів, валідацію JSON та виклик відповідних сервісів. Вони не містять бізнес-логіки
	userH := handlers.NewUserHandler(userSvc)
	modelH := handlers.NewModelHandler(modelSvc)
	taskH := handlers.NewTaskHandler(inferenceSvc)

	// 5. --- Роутер (Chi) ---
	r := chi.NewRouter()
	
	// Middlewares (проміжні обробники)
	r.Use(chimw.Logger)                // Логування кожного запиту у консоль
	r.Use(chimw.Recoverer)             // Захист від критичних помилок (panic)
	r.Use(handlers.RecoveryMiddleware) // Власний обробник для форматування помилок у єдиний JSON
	
	// Налаштування CORS, щоб фронтенд (наприклад, React на порту 5173) міг робити API-запити
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	}))

	// /healthz (Health-Check) ендпоінт для перевірки життєздатності сервера
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Реєстрація публічних API-маршрутів
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

// seedUsers створює 3-ох користувачів з різним балансом токенів для симуляції як успішних сценаріїв, так і браку коштів
func seedUsers(repo *repositories.UserRepository) {
	repo.Create(&models.User{ID: "user-1", Username: "alice", TokenBalance: 100})
	repo.Create(&models.User{ID: "user-2", Username: "bob", TokenBalance: 5})
	repo.Create(&models.User{ID: "user-3", Username: "charlie", TokenBalance: 200})
	log.Println("Seed: Додано 3 користувачі")
}

// loadModels перевіряє доступність локал. Ollama. Якщо вона працює — завантаж. список реальних моделей. Якщо ні — використов. fallback до базових симульованих моделей
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
		cost := costBySize(m.Size) // Вартість динамічно розраховується від розміру
		repo.Create(&models.AIModel{
			ID:          id,
			Name:        m.Name,
			Description: fmt.Sprintf("Ollama модель · %s", formatSize(m.Size)),
			TokenCost:   cost,
		})
	}
	log.Printf("Завантажено %d моделей з Ollama", len(ollamaModels))
}

// seedDefaultModels створює список доступних ШІ-моделей (заглушки) із різною вартістю. Викликається лише тоді, коли реальна нейромережа недоступна
func seedDefaultModels(repo *repositories.ModelRepository) {
	repo.Create(&models.AIModel{
		ID: "model-1", Name: "Llama-3",
		Description: "Велика мовна модель для генерації тексту",
		TokenCost:   5,
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

// seedWorkers створює три обчислювальні вузли (Воркери).
// Вони призначені для асинхронної обробки задач. Кожен з них підтримує ВСІ створені раніше моделі, стартовий статус - Idle (Вільний)
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

// --- Допоміжні утиліти ---

// sanitizeID перетворює імена моделей Ollama (наприклад "llama3:latest") у безпечний для URL формат ID ("llama3-latest")
func sanitizeID(name string) string {
	s := strings.ReplaceAll(name, ":", "-")
	s = strings.ReplaceAll(s, "/", "-")
	return s
}

// costBySize автоматично призначає вартість запиту до моделі на основі її розміру. Чим "важча" модель (більше ГБ), тим більше токенів коштує її виклик
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

// formatSize перетворює байти у зрозумілий рядок розміру для відображення на UI
func formatSize(bytes int64) string {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	if gb >= 1 {
		return fmt.Sprintf("%.1f GB", gb)
	}
	mb := float64(bytes) / (1024 * 1024)
	return fmt.Sprintf("%.0f MB", mb)
}