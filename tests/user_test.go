package tests

import (
	"SM/models"
	"SM/services"
	"SM/tests/mocks"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.CreateUserDTO{
		FirstName:      "John",
		LastName:       "Doe",
		Username:       "johndoe",
		Password:       "password",
		RetypePassword: "password",
		Email:          "email@email",
	}

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.CreateUser(&dto)

	assert.Empty(t, result.Error)
	assert.Equal(t, result.IsSuccess, true)
}

func TestFailedToCreateUserIfEmailExist(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	currentUser := models.User{
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

	dto := models.CreateUserDTO{
		FirstName:      "John",
		LastName:       "Doe",
		Username:       "johndoe",
		Password:       "password",
		RetypePassword: "password",
		Email:          "email@email",
	}

	_ = database.CreateUser(&currentUser)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.CreateUser(&dto)

	assert.Equal(t, result.Error[0], "Email already exists")
	assert.Equal(t, result.IsSuccess, false)
}

func TestFailedToCreateUserIfUsernameExist(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	currentUser := models.User{
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

	dto := models.CreateUserDTO{
		FirstName:      "John",
		LastName:       "Doe",
		Username:       "johndoe",
		Password:       "password",
		RetypePassword: "password",
		Email:          "email12@email",
	}

	_ = database.CreateUser(&currentUser)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.CreateUser(&dto)

	assert.Equal(t, result.Error[0], "Username already exists")
	assert.Equal(t, result.IsSuccess, false)
}

func TestCreateAdmin(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	adminId := database.AddNewAdmin()

	dto := models.CreateUserDTO{
		FirstName:      "John",
		LastName:       "Doe",
		Username:       "johndoe",
		Password:       "password",
		RetypePassword: "password",
		Email:          "email@email",
	}

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.CreateAdmin(adminId, &dto)

	assert.Empty(t, result.Error)
	assert.Equal(t, result.IsSuccess, true)
}

func TestFailedToCreateAdminIfAdminNotExist(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.CreateUserDTO{
		FirstName:      "John",
		LastName:       "Doe",
		Username:       "johndoe",
		Password:       "password",
		RetypePassword: "password",
		Email:          "email@email",
	}

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.CreateAdmin(1, &dto)

	assert.Equal(t, result.Error[0], "Failed to create user")
	assert.Equal(t, result.IsSuccess, false)
}

func TestGetUser(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	currentUser := useCase.GetUserById(1)

	assert.Equal(t, dto.FirstName, currentUser.Data.FirstName)
	assert.Equal(t, dto.LastName, currentUser.Data.LastName)
	assert.Equal(t, dto.Username, currentUser.Data.Username)
	assert.Equal(t, dto.Email, currentUser.Data.Email)
}

func TestFailedToCreateUser(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.GetUserById(1)

	assert.Equal(t, result.Error[0], "User not found")
	assert.Equal(t, result.IsSuccess, false)
}

func TestDeleteUser(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.DeleteUser(1)

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToDeleteUserIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.DeleteUser(1)

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestSetEmail(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.SetEmail(1, "johndoe@email")

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToSetEmailIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.SetEmail(1, "johndoe@email")

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestConfirmEmail(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.ConfirmEmail("email@email")

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToConfirmEmailIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.ConfirmEmail("johndoe@email")

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestSetPhoneNumber(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.SetPhoneNumber(1, "00000000001")

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToSetPhonenumberIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.SetPhoneNumber(1, "00000000001")

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestConfirmPhoneNumber(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.ConfirmPhoneNumber(1)

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToConfirmPhonenumberIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.ConfirmPhoneNumber(1)

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestSetUsername(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.SetUsername(1, "john_doe")

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToSetUsernameIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.SetUsername(1, "johnny")

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestSetFirstName(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.SetFirstName(1, "Johnny")

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToSetFirstNameIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.SetFirstName(1, "johns")

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestSetLastName(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.SetLastName(1, "Doee")

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToSetLastNameIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.SetLastName(1, "Doeey")

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestChangePassword(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.ChangePassword(1, "password", "password1234")

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToChangePasswordIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.ChangePassword(1, "password", "password1234")

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "Failed to change password")
}

func TestChangeAnonymousPassword(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.ChangeAnonPassword("email@email", "password1234")

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToChangeAnonyPasswordIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.ChangeAnonPassword("email@email", "password1234")

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestChangeRole(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.ChangeRole(1, "admin")

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToChangeRoleIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.ChangeRole(1, "user")

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestEnableLockout(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.EnableLockout("johndoe")

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToEnableLockoutIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.EnableLockout("johnnyyy")

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}

func TestDisableLockout(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	dto := models.User{
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

	_ = database.CreateUser(&dto)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.DisableLockout(1)

	assert.Equal(t, result.IsSuccess, true)
	assert.Empty(t, result.Error)
}

func TestFailedToDisableLockoutIfUserNotFound(t *testing.T) {
	database := mocks.NewUserDatabase()

	handler := slog.NewJSONHandler(os.Stderr, nil)

	passwordManager := mocks.NewPasswordHandler(handler)

	useCase := services.NewUserService(database, passwordManager)

	result := useCase.DisableLockout(1)

	assert.Equal(t, result.IsSuccess, false)
	assert.Equal(t, result.Error[0], "User not found")
}
