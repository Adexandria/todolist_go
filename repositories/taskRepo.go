package repositories

import (
	"SM/models"
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type TaskRepository struct {
	Db  *gorm.DB
	Log *slog.Logger
}

// Constructor function to create a new TaskRepository instance with the provided database connection

func TaskRepo(db *gorm.DB, jsonHandler *slog.JSONHandler) *TaskRepository {
	return &TaskRepository{
		Db:  db,
		Log: slog.New(jsonHandler),
	}
}

// Ensures that TaskRepository implements the Tasker interface
var _ Tasker = (*TaskRepository)(nil)

// Creates a new task
func (r *TaskRepository) Create(newTask *models.Task) (uint, string) {

	ctx := context.Background()

	err := gorm.G[models.Task](r.Db).Create(ctx, newTask)

	if err != nil {
		r.Log.Error(err.Error())
		return 0, err.Error()
	}

	return newTask.ID, ""
}

// Updates an existing task by its ID
func (r *TaskRepository) Update(id int, updateTask *models.Task) bool {

	ctx := context.Background()

	_, err := gorm.G[models.Task](r.Db).Where("id = ?", id).Updates(ctx, *updateTask)

	if err != nil {
		r.Log.Error("Failed to update task: " + err.Error())
		return false
	}

	return true
}

// Deletes a task by its ID
func (r *TaskRepository) Delete(id int) bool {
	ctx := context.Background()
	_, err := gorm.G[models.Task](r.Db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		r.Log.Error("Failed to delete task: " + err.Error())
		return false
	}

	return true
}

// Retrieves a task by its ID
func (r *TaskRepository) Get(id int) models.Task {
	ctx := context.Background()
	currentTask, err := gorm.G[models.Task](r.Db).Where("id = ?", id).First(ctx)
	if err != nil {
		r.Log.Error("Failed to get task: " + err.Error())
		return models.Task{}
	}

	return currentTask
}

// Retrieves all tasks with pagination
func (r *TaskRepository) GetAll(page int, pageSize int) []models.Task {
	ctx := context.Background()

	tasks, err := gorm.G[models.Task](r.Db).Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at").Find(ctx)

	if err != nil {
		r.Log.Error("Failed to get all tasks" + err.Error())
		return []models.Task{}
	}
	return tasks

}

// Retrieves tasks based on the specified year with pagination
func (r *TaskRepository) GetByYear(year int, page int, pageSize int) []models.Task {
	ctx := context.Background()

	currentTasks, err := gorm.G[models.Task](r.Db).Limit(pageSize).Offset((page-1)*pageSize).
		Where("strftime('%Y', created_at) = ?", fmt.Sprintf("%d", year)).Find(ctx)

	if err != nil {
		r.Log.Error("Failed to get task by year" + err.Error())
		return []models.Task{}
	}

	return currentTasks
}

// Retrieves tasks based on the specified month and year with pagination
func (r *TaskRepository) GetByMonthAndYear(month int, year int, page int, pageSize int) []models.Task {
	ctx := context.Background()

	currentTasks, err := gorm.G[models.Task](r.Db).Limit(pageSize).Offset((page-1)*pageSize).
		Where("strftime('%Y', created_at) = ? AND strftime('%m', created_at) = ?",
			fmt.Sprintf("%d", year),
			fmt.Sprintf("%02d", month),
		).Find(ctx)

	if err != nil {
		r.Log.Error("Failed to get task by month and year" + err.Error())
		return []models.Task{}
	}

	return currentTasks
}

// Retrieves tasks based on the specified created date, month, and year with pagination
func (r *TaskRepository) GetByCreatedDate(day int, month int, year int, page int, pageSize int) []models.Task {
	ctx := context.Background()

	currentTasks, err := gorm.G[models.Task](r.Db).Limit(pageSize).Offset((page-1)*pageSize).
		Where("strftime('%Y-%m-%d', created_at) = ?",
			fmt.Sprintf("%d-%02d-%02d", year, month, day),
		).Find(ctx)

	if err != nil {
		r.Log.Error("Failed to get task by date" + err.Error())
		return []models.Task{}
	}

	return currentTasks

}

// Searches for tasks by name with pagination
func (r *TaskRepository) SearchByName(name string, page int, pageSize int) []models.Task {

	ctx := context.Background()

	currentTasks, err := gorm.G[models.Task](r.Db).Limit(pageSize).Offset((page-1)*pageSize).
		Where("LOWER(name) LIKE LOWER(?)", name+"%").
		Find(ctx)

	if err != nil {
		r.Log.Error("Failed to get task by name" + err.Error())
		return []models.Task{}
	}

	return currentTasks
}
