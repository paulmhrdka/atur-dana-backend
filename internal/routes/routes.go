package routes

import (
	"atur-dana/internal/handlers"
	"atur-dana/internal/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(r *mux.Router) {
	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Auth Routes
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", handlers.Register).Methods("POST")
	auth.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected Routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.HandleFunc("/transactions", handlers.GetTransactions).Methods("GET")
	protected.HandleFunc("/transactions", handlers.CreateTransaction).Methods("POST")

	protected.HandleFunc("/budgets", handlers.GetBudgets).Methods("GET")
	protected.HandleFunc("/budgets", handlers.CreateBudget).Methods("POST")
}
