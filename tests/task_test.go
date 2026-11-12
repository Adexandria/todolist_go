package tests

import (
	"SM/models"
	"SM/services"
	"SM/tests/mocks"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentTasks(t *testing.T) {

	database := mocks.NewMockDatabase()

	dto := models.Task{
		ID:          uint(1),
		Name:        "test",
		Description: "test",
		UserID:      uint(1),
		DueDate:     new(time.Now().Add(50 * 60 * 60)),
	}

	userDatabase := mocks.NewUserDatabase()

	userdto := models.User{
		FirstName:            "John",
		LastName:             "Doe",
		Username:             "johndoe",
		Password:             "password",
		Email:                "email@email",
		PhoneNumberConfirmed: true,
		EmailConfirmed:       true,
		PhoneNumber:          "00000000000",
		Role:                 0,
		AuthenticationType:   2,
		AuthenticationKey:    "",
		RefreshToken:         "",
		Tasks:                []models.Task{},
	}

	_ = userDatabase.CreateUser(&userdto)

	database.Create(&dto)

	useCase := services.NewTaskService(database)

	currentTask := useCase.GetTaskById(1, 1)

	assert.Equal(t, strconv.Itoa(int(dto.ID)), currentTask.ID)
	assert.Equal(t, dto.Name, currentTask.Name)
	assert.Equal(t, dto.Description, currentTask.Description)

}

func TestGetPaginatedTasks(t *testing.T) {
	database := mocks.NewMockDatabase()

	userDatabase := mocks.NewUserDatabase()

	userdto := models.User{
		FirstName:            "John",
		LastName:             "Doe",
		Username:             "johndoe",
		Password:             "password",
		Email:                "email@email",
		PhoneNumberConfirmed: true,
		EmailConfirmed:       true,
		PhoneNumber:          "00000000000",
		Role:                 0,
		AuthenticationType:   2,
		AuthenticationKey:    "",
		RefreshToken:         "",
		Tasks:                []models.Task{},
	}

	_ = userDatabase.CreateUser(&userdto)

	dto := models.Task{
		ID:          uint(1),
		Name:        "test",
		Description: "test",
		UserID:      uint(1),
		DueDate:     new(time.Now().Add(50 * 60 * 60)),
	}
	dto2 := models.Task{
		ID:          uint(2),
		Name:        "test",
		Description: "test",
		UserID:      uint(1),
		DueDate:     new(time.Now().Add(50 * 60 * 60)),
	}
	database.Create(&dto)
	database.Create(&dto2)

	useCase := services.NewTaskService(database)
	userTasks := useCase.GetAllTasks(1, 1, 20)

	assert.Equal(t, userTasks.Total, 2)
	assert.Equal(t, strconv.Itoa(int(dto.ID)), userTasks.Tasks[0].ID)
	assert.Equal(t, userTasks.Tasks[1].ID, strconv.Itoa(int(dto2.ID)))
}

func TestSearchTasks(t *testing.T) {
	database := mocks.NewMockDatabase()
	userDatabase := mocks.NewUserDatabase()
	userdto := models.User{
		FirstName:            "John",
		LastName:             "Doe",
		Username:             "johndoe",
		Password:             "password",
		Email:                "email@email",
		PhoneNumberConfirmed: true,
		EmailConfirmed:       true,
		PhoneNumber:          "00000000000",
		Role:                 0,
		AuthenticationType:   2,
		AuthenticationKey:    "",
		RefreshToken:         "",
		Tasks:                []models.Task{},
	}

	_ = userDatabase.CreateUser(&userdto)
	dto := models.Task{
		ID:          uint(1),
		Name:        "test",
		Description: "test",
		UserID:      uint(1),
		DueDate:     new(time.Now().Add(50 * 60 * 60)),
	}

	database.Create(&dto)
	useCase := services.NewTaskService(database)
	userTasks := useCase.SearchTaskByName(1, "test", 1, 20)

	assert.Equal(t, userTasks.Total, 1)
	assert.Equal(t, strconv.Itoa(int(dto.ID)), userTasks.Tasks[0].ID)
}

func TestCreateTask(t *testing.T) {
	database := mocks.NewMockDatabase()

	userDatabase := mocks.NewUserDatabase()

	userdto := models.User{
		FirstName:            "John",
		LastName:             "Doe",
		Username:             "johndoe",
		Password:             "password",
		Email:                "email@email",
		PhoneNumberConfirmed: true,
		EmailConfirmed:       true,
		PhoneNumber:          "00000000000",
		Role:                 0,
		AuthenticationType:   2,
		AuthenticationKey:    "",
		RefreshToken:         "",
		Tasks:                []models.Task{},
	}

	_ = userDatabase.CreateUser(&userdto)

	useCase := services.NewTaskService(database)

	dto := models.CreateTaskDTO{
		Name:        "test",
		Description: "test",
		DueDate:     time.Now().Add(50 * time.Hour).Format("2006-01-02 15:04:05"),
	}

	id, err := useCase.CreateTask(1, &dto)

	assert.Equal(t, err, "")
	assert.Equal(t, id, uint(1))
}

func TestDeleteTask(t *testing.T) {
	database := mocks.NewMockDatabase()
	userDatabase := mocks.NewUserDatabase()
	userdto := models.User{
		FirstName:            "John",
		LastName:             "Doe",
		Username:             "johndoe",
		Password:             "password",
		Email:                "email@email",
		PhoneNumberConfirmed: true,
		EmailConfirmed:       true,
		PhoneNumber:          "00000000000",
		Role:                 0,
		AuthenticationType:   2,
		AuthenticationKey:    "",
		RefreshToken:         "",
		Tasks:                []models.Task{},
	}

	_ = userDatabase.CreateUser(&userdto)

	dto := models.Task{
		ID:          uint(1),
		Name:        "test",
		Description: "test",
		UserID:      uint(1),
		DueDate:     new(time.Now().Add(50 * 60 * 60)),
	}

	database.Create(&dto)
	useCase := services.NewTaskService(database)

	isDeleted := useCase.DeleteTask(1, 1)

	assert.Equal(t, isDeleted, true)
}

func TestUpdateTask(t *testing.T) {
	database := mocks.NewMockDatabase()
	userDatabase := mocks.NewUserDatabase()
	userdto := models.User{
		FirstName:            "John",
		LastName:             "Doe",
		Username:             "johndoe",
		Password:             "password",
		Email:                "email@email",
		PhoneNumberConfirmed: true,
		EmailConfirmed:       true,
		PhoneNumber:          "00000000000",
		Role:                 0,
		AuthenticationType:   2,
		AuthenticationKey:    "",
		RefreshToken:         "",
		Tasks:                []models.Task{},
	}

	_ = userDatabase.CreateUser(&userdto)
	dto := models.Task{
		ID:          uint(1),
		Name:        "test",
		Description: "test",
		UserID:      uint(1),
		DueDate:     new(time.Now().Add(50 * 60 * 60)),
	}

	database.Create(&dto)
	useCase := services.NewTaskService(database)

	updateTask := models.UpdateTaskDTO{
		Name:        "test",
		Description: "test",
		DueDate:     time.Now().Add(50 * time.Hour).Format("2006-01-02 15:04:05"),
	}

	isUpdated := useCase.UpdateTask(1, 1, &updateTask)

	assert.Equal(t, isUpdated, true)
}

func TestGetTaskByYear(t *testing.T) {
	database := mocks.NewMockDatabase()
	userDatabase := mocks.NewUserDatabase()
	userdto := models.User{
		FirstName:            "John",
		LastName:             "Doe",
		Username:             "johndoe",
		Password:             "password",
		Email:                "email@email",
		PhoneNumberConfirmed: true,
		EmailConfirmed:       true,
		PhoneNumber:          "00000000000",
		Role:                 0,
		AuthenticationType:   2,
		AuthenticationKey:    "",
		RefreshToken:         "",
		Tasks:                []models.Task{},
	}

	_ = userDatabase.CreateUser(&userdto)
	dto := models.Task{
		ID:          uint(1),
		Name:        "test",
		Description: "test",
		UserID:      uint(1),
		DueDate:     new(time.Now().Add(50 * 60 * 60)),
		CreatedAt:   time.Now(),
	}

	database.Create(&dto)
	useCase := services.NewTaskService(database)

	currentTime := time.Now()

	tasks := useCase.GetTaskByYear(1, currentTime.Year(), 1, 20)

	assert.Equal(t, 1, tasks.Total)
	assert.Equal(t, tasks.Tasks[0].ID, strconv.Itoa(int(dto.ID)))
}

func TestGetTaskByMonth(t *testing.T) {
	database := mocks.NewMockDatabase()
	userDatabase := mocks.NewUserDatabase()
	userdto := models.User{
		FirstName:            "John",
		LastName:             "Doe",
		Username:             "johndoe",
		Password:             "password",
		Email:                "email@email",
		PhoneNumberConfirmed: true,
		EmailConfirmed:       true,
		PhoneNumber:          "00000000000",
		Role:                 0,
		AuthenticationType:   2,
		AuthenticationKey:    "",
		RefreshToken:         "",
		Tasks:                []models.Task{},
	}

	_ = userDatabase.CreateUser(&userdto)
	dto := models.Task{
		ID:          uint(1),
		Name:        "test",
		Description: "test",
		UserID:      uint(1),
		DueDate:     new(time.Now().Add(50 * 60 * 60)),
		CreatedAt:   time.Now(),
	}
	database.Create(&dto)

	useCase := services.NewTaskService(database)

	currentTime := time.Now()

	tasks := useCase.GetTaskByMonthAndYear(1, int(time.Month(currentTime.Month())), currentTime.Year(), 1, 20)

	assert.Equal(t, tasks.Total, 1)
	assert.Equal(t, tasks.Tasks[0].ID, strconv.Itoa(int(dto.ID)))
}

func TestGetTaskByCreatedDate(t *testing.T) {
	database := mocks.NewMockDatabase()
	userDatabase := mocks.NewUserDatabase()
	userdto := models.User{
		FirstName:            "John",
		LastName:             "Doe",
		Username:             "johndoe",
		Password:             "password",
		Email:                "email@email",
		PhoneNumberConfirmed: true,
		EmailConfirmed:       true,
		PhoneNumber:          "00000000000",
		Role:                 0,
		AuthenticationType:   2,
		AuthenticationKey:    "",
		RefreshToken:         "",
		Tasks:                []models.Task{},
	}

	_ = userDatabase.CreateUser(&userdto)
	dto := models.Task{
		ID:          uint(1),
		Name:        "test",
		Description: "test",
		UserID:      uint(1),
		DueDate:     new(time.Now().Add(50 * 60 * 60)),
		CreatedAt:   time.Now(),
	}
	database.Create(&dto)

	useCase := services.NewTaskService(database)

	currentTime := time.Now()

	tasks := useCase.GetTaskByCreatedDate(1, currentTime.Day(), int(time.Month(currentTime.Month())), currentTime.Year(), 1, 20)

	assert.Equal(t, tasks.Total, 1)
	assert.Equal(t, tasks.Tasks[0].ID, strconv.Itoa(int(dto.ID)))
}
