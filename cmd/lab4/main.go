package main

import (
	"context"
	"log"
	app "lab4/internal/app/pkg"
	"os"
)

// @title Lab4 API
// @version 1.0
// @description API with JWT Auth

// @contact.name API Support

// @license.name AS IS

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

func main() {
	log.Println("Initializing server")
	application, err := app.New(context.Background())
	if err != nil {
		log.Println("cant create application")
		os.Exit(2)
	}

	log.Println("Application start!")
	if err := application.Run(); err != nil {
		log.Println("application error:", err)
		os.Exit(1)
	}
	log.Println("Application terminated!")
}
