package main

import (
	"atur-dana/internal/common"
	"atur-dana/internal/db"
	"atur-dana/internal/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

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
