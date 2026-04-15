package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	dbpkg "ai-inference-gateway/internal/db"
	"ai-inference-gateway/internal/handlers"
	"ai-inference-gateway/internal/repositories"
	"ai-inference-gateway/internal/services"
)

func main() {
	// Database
	postgresDB, err := dbpkg.InitDB()
	if err != nil {
		log.Fatalf("failed to initialize PostgreSQL connection: %v", err)
	}
	defer postgresDB.Close()

	// Repositories
	userRepo := repositories.NewUserRepository(postgresDB)
	modelRepo := repositories.NewModelRepository(postgresDB)
	taskRepo := repositories.NewTaskRepository(postgresDB)
	txRepo := repositories.NewTransactionRepository(postgresDB)
	workerRepo := repositories.NewWorkerRepository(postgresDB)

	// External clients
	ollama := services.NewOllamaClient("http://localhost:11434")

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
		r.Put("/users/{id}", userH.Update)

		r.Get("/models", modelH.GetAll)

		r.Post("/tasks", taskH.Submit)
		r.Get("/tasks/{id}", taskH.GetByID)
		r.Put("/tasks/{id}", taskH.UpdateTask)
		r.Get("/tasks", taskH.List)
	})

	log.Println("AI Inference Gateway запущено на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("помилка запуску сервера: %v", err)
	}
}
