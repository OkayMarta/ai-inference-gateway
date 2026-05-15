package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	dbpkg "ai-inference-gateway/internal/db"
	"ai-inference-gateway/internal/handlers"
	"ai-inference-gateway/internal/repositories"
	"ai-inference-gateway/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	appEnv := appEnv()
	frontendOrigin := envOrDefault("FRONTEND_ORIGIN", "http://localhost:5173")
	jwtSecret := requiredSecret("JWT_SECRET", "dev-secret", appEnv)

	postgresDB, err := dbpkg.InitDB()
	if err != nil {
		log.Fatalf("failed to initialize PostgreSQL connection: %v", err)
	}
	defer postgresDB.Close()

	log.Println("initializing PostgreSQL-backed repositories...")
	userRepo := repositories.NewUserRepository(postgresDB)
	modelRepo := repositories.NewModelRepository(postgresDB)
	taskRepo := repositories.NewTaskRepository(postgresDB)
	txRepo := repositories.NewTransactionRepository(postgresDB)
	workerRepo := repositories.NewWorkerRepository(postgresDB)

	log.Println("initializing Ollama client...")
	ollama := services.NewOllamaClient("http://localhost:11434")

	// Сервіси для startup sync створюємо раніше, бо вони потрібні ще до запуску HTTP-шару.
	modelSvc := services.NewModelService(modelRepo, ollama)
	workerSvc := services.NewWorkerService(workerRepo, taskRepo, modelRepo, ollama)

	log.Println("synchronizing models from Ollama...")
	if err := modelSvc.SyncFromOllama(); err != nil {
		log.Printf("failed to synchronize models from Ollama: %v", err)
	} else {
		syncedModels, err := modelSvc.GetAll()
		if err != nil {
			log.Printf("models synchronized from Ollama, but failed to load synced catalog: %v", err)
		} else {
			log.Printf("successfully synchronized %d models from Ollama", len(syncedModels))
		}

		log.Println("refreshing worker/model mappings...")
		if err := workerSvc.RefreshSupportedModels(); err != nil {
			log.Printf("failed to refresh worker/model mappings: %v", err)
		} else {
			log.Println("successfully refreshed worker/model mappings")
		}
	}

	log.Println("initializing application services...")
	authSvc := services.NewAuthService(userRepo, jwtSecret)
	userSvc := services.NewUserService(userRepo)
	inferenceSvc := services.NewInferenceService(postgresDB, userRepo, modelRepo, taskRepo, txRepo)

	log.Println("starting background worker service...")
	workerSvc.Start()

	log.Println("initializing HTTP handlers...")
	authH := handlers.NewAuthHandler(authSvc)
	userH := handlers.NewUserHandler(userSvc)
	modelH := handlers.NewModelHandler(modelSvc)
	taskH := handlers.NewTaskHandler(inferenceSvc)

	log.Println("configuring HTTP router...")
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(handlers.RecoveryMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{frontendOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", authH.Register)
		r.Post("/auth/login", authH.Login)
		r.Get("/models", modelH.GetAll)

		r.Group(func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(authSvc))

			r.Get("/auth/me", authH.Me)

			r.Get("/users", userH.GetAll)
			r.Get("/users/{id}", userH.GetByID)
			r.Put("/users/{id}", userH.Update)

			r.Post("/tasks", taskH.Submit)
			r.Get("/tasks", taskH.List)
			r.Get("/tasks/{id}", taskH.GetByID)
			r.Put("/tasks/{id}", taskH.UpdateTask)
			r.Delete("/tasks/{id}", taskH.DeleteTask)
		})
	})

	log.Println("AI Inference Gateway startup complete")
	log.Println("AI Inference Gateway запущено на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("помилка запуску сервера: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func appEnv() string {
	return strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
}

func isDevelopment(appEnv string) bool {
	return appEnv == "" || appEnv == "development"
}

func requiredSecret(key, fallback, appEnv string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	if isDevelopment(appEnv) {
		return fallback
	}

	log.Fatalf("%s is required when APP_ENV=%s", key, appEnv)
	return ""
}
