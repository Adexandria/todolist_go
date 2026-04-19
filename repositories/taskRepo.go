package repositories

import (
	"SM/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TaskRepository struct {
	Db *gorm.DB
}

func TaskRepo(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		Db: db,
	}
}

var _ Tasker = (*TaskRepository)(nil)

func (r *TaskRepository) Create(dto *models.CreateTaskDTO) (uint, string) {
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

	ctx := context.Background()

	err := gorm.G[models.Task](r.Db).Create(ctx, newTask)

	if err != nil {
		return 0, err.Error()
	}

	return newTask.ID, ""
}

func (r *TaskRepository) Update(id int, dto *models.UpdateTaskDTO) bool {
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

	ctx := context.Background()

	_, err := gorm.G[models.Task](r.Db).Where("id = ?", id).Updates(ctx, *updateTask)

	if err != nil {
		return false
	}

	return true
}

func (r *TaskRepository) Delete(id int) bool {
	ctx := context.Background()
	_, err := gorm.G[models.Task](r.Db).Where("id = ?", id).Delete(ctx)
	return err == nil
}

func (r *TaskRepository) Get(id int) models.TaskDTO {

	ctx := context.Background()
	currentTask, err := gorm.G[models.Task](r.Db).Where("id = ?", id).First(ctx)
	if err != nil {
		return models.TaskDTO{}
	}

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
func (r *TaskRepository) GetAll(page int, pageSize int) models.PaginatedDTO {
	var total int

	ctx := context.Background()

	tasks, err := gorm.G[models.Task](r.Db).Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at").Find(ctx)

	if err != nil {
		fmt.Println(err)
		return models.PaginatedDTO{}
	}

	total = len(tasks)

	return models.ToPaginatedDTO(tasks, total, pageSize, pageSize)
}

func (r *TaskRepository) GetByYear(year int, page int, pageSize int) models.PaginatedDTO {

	ctx := context.Background()

	currentTasks, err := gorm.G[models.Task](r.Db).Limit(pageSize).Offset((page-1)*pageSize).
		Where("strftime('%Y', created_at) = ?", fmt.Sprintf("%d", year)).Find(ctx)

	if err != nil {
		return models.PaginatedDTO{}
	}

	total := len(currentTasks)

	return models.ToPaginatedDTO(currentTasks, total, pageSize, pageSize)
}

func (r *TaskRepository) GetByMonthAndYear(month int, year int, page int, pageSize int) models.PaginatedDTO {
	ctx := context.Background()

	currentTasks, err := gorm.G[models.Task](r.Db).Limit(pageSize).Offset((page-1)*pageSize).
		Where("strftime('%Y', created_at) = ? AND strftime('%m', created_at) = ?",
			fmt.Sprintf("%d", year),
			fmt.Sprintf("%02d", month),
		).Find(ctx)

	if err != nil {
		return models.PaginatedDTO{}
	}

	total := len(currentTasks)

	return models.ToPaginatedDTO(currentTasks, total, pageSize, pageSize)
}

func (r *TaskRepository) GetByCreatedDate(day int, month int, year int, page int, pageSize int) models.PaginatedDTO {
	ctx := context.Background()

	currentTasks, err := gorm.G[models.Task](r.Db).Limit(pageSize).Offset((page-1)*pageSize).
		Where("strftime('%Y-%m-%d', created_at) = ?",
			fmt.Sprintf("%d-%02d-%02d", year, month, day),
		).Find(ctx)

	if err != nil {
		return models.PaginatedDTO{}
	}

	total := len(currentTasks)

	return models.ToPaginatedDTO(currentTasks, total, pageSize, pageSize)

}

func (r *TaskRepository) SearchByName(name string, page int, pageSize int) models.PaginatedDTO {

	ctx := context.Background()

	currentTasks, err := gorm.G[models.Task](r.Db).Limit(pageSize).Offset((page-1)*pageSize).
		Where("LOWER(name) LIKE LOWER(?)", name+"%").
		Find(ctx)

	if err != nil {
		return models.PaginatedDTO{}
	}

	total := len(currentTasks)

	return models.ToPaginatedDTO(currentTasks, total, pageSize, pageSize)
}
