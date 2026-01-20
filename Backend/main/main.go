package main

import (
	"log"
	"net/http"
	"os"
	"ticketapp/internal/db"
	"ticketapp/internal/handlers"
	"ticketapp/internal/middlewares"
	"ticketapp/internal/repositories"
	"ticketapp/internal/router"
	"ticketapp/internal/services"
	"github.com/joho/godotenv"
)
func main() {
	// -------------------------
	// LOAD ENV
	// -------------------------
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// -------------------------
	// DATABASE
	// -------------------------
	database, err := db.Connect()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// -------------------------
	// REPOSITORIES
	// -------------------------
	userRepo := repositories.NewPostgresUserRepo(database)
	tokenRepo := repositories.NewPostgresRefreshTokenRepo(database)
	// auditRepo := repositories.NewAuditRepo(database)

	// -------------------------
	// SERVICES
	// -------------------------
	jwtService := services.NewJWTService(os.Getenv("JWT_SECRET"))
	otpService := services.NewOTPService()

	// -------------------------
	// HANDLERS
	// -------------------------
	authHandler := handlers.NewAuthHandler(
		userRepo,
		tokenRepo,
		jwtService,
		otpService,
		
	)

	emailSvc := services.NewEmailService()

adminHandler := handlers.NewAdminHandler(
	userRepo,
	emailSvc,
)

	// -------------------------
	// ROUTER
	// -------------------------
	appRouter := router.NewRouter(
		authHandler,
		adminHandler,
		jwtService,
	)

	// -------------------------
	// GLOBAL MIDDLEWARES
	// -------------------------
	handler := middlewares.CORS(appRouter)

	// -------------------------
	// SERVER
	// -------------------------
	log.Println("🚀 Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
