package main

import (
	"os"

	"github.com/PavelBradnitski/WbTechL2/internal/handlers"
	"github.com/PavelBradnitski/WbTechL2/internal/repositories"
	"github.com/PavelBradnitski/WbTechL2/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	eventRepo := repositories.NewEventRepo()
	eventService := services.NewEventService(eventRepo)
	eventHandler := handlers.NewEventHandler(eventService)
	router := gin.Default()
	eventHandler.RegisterRoutes(router)

	router.Run(":" + os.Getenv("API_PORT"))
}
