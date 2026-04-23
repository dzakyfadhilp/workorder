package main

import (
	"log"
	"net/http"
	"os"
	"workorder-api/config"
	"workorder-api/handler"
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

	// Initialize repositories
	workorderRepo := repository.NewWorkorderRepository(db)

	// Initialize dispatcher handler
	dispatcher := handler.NewDispatcherHandler(workorderRepo)

	r := mux.NewRouter()

	// Main endpoint - dispatcher akan route berdasarkan function name
	r.HandleFunc("/api/execute", dispatcher.HandleRequest).Methods("POST")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Println("Available functions:")
	log.Println("  - ff_updateWorkorder")

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
