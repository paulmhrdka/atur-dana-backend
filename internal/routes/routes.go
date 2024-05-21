package routes

import (
	"atur-dana/internal/handlers"
	"atur-dana/internal/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	// Auth Routes
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", handlers.Signup).Methods("POST")
	auth.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected Routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.HandleFunc("/transactions", handlers.GetTransactions).Methods("GET")
	protected.HandleFunc("/transactions", handlers.CreateTransaction).Methods("POST")

	protected.HandleFunc("/budgets", handlers.GetBudgets).Methods("GET")
	protected.HandleFunc("/budgets", handlers.CreateBudget).Methods("POST")
}
