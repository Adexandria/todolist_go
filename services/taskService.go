package services

import (
	"SM/models"
	"SM/repositories"
	"fmt"
	"time"
)

type TaskService struct {
	repo repositories.Tasker
}

// Injects the TaskRepository dependency into the TaskService
func TaskServiceCon(taskRepo *repositories.TaskRepository) *TaskService {
	return &TaskService{
		repo: taskRepo,
	}
}

// Ensures that TaskService implements the ITaskService interface
var _ ITaskService = (*TaskService)(nil)

// Creates a new task
func (t *TaskService) CreateTask(dto *models.CreateTaskDTO) (uint, string) {
	newTask := &models.Task{
		Description: dto.Description,
		Name:        dto.Name,
		CreatedAt:   time.Now(),
	}
	if dto.DueDate != "" {
		t, err := time.Parse("2006-01-02 15:04:05", dto.DueDate)
		if err != nil {
			return 0, "Failed to parse due date"
		}
		newTask.DueDate = &t
	}

	return t.repo.Create(newTask)
}

// Updates existing task by id
func (t *TaskService) UpdateTask(id int, dto *models.UpdateTaskDTO) bool {
	currentTask := t.repo.Get(id)

	if currentTask == (models.Task{}) {
		return false
	}

	updateTask := &models.Task{
		Description: dto.Description,
		Name:        dto.Name,
	}

	if dto.DueDate != "" {
		t, err := time.Parse("2006-01-02 15:04:05", dto.DueDate)
		if err != nil {
			return false
		}
		updateTask.DueDate = &t
	}

	return t.repo.Update(id, updateTask)
}

// Deletes a task by id
func (t *TaskService) DeleteTask(id int) bool {
	currentTask := t.repo.Get(id)

	if currentTask == (models.Task{}) {
		return false
	}
	return t.repo.Delete(id)
}

// Retrieves all tasks with pagination
func (t *TaskService) GetAllTasks(page int, pageSize int) models.PaginatedDTO {
	tasks := t.repo.GetAll(page, pageSize)

	var total int

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, page, pageSize)
}

// Retrieves a task by its id
func (t *TaskService) GetTaskById(id int) models.TaskDTO {
	currentTask := t.repo.Get(id)

	var dueDate string

	if currentTask.DueDate != nil {
		dueDate = currentTask.DueDate.Format("2006-01-02")
	}
	return models.TaskDTO{
		ID:          fmt.Sprintf("%d", currentTask.ID),
		Description: currentTask.Description,
		Name:        currentTask.Name,
		DueDate:     dueDate,
		CreatedDate: currentTask.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// Retrieves tasks based on the specified month and year with pagination
func (t *TaskService) GetTaskByMonthAndYear(month int, year int, page int, pageSize int) models.PaginatedDTO {
	tasks := t.repo.GetByMonthAndYear(month, year, page, pageSize)

	var total int

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, page, pageSize)
}

// Retrieves tasks based on the created date with pagination
func (t *TaskService) GetTaskByCreatedDate(date int, month int, year int, page int, pageSize int) models.PaginatedDTO {
	tasks := t.repo.GetByCreatedDate(date, month, year, page, pageSize)

	var total int

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, page, pageSize)
}

// Retrieves tasks based on the specified year with pagination
func (t *TaskService) GetTaskByYear(year int, page int, pageSize int) models.PaginatedDTO {
	tasks := t.repo.GetByYear(year, page, pageSize)

	var total int

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, page, pageSize)
}

// Searches for tasks by name with paginations
func (t *TaskService) SearchTaskByName(name string, page int, pageSize int) models.PaginatedDTO {
	tasks := t.repo.SearchByName(name, page, pageSize)

	var total int

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, page, pageSize)
}
