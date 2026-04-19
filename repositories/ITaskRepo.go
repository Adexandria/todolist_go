package repositories

import "SM/models"

type Tasker interface {
	Create(dto *models.CreateTaskDTO) (uint, string)
	Update(id int, dto *models.UpdateTaskDTO) bool
	Delete(id int) bool
	Get(id int) models.TaskDTO
	GetAll(page int, pageSize int) models.PaginatedDTO
	GetByYear(year int, page int, pageSize int) models.PaginatedDTO
	GetByMonthAndYear(month int, year int, page int, pageSize int) models.PaginatedDTO
	GetByCreatedDate(day int, month int, year int, page int, pageSize int) models.PaginatedDTO
	SearchByName(name string, page int, pageSize int) models.PaginatedDTO
}
