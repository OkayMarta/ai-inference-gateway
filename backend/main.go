package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ai-inference-gateway/internal/handlers"
	"ai-inference-gateway/internal/repositories"
	"ai-inference-gateway/internal/services"
)

func main() {
	// 1. --- Repositories (Сховища даних) ---
	userRepo := repositories.NewUserRepository()
	modelRepo := repositories.NewModelRepository()
	taskRepo := repositories.NewTaskRepository()
	txRepo := repositories.NewTransactionRepository()
	// workerRepo := repositories.NewWorkerRepository() // Розкоментув. пізніше

	// 2. --- Seed data (Тестові дані) ---
	// TODO: Додам пізніше

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
	r.Use(chimw.Logger)          // Логує кожен запит у консоль
	r.Use(chimw.Recoverer)       // Додатковий захист від падінь
	r.Use(handlers.RecoveryMiddleware) // Наш власний middleware для JSON-помилок
	
	// Налаштування CORS (щоб фронтенд на іншому порту міг робити запити)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	}))

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