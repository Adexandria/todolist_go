package services

import "SM/models"

type ITaskService interface {
	SearchTaskByName(name string, page int, pageSize int) models.PaginatedDTO
	CreateTask(dto *models.CreateTaskDTO) (uint, string)
	UpdateTask(id int, dto *models.UpdateTaskDTO) bool
	DeleteTask(id int) bool
	GetAllTasks(page int, pageSize int) models.PaginatedDTO
	GetTaskById(id int) models.TaskDTO
	GetTaskByMonthAndYear(month int, year int, page int, pageSize int) models.PaginatedDTO
	GetTaskByCreatedDate(date int, month int, year int, page int, pageSize int) models.PaginatedDTO
	GetTaskByYear(year int, page int, pageSize int) models.PaginatedDTO
}
