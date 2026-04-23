package main

import (
	"log"
	"net/http"
	"os"
	"workorder-api/config"
	"workorder-api/handler"
	"workorder-api/middleware"
	"workorder-api/repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	repo := repository.NewWorkorderRepository(db)
	workorderHandler := handler.NewWorkorderHandler(repo)

	r := mux.NewRouter()
	r.Use(middleware.Logger)

	r.HandleFunc("/api/workorder/update", workorderHandler.UpdateWorkorder).Methods("POST")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
