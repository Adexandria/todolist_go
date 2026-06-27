package main

import (
	"SM/handlers"
	"SM/models"
	"SM/repositories/Utilities"
	"SM/services"
	"SM/services/Middlewares"
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AttachRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(Middlewares.RolesKey, roles)
		c.Next()
	}
}

// Main function to set up the application, including dependency injection and route handling
func main() {
	container := services.CreateContainer()

	err := container.Provide(handlers.TaskHandler)
	if err != nil {
		panic(err)
	}

	passwordRule := Utilities.AccessRule{}.PasswordRule.SetMaximumPasswordLength(7).
		SetMinimumPasswordLength(3).EnableSmallLetter().EnableCapitalLetter().EnableNumber().EnableSpecialCharacter()

	accessRule := Utilities.AccessRule{}.EnableEmailConfirmation().EnablePasswordValidation(passwordRule)

	err = container.Provide(accessRule)
	if err != nil {
		panic(err)
	}

	err = container.Provide(handlers.UserHandlerCon)
	if err != nil {
		panic(err)
	}

	err = container.Provide(Utilities.ValidatorCon)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.StaticFile("/docs/openapi.yml", "./docs/openapi.yml")

	url := ginSwagger.URL("/docs/openapi.yml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	err = container.Invoke(func(h *handlers.Handler, jsonHandler *slog.JSONHandler, tokenManager *Utilities.ITokenManager) {
		logger := slog.New(jsonHandler)

		logger.Info("Application is running")

		publicRoutes := router.Group("/api/tasks")

		publicRoutes.Use(Middlewares.AuthenticateMiddleware(*tokenManager),
			Middlewares.AuthorizeMiddleware())

		publicRoutes.GET("/:id", AttachRoles(models.Roles[0], models.Roles[1]), h.GetTaskById)
		publicRoutes.POST("/", AttachRoles(models.Roles[0], models.Roles[1]), h.CreateTask)
		publicRoutes.PUT("/:id", AttachRoles(models.Roles[0], models.Roles[1]), h.UpdateTask)
		publicRoutes.DELETE("/:id", AttachRoles(models.Roles[0], models.Roles[1]), h.DeleteTask)
		publicRoutes.GET("/", AttachRoles(models.Roles[0], models.Roles[1]), h.GetAllTasks)
		publicRoutes.GET("/filter", AttachRoles(models.Roles[0], models.Roles[1]), h.FilterTasks)
		publicRoutes.GET("/search", AttachRoles(models.Roles[0], models.Roles[1]), h.SearchByTask)
	})

	if err != nil {
		panic(err)
	}

	err = router.Run(":8080")
	if err != nil {
		return
	}
}
