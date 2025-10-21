package routes

import (
	"github.com/ADEMOLA200/string-analyzer/cmd/internal/controllers"
	"github.com/ADEMOLA200/string-analyzer/cmd/internal/repository"
	"github.com/ADEMOLA200/string-analyzer/cmd/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Initialize dependencies
	repo := repository.NewInMemoryRepository()
	service := services.NewAnalyzerService(repo)
	controller := controllers.NewStringController(service)

	// Routes
	api := router.Group("/api/v1")
	{
		strings := api.Group("/strings")
		{
			strings.POST("", controller.CreateAnalyzeString)
			strings.GET("", controller.GetAllStrings)
			strings.GET("/filter-by-natural-language", controller.FilterByNaturalLanguage)
			strings.GET("/:string_value", controller.GetString)
			strings.DELETE("/:string_value", controller.DeleteString)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Debug endpoint to see all strings in repository
	router.GET("/debug/strings", func(c *gin.Context) {
		allStrings, err := repo.FindAll()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"total_strings": len(allStrings),
			"strings":       allStrings,
		})
	})

	return router
}
