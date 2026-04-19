package models

import (
	"time"

	"gorm.io/gorm"
)

type CreateTaskDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}
type UpdateTaskDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type TaskDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	CreatedDate string `json:"created_date"`
}

type TaskPaginationDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	CreatedDate string `json:"created_date"`
}
type PaginatedDTO struct {
	Tasks    []TaskPaginationDTO
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type Task struct {
	gorm.Model
	ID          uint `Gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string
	CreatedAt   time.Time `Gorm:"autoCreateTime"`
	DueDate     *time.Time
	UpdatedAt   time.Time `Gorm:"autoUpdateTime"`
}
