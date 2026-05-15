package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"billing-service/internal/config"
	dbpkg "billing-service/internal/db"
	"billing-service/internal/handlers"
	"billing-service/internal/repositories"
	"billing-service/internal/services"
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

	userRepo := repositories.NewUserRepository(postgresDB)
	txRepo := repositories.NewTransactionRepository(postgresDB)
	resetRepo := repositories.NewPasswordResetRepository(postgresDB)
	emailSvc := services.NewEmailServiceFromEnv()

	authSvc := services.NewAuthService(userRepo, resetRepo, emailSvc, cfg.JWTSecret)
	userSvc := services.NewUserService(userRepo)
	billingSvc := services.NewBillingService(postgresDB, userRepo, txRepo)

	authH := handlers.NewAuthHandler(authSvc)
	userH := handlers.NewUserHandler(userSvc)
	billingH := handlers.NewBillingHandler(billingSvc)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(handlers.RecoveryMiddleware)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", authH.Register)
		r.Post("/auth/login", authH.Login)
		r.Post("/auth/forgot-password", authH.ForgotPassword)
		r.Post("/auth/reset-password", authH.ResetPassword)

		r.Group(func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(authSvc))
			r.Get("/users/me", userH.Me)
			r.Get("/users", userH.GetAll)
			r.Put("/users/{id}", userH.Update)
		})

		r.Get("/users/{id}", userH.GetByID)
	})

	r.Route("/internal", func(r chi.Router) {
		r.Use(handlers.InternalServiceTokenMiddleware(cfg.InternalServiceToken))
		r.Get("/users/{id}", userH.GetByID)
		r.Post("/billing/charge", billingH.Charge)
		r.Post("/billing/refund", billingH.Refund)
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("billing-service listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("failed to start billing-service: %v", err)
	}
}
