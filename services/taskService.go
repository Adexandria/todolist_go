package services

import (
	"SM/models"
	"SM/repositories"
)

type TaskService struct {
	repo repositories.Tasker
}

func TaskServiceCon(taskRepo *repositories.TaskRepository) *TaskService {
	return &TaskService{
		repo: taskRepo,
	}
}

var _ ITaskService = (*TaskService)(nil)

func (t *TaskService) CreateTask(dto *models.CreateTaskDTO) (uint, string) {
	return t.repo.Create(dto)
}

func (t *TaskService) UpdateTask(id int, dto *models.UpdateTaskDTO) bool {
	currentTask := t.repo.Get(id)

	if currentTask == (models.TaskDTO{}) {
		return false
	}
	return t.repo.Update(id, dto)
}

func (t *TaskService) DeleteTask(id int) bool {
	currentTask := t.repo.Get(id)

	if currentTask == (models.TaskDTO{}) {
		return false
	}
	return t.repo.Delete(id)
}

func (t *TaskService) GetAllTasks(page int, pageSize int) models.PaginatedDTO {
	return t.repo.GetAll(page, pageSize)
}

func (t *TaskService) GetTaskById(id int) models.TaskDTO {
	return t.repo.Get(id)
}

func (t *TaskService) GetTaskByMonthAndYear(month int, year int, page int, pageSize int) models.PaginatedDTO {
	return t.repo.GetByMonthAndYear(month, year, page, pageSize)
}

func (t *TaskService) GetTaskByCreatedDate(date int, month int, year int, page int, pageSize int) models.PaginatedDTO {
	return t.repo.GetByCreatedDate(date, month, year, page, pageSize)
}
func (t *TaskService) GetTaskByYear(year int, page int, pageSize int) models.PaginatedDTO {
	return t.repo.GetByYear(year, page, pageSize)
}

func (t *TaskService) SearchTaskByName(name string, page int, pageSize int) models.PaginatedDTO {
	return t.repo.SearchByName(name, page, pageSize)
}
