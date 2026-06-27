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
func (t *TaskService) CreateTask(userId int, dto *models.CreateTaskDTO) (uint, string) {
	newTask := &models.Task{
		Description: dto.Description,
		Name:        dto.Name,
		CreatedAt:   time.Now(),
		UserID:      uint(userId),
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
func (t *TaskService) UpdateTask(userId int, id int, dto *models.UpdateTaskDTO) bool {
	currentTask := t.repo.Get(userId, id)

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
func (t *TaskService) DeleteTask(userId int, id int) bool {
	currentTask := t.repo.Get(userId, id)

	if currentTask == (models.Task{}) {
		return false
	}
	return t.repo.Delete(userId, id)
}

// Retrieves all tasks with pagination
func (t *TaskService) GetAllTasks(userId int, page int, pageSize int) models.PaginatedDTO {
	tasks := t.repo.GetAll(userId, page, pageSize)

	var total int

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, page, pageSize)
}

// Retrieves a task by its id
func (t *TaskService) GetTaskById(userId int, id int) models.TaskDTO {
	currentTask := t.repo.Get(userId, id)

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
func (t *TaskService) GetTaskByMonthAndYear(userId int, month int, year int, page int, pageSize int) models.PaginatedDTO {
	tasks := t.repo.GetByMonthAndYear(userId, month, year, page, pageSize)

	var total int

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, page, pageSize)
}

// Retrieves tasks based on the created date with pagination
func (t *TaskService) GetTaskByCreatedDate(userId int, date int, month int, year int, page int, pageSize int) models.PaginatedDTO {
	tasks := t.repo.GetByCreatedDate(userId, date, month, year, page, pageSize)

	var total int

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, page, pageSize)
}

// Retrieves tasks based on the specified year with pagination
func (t *TaskService) GetTaskByYear(userId int, year int, page int, pageSize int) models.PaginatedDTO {
	tasks := t.repo.GetByYear(userId, year, page, pageSize)

	var total int

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, page, pageSize)
}

// Searches for tasks by name with paginations
func (t *TaskService) SearchTaskByName(userId int, name string, page int, pageSize int) models.PaginatedDTO {
	tasks := t.repo.SearchByName(userId, name, page, pageSize)

	var total int

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, page, pageSize)
}
