package routes

import (
	"github.com/ADEMOLA200/string-analyzer/cmd/internal/controllers"
	"github.com/ADEMOLA200/string-analyzer/cmd/internal/repository"
	"github.com/ADEMOLA200/string-analyzer/cmd/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	repo := repository.NewInMemoryRepository()
	service := services.NewAnalyzerService(repo)
	controller := controllers.NewStringController(service)

	router.POST("/strings", controller.CreateAnalyzeString)
	router.GET("/strings", controller.GetAllStrings)
	router.GET("/strings/filter-by-natural-language", controller.FilterByNaturalLanguage)
	router.GET("/strings/:string_value", controller.GetString)
	router.DELETE("/strings/:string_value", controller.DeleteString)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
