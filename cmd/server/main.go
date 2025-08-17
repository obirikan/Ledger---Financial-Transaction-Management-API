package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/obirikan/Ledger---Financial-Transaction-Management-API/internal/config"
	"github.com/obirikan/Ledger---Financial-Transaction-Management-API/internal/handlers"
	"github.com/obirikan/Ledger---Financial-Transaction-Management-API/internal/middleware"
	"github.com/obirikan/Ledger---Financial-Transaction-Management-API/internal/repository"
	"github.com/obirikan/Ledger---Financial-Transaction-Management-API/internal/service"
)

func main() {

	db := config.SetupDB()

	repo := repository.NewPostgresRepo(db)
	authSvc := service.NewAuthService(repo)
	ledgerSvc := service.NewLedgerService(repo)

	authHandler := handlers.NewAuthHandler(authSvc)
	ledgerHandler := handlers.NewLedgerHandler(ledgerSvc)

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.JWTAuthMiddleware())
	{
		authorized.GET("/accounts", ledgerHandler.ListAccounts)
		authorized.POST("/transfer", ledgerHandler.TransferFunds)
	}

	log.Println("Starting server at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
