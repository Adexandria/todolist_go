package services

import "SM/models"

// ITaskService defines the interface for task-related operations
// It includes methods for searching, creating, updating, deleting, and retrieving tasks based on various criteria.
// This interface allows for abstraction and decoupling of the task service implementation from its usage, enabling easier testing and maintenance.
type ITaskService interface {

	// Searches for tasks by name with pagination
	SearchTaskByName(userId int, name string, page int, pageSize int) models.PaginatedDTO

	// Creates a new task
	CreateTask(userId int, dto *models.CreateTaskDTO) (uint, string)

	// Updates an existing task by its ID
	UpdateTask(userId int, id int, dto *models.UpdateTaskDTO) bool

	// Deletes a task by its ID
	DeleteTask(userId int, id int) bool

	// Retrieves all tasks with pagination
	GetAllTasks(userId int, page int, pageSize int) models.PaginatedDTO

	// Retrieves a task by its ID
	GetTaskById(userId int, id int) models.TaskDTO

	// Retrieves tasks based on the specified month and year with pagination
	GetTaskByMonthAndYear(userId int, month int, year int, page int, pageSize int) models.PaginatedDTO

	// Retrieves tasks based on the specified created date, month, and year with pagination
	GetTaskByCreatedDate(userId int, date int, month int, year int, page int, pageSize int) models.PaginatedDTO

	// Retrieves tasks based on the specified year with pagination
	GetTaskByYear(userId int, year int, page int, pageSize int) models.PaginatedDTO
}
