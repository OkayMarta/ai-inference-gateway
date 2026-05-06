package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"gateway-service/internal/clients"
	"gateway-service/internal/config"
	"gateway-service/internal/handlers"
	"gateway-service/internal/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()
	billingClient := clients.NewBillingClient(cfg.BillingServiceURL)
	taskClient := clients.NewTaskClient(cfg.TaskServiceURL)
	proxy := handlers.NewProxyHandler(billingClient, taskClient)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(handlers.RecoveryMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", proxy.Register)
		r.Post("/auth/login", proxy.Login)
		r.Get("/models", proxy.Models)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWT(cfg.JWTSecret))

			r.Get("/auth/me", proxy.Me)
			r.Post("/tasks", proxy.Tasks)
			r.Get("/tasks", proxy.Tasks)
			r.Get("/tasks/{id}", proxy.TaskByID)
			r.Put("/tasks/{id}", proxy.TaskByID)
			r.Delete("/tasks/{id}", proxy.TaskByID)
		})
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("gateway-service listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("failed to start gateway-service: %v", err)
	}
}
