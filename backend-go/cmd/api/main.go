package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KritSuth/fullstack-coding-test/backend-go/internal/handler"
	"github.com/KritSuth/fullstack-coding-test/backend-go/internal/middleware"
	"github.com/KritSuth/fullstack-coding-test/backend-go/internal/repository"
	"github.com/KritSuth/fullstack-coding-test/backend-go/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("userapi")
	repo := repository.NewMongoUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	// Background goroutine: log user count every 10 seconds
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			count, err := svc.Count(context.Background())
			if err != nil {
				log.Printf("[background] count error: %v", err)
				continue
			}
			log.Printf("[background] total users in DB: %d", count)
		}
	}()

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)

	// Protected routes — wrap with Auth middleware
	protected := http.NewServeMux()
	protected.HandleFunc("GET /users", h.ListUsers)
	protected.HandleFunc("GET /users/{id}", h.GetUser)
	protected.HandleFunc("PUT /users/{id}", h.UpdateUser)
	protected.HandleFunc("DELETE /users/{id}", h.DeleteUser)
	mux.Handle("/users", middleware.Auth(protected))
	mux.Handle("/users/", middleware.Auth(protected))

	// Apply logging middleware to everything
	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware.Logger(mux),
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("Server exited")
}
