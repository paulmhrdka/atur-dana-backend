package main

import (
	_ "atur-dana/docs"
	"atur-dana/internal/common"
	"atur-dana/internal/db"
	"atur-dana/internal/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// @title           Atur Dana API
// @version         1.0
// @description     Personal finance management API for tracking transactions and budgets.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Type "Bearer" followed by a space and the JWT token.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.Init()
	common.InitValidator()

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	http.Handle("/", r)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
