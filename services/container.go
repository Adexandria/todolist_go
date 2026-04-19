package services

import (
	"SM/repositories"

	"go.uber.org/dig"
)

func CreateContainer() *dig.Container {
	container := RegisterDb()

	err := container.Provide(repositories.TaskRepo)
	if err != nil {
		panic(err)
	}

	err = container.Provide(TaskServiceCon)
	if err != nil {
		panic(err)
	}

	return container
}
