package repositories

import "SM/models"

// Tasker defines the interface for task-related database operations
type Tasker interface {

	// Creates a new task
	Create(dto *models.Task) (uint, string)

	// Updates an existing task by its ID
	Update(id int, dto *models.Task) bool

	// Deletes a task by its ID
	Delete(id int) bool

	// Retrieves a task by its ID
	Get(id int) models.Task

	// Retrieves all tasks with pagination
	GetAll(page int, pageSize int) []models.Task

	// Retrieves tasks based on the specified year with pagination
	GetByYear(year int, page int, pageSize int) []models.Task

	// Retrieves tasks based on the specified month and year with pagination
	GetByMonthAndYear(month int, year int, page int, pageSize int) []models.Task

	// Retrieves tasks based on the specified created date, month, and year with pagination
	GetByCreatedDate(day int, month int, year int, page int, pageSize int) []models.Task

	// Searches for tasks by name with pagination
	SearchByName(name string, page int, pageSize int) []models.Task
}
