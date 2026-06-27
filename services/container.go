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

	err = container.Provide(repositories.TaskRepo)
	if err != nil {
		panic(err)
	}

	err = container.Provide(TaskServiceCon)
	if err != nil {
		panic(err)
	}

	err = container.Provide(repositories.UserRepoCon)
	if err != nil {
		panic(err)
	}
	err = container.Provide(Utilities.TokenManagerCons)
	if err != nil {
		panic(err)
	}

	err = container.Provide(Utilities.PasswordManagerCons)
	if err != nil {
		panic(err)
	}

	err = container.Provide(UserServiceCon)
	if err != nil {
		panic(err)
	}

	err = container.Provide(AuthenticationServiceCons)
	if err != nil {
		panic(err)
	}

	err = container.Provide(AuthenticationServiceCons)
	if err != nil {
		panic(err)
	}

	return container
}

func generateHandler() *slog.JSONHandler {
	return slog.NewJSONHandler(os.Stderr, nil)
}
