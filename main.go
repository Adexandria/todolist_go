package main

import (
	"SM/handlers"
	"SM/services"

	"github.com/gin-gonic/gin"
)

func main() {
	container := services.CreateContainer()

	err := container.Provide(handlers.TaskHandler)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	err = container.Invoke(func(h *handlers.Handler) {
		publicRoutes := router.Group("/tasks")

		publicRoutes.GET("/:id", h.GetTaskById)
		publicRoutes.POST("/create", h.CreateTask)
		publicRoutes.PUT("/:id", h.UpdateTask)
		publicRoutes.DELETE("/:id", h.DeleteTask)
		publicRoutes.GET("/", h.GetAllTasks)
		publicRoutes.GET("/year", h.GetTaskByYear)
		publicRoutes.GET("/month", h.GetTaskByMonthAndYear)
		publicRoutes.GET("/date", h.GetTaskByCreatedDate)
	})

	if err != nil {
		panic(err)
	}

	err = router.Run(":8080")
	if err != nil {
		return
	}

}
