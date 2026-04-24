package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"workorder-api/config"
	"workorder-api/handler"
	"workorder-api/queue"
	"workorder-api/repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize RabbitMQ
	rabbitConn, err := config.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// Initialize publisher
	publisher, err := queue.NewPublisher(rabbitConn)
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}
	defer publisher.Close()

	// Initialize repositories
	workorderRepo := repository.NewWorkorderRepository(db)

	// Initialize consumer (background worker)
	consumer, err := queue.NewConsumer(rabbitConn, workorderRepo)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Start consumer
	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// Initialize dispatcher handler
	dispatcher := handler.NewDispatcherHandler(publisher)

	r := mux.NewRouter()

	// Main endpoint - dispatcher akan route berdasarkan function name
	r.HandleFunc("/api/execute", dispatcher.HandleRequest).Methods("POST")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Println("Available functions:")
	log.Println("  - ff_updateWorkorder")
	log.Println("Mode: Async processing with RabbitMQ")

	// Graceful shutdown
	go func() {
		if err := http.ListenAndServe(":"+port, r); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
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
