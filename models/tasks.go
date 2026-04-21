package models

import (
	"time"

	"gorm.io/gorm"
)

// Data Transfer Object (DTO) for creating a new task
type CreateTaskDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

// Data Transfer Object (DTO) for updating an existing task
type UpdateTaskDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

// Data Transfer Object (DTO) for representing a task in API responses
type TaskDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	CreatedDate string `json:"created_date"`
}

// Data Transfer Object (DTO) for representing a task in paginated API responses
type TaskPaginationDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	CreatedDate string `json:"created_date"`
}

// Data Transfer Object (DTO) for representing paginated task responses
type PaginatedDTO struct {
	Tasks    []TaskPaginationDTO
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// Task model representing the database structure for tasks
type Task struct {
	gorm.Model
	ID          uint `Gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string
	CreatedAt   time.Time `Gorm:"autoCreateTime"`
	DueDate     *time.Time
	UpdatedAt   time.Time `Gorm:"autoUpdateTime"`
}
