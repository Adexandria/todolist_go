package services

import (
	"SM/repositories"
	"SM/repositories/Utilities"
	"log/slog"
	"os"

	"go.uber.org/dig"
)

// CreateContainer sets up the dependency injection container with all necessary services and repositories
func CreateContainer() *dig.Container {
	container := RegisterDb()

	err := container.Provide(generateHandler)
	if err != nil {
		panic(err)
	}

	err = container.Provide(repositories.NewTaskRepo)
	if err != nil {
		panic(err)
	}

	err = container.Provide(NewTaskService)
	if err != nil {
		panic(err)
	}

	err = container.Provide(repositories.NewUserRepo)
	if err != nil {
		panic(err)
	}

	err = container.Provide(Utilities.NewTokenManager)
	if err != nil {
		panic(err)
	}

	err = container.Provide(Utilities.NewPasswordManager)
	if err != nil {
		panic(err)
	}

	err = container.Provide(NewUserService)
	if err != nil {
		panic(err)
	}

	err = container.Provide(NewAuthenticationService)
	if err != nil {
		panic(err)
	}

	err = container.Provide(NewEmailService)
	if err != nil {
		panic(err)
	}

	err = container.Provide(NewSmsService)
	if err != nil {
		panic(err)
	}

	return container
}

func generateHandler() *slog.JSONHandler {
	return slog.NewJSONHandler(os.Stderr, nil)
}
