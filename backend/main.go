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
	// Repositories
	userRepo := repositories.NewUserRepository()
	modelRepo := repositories.NewModelRepository()
	taskRepo := repositories.NewTaskRepository()
	txRepo := repositories.NewTransactionRepository()
	workerRepo := repositories.NewWorkerRepository()
	ollama := services.NewOllamaClient("http://localhost:11434")

	// Temporary in-memory bootstrap.
	// Users and models are expected to come from SQL migrations in Lab 3.
	seedWorkers(workerRepo)

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
		r.Get("/tasks", taskH.List)
	})

	// Server start
	log.Println("AI Inference Gateway Р·Р°РїСѓС‰РµРЅРѕ РЅР° http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("РџРѕРјРёР»РєР° Р·Р°РїСѓСЃРєСѓ СЃРµСЂРІРµСЂР°: %v", err)
	}
}

func seedWorkers(workerRepo *repositories.WorkerRepository) {
	// Тимчасово лишаємо in-memory воркерів, але перелік model ID більше не формується під час старту застосунку. Ці ID мають відповідати seed-даним БД.
	supportedModelIDs := []string{"model-1", "model-2", "model-3", "model-4"}

	workerRepo.Create(&models.WorkerNode{ID: "worker-1", SupportedModels: supportedModelIDs, Status: models.WorkerIdle})
	workerRepo.Create(&models.WorkerNode{ID: "worker-2", SupportedModels: supportedModelIDs, Status: models.WorkerIdle})
	workerRepo.Create(&models.WorkerNode{ID: "worker-3", SupportedModels: supportedModelIDs, Status: models.WorkerIdle})
	log.Println("Seed: Р”РѕРґР°РЅРѕ 3 С„РѕРЅРѕРІС– РІРѕСЂРєРµСЂРё")
}
