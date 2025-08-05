package main

import (
	"github.com/PavelBradnitski/WbTechL2/internal/handlers"
	"github.com/PavelBradnitski/WbTechL2/internal/repositories"
	"github.com/PavelBradnitski/WbTechL2/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	eventRepo := repositories.NewEventRepo()
	eventService := services.NewEventService(eventRepo)
	eventHandler := handlers.NewEventHandler(eventService)
	router := gin.Default()
	eventHandler.RegisterRoutes(router)
	router.Run(":8080")
}
