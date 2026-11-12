package mocks

import (
	"SM/models"
	"SM/repositories"
	"sort"
	"time"
)

type Database struct {
	data map[uint]*models.Task
}

func NewMockDatabase() *Database {
	return &Database{
		data: make(map[uint]*models.Task),
	}
}
func (d *Database) Create(dto *models.Task) (uint, string) {
	id := uint(len(d.data) + 1)
	d.data[id] = dto
	return id, ""
}

func (d *Database) Update(id int, dto *models.Task) bool {
	d.data[uint(id)] = dto
	return true
}

func (d *Database) Delete(userId int, id int) bool {
	delete(d.data, uint(id))
	return true
}

func (d *Database) Get(userId int, id int) models.Task {
	x := d.data[uint(id)]
	if x.UserID != uint(userId) && x.ID != uint(id) {
		return models.Task{}
	}
	return *x
}

func (d *Database) GetAll(userId int, page int, pageSize int) []models.Task {
	userTasks := make([]models.Task, 0)

	// 1. filter
	for _, task := range d.data {
		if task.UserID == uint(userId) {
			userTasks = append(userTasks, *task)
		}
	}
	return getPagination(userTasks, page, pageSize)
}

func getPagination(userTasks []models.Task, page int, pageSize int) []models.Task {

	// 2. stable order (map iteration is random)
	sort.Slice(userTasks, func(i, j int) bool {
		return userTasks[i].ID < userTasks[j].ID
	})

	// 3. paginate the filtered slice
	if page < 1 || pageSize < 1 {
		return userTasks[:0]
	}
	start := (page - 1) * pageSize
	if start >= len(userTasks) {
		return []models.Task{}
	}
	end := start + pageSize
	if end > len(userTasks) {
		end = len(userTasks)
	}
	return userTasks[start:end]
}

func (d *Database) GetByYear(userId int, year int, page int, pageSize int) []models.Task {
	userTasks := make([]models.Task, 0)

	// 1. filter
	for _, task := range d.data {
		if task.UserID == uint(userId) && task.CreatedAt.Year() == year {
			userTasks = append(userTasks, *task)
		}
	}
	return getPagination(userTasks, page, pageSize)
}

func (d *Database) GetByMonthAndYear(userId int, month int, year int, page int, pageSize int) []models.Task {
	userTasks := make([]models.Task, 0)

	// 1. filter
	for _, task := range d.data {
		if task.UserID == uint(userId) && task.CreatedAt.Year() == year && task.CreatedAt.Month() == time.Month(month) {
			userTasks = append(userTasks, *task)
		}
	}
	return getPagination(userTasks, page, pageSize)
}

func (d *Database) GetByCreatedDate(userId int, day int, month int, year int, page int, pageSize int) []models.Task {
	userTasks := make([]models.Task, 0)

	// 1. filter
	for _, task := range d.data {
		if task.UserID == uint(userId) && task.CreatedAt.Year() == year && task.CreatedAt.Month() == time.Month(month) && task.CreatedAt.Day() == day {
			userTasks = append(userTasks, *task)
		}
	}
	return getPagination(userTasks, page, pageSize)
}

func (d *Database) SearchByName(userId int, name string, page int, pageSize int) []models.Task {
	userTasks := make([]models.Task, 0)

	for _, task := range d.data {
		if task.UserID == uint(userId) && task.Name == name {
			userTasks = append(userTasks, *task)
		}
	}
	return getPagination(userTasks, page, pageSize)
}

var _ repositories.TaskRepository = (*Database)(nil)
