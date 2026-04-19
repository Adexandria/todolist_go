package services

import (
	"SM/models"

	"github.com/glebarez/sqlite"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func RegisterDb() *dig.Container {
	container := dig.New()

	err := container.Provide(connect)
	if err != nil {
		panic(err)
	}

	return container
}
