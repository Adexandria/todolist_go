package main

import (
	"SM/handlers"
	"SM/services"
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Main function to set up the application, including dependency injection and route handling
func main() {
	container := services.CreateContainer()

	err := container.Provide(handlers.TaskHandler)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.StaticFile("/docs/openapi.yml", "./docs/openapi.yml")

	url := ginSwagger.URL("/docs/openapi.yml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	err = container.Invoke(func(h *handlers.Handler, jsonHandler *slog.JSONHandler) {
		logger := slog.New(jsonHandler)

		logger.Info("Application is running")

		publicRoutes := router.Group("/api/tasks")

		publicRoutes.GET("/:id", h.GetTaskById)
		publicRoutes.POST("/", h.CreateTask)
		publicRoutes.PUT("/:id", h.UpdateTask)
		publicRoutes.DELETE("/:id", h.DeleteTask)
		publicRoutes.GET("/", h.GetAllTasks)
		publicRoutes.GET("/filter", h.FilterTasks)
		publicRoutes.GET("/search", h.SearchByTask)
	})

	if err != nil {
		panic(err)
	}

	err = router.Run(":8080")
	if err != nil {
		return
	}
}
