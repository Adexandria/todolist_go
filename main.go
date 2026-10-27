package main

import (
	"SM/handlers"
	"SM/models"
	"SM/repositories"
	"SM/repositories/Utilities"
	"SM/services"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found:", err)
	}

	container := services.CreateContainer()

	passwordRule := Utilities.PasswordRule{}.SetMaximumPasswordLength(7).
		SetMinimumPasswordLength(3).EnableSmallLetter().EnableCapitalLetter().EnableNumber().EnableSpecialCharacter()

	accessRule := Utilities.AccessRule{}.EnableEmailConfirmation().EnablePasswordValidation(passwordRule)

	must(container.Provide(func() *Utilities.AccessRule { return &accessRule }))

	must(container.Provide(handlers.TaskHandler))

	must(container.Provide(handlers.UserHandlerCon))

	must(container.Provide(Utilities.ValidatorCon))

	router := gin.Default()

	router.StaticFile("/docs/openapi.yml", "./docs/openapi.yml")

	url := ginSwagger.URL("/docs/openapi.yml")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	must(container.Invoke(func(h *handlers.Handler, jsonHandler *slog.JSONHandler,
		tokenManager *Utilities.TokenManager, u *repositories.UserRepository,
		a *handlers.UserHandler, passwordManager *Utilities.PasswordManager) {
		logger := slog.New(jsonHandler)

		CreateDefaultAdmin(u, passwordManager, logger)

		logger.Info("Application is running")

		registerRoutes(router, h, tokenManager, a)

	}))

	err := router.Run(":8081")
	if err != nil {
		panic(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateDefaultAdmin(userRepos *repositories.UserRepository, passwordManager *Utilities.PasswordManager, logger *slog.Logger) {

	passwordHash, err := passwordManager.HashPassword(os.Getenv("ADMIN_PASSWORD"))

	if err != nil {
		logger.Error(err.Error())
	}

	user := &models.User{
		Username:           "admin",
		Password:           passwordHash,
		FirstName:          "admin",
		LastName:           "admin",
		Email:              "admin@admin",
		Role:               1,
		EmailConfirmed:     true,
		AuthenticationType: 2,
	}
	isExist, _ := userRepos.VerifyUserByEmail(user.Email)
	if isExist {
		return
	}
	err = userRepos.CreateUser(user)

	if err != nil {
		logger.Error("Error creating default admin user")
	}

	logger.Info("Created Admin Successfully")
}
