package main

import (
	"chanterelle/internal/handlers"
	"chanterelle/internal/repositories"
	"chanterelle/internal/services"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chanterelle/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate *validator.Validate

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	validate = validator.New()

	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Get database handle
	db := client.Database(cfg.MongoDatabase)

	// Initialize repositories with MongoDB
	contactRepo := repositories.NewMongoContactRepository(db)
	verificationRepo := repositories.NewMongoVerificationRepository(db)
	contactService := services.NewContactService(contactRepo)
	notificationService := services.NewNotificationService(cfg)
	verificationService := services.NewVerificationService(cfg, verificationRepo)

	// Initialize handlers
	handlers := handlers.NewHandlers(contactService, notificationService, verificationService, cfg)

	// Set up router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// Public routes
	r := router.Group("/api")
	// Contact creation (public)
	r.POST("/contacts", handlers.CreateContact)
	// Authentication endpoints
	r.POST("/send-verification", handlers.SendVerification)
	r.POST("/verify-code", handlers.VerifyCode)

	// Protected routes
	authGroup := r.Group("")
	authGroup.Use(handlers.JWTAuth())

	// Get all contacts
	authGroup.GET("/contacts", handlers.GetContacts)
	authGroup.DELETE("/contacts/:id", handlers.DeleteContact)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", cfg.Port)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
