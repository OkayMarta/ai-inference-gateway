package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"task-service/internal/cache"
	"task-service/internal/clients"
	"task-service/internal/config"
	dbpkg "task-service/internal/db"
	"task-service/internal/handlers"
	"task-service/internal/repositories"
	"task-service/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	postgresDB, err := dbpkg.InitDB(cfg.DB)
	if err != nil {
		log.Fatalf("failed to initialize PostgreSQL connection: %v", err)
	}
	defer postgresDB.Close()

	modelRepo := repositories.NewModelRepository(postgresDB)
	taskRepo := repositories.NewTaskRepository(postgresDB)
	workerRepo := repositories.NewWorkerRepository(postgresDB)

	ollama := services.NewOllamaClient(cfg.OllamaURL)
	billingClient := clients.NewBillingClient(cfg.BillingServiceURL, cfg.InternalServiceToken)
	modelCache := cache.NewRedisCache(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
		time.Duration(cfg.Redis.CacheTTLSeconds)*time.Second,
	)

	modelSvc := services.NewModelService(modelRepo, ollama, modelCache)
	workerSvc := services.NewWorkerService(workerRepo, taskRepo, modelRepo, ollama)
	inferenceSvc := services.NewInferenceService(postgresDB, modelRepo, taskRepo, billingClient)

	log.Println("synchronizing models from Ollama")
	if err := modelSvc.SyncFromOllama(); err != nil {
		log.Printf("failed to synchronize models from Ollama: %v", err)
	}

	log.Println("ensuring default worker exists")
	if err := workerSvc.EnsureDefaultWorker(); err != nil {
		log.Fatalf("failed to ensure default worker: %v", err)
	}

	log.Println("refreshing worker supported models")
	if err := workerSvc.RefreshSupportedModels(); err != nil {
		log.Printf("failed to refresh worker/model mappings: %v", err)
	}

	log.Println("starting background worker service")
	workerSvc.Start()

	modelH := handlers.NewModelHandler(modelSvc)
	taskH := handlers.NewTaskHandler(inferenceSvc)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(handlers.RecoveryMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-User-ID", "X-User-Role"},
		AllowCredentials: true,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/models", modelH.GetAll)

		r.Group(func(r chi.Router) {
			r.Use(handlers.InternalServiceTokenMiddleware(cfg.InternalServiceToken))
			r.Post("/tasks", taskH.Submit)
			r.Get("/tasks", taskH.List)
			r.Get("/tasks/{id}", taskH.GetByID)
			r.Put("/tasks/{id}", taskH.UpdateTask)
			r.Delete("/tasks/{id}", taskH.DeleteTask)
		})
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("task-service listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("failed to start task-service: %v", err)
	}
}
