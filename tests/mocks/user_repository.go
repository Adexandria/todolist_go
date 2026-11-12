package mocks

import (
	"SM/models"
	"SM/repositories"
	"errors"
)

type UserDatabase struct {
	data map[uint]*models.User
}

func NewUserDatabase() *UserDatabase {
	return &UserDatabase{make(map[uint]*models.User)}
}

func (u *UserDatabase) AddNewAdmin() int {
	id := len(u.data) + 1

	userId := uint(id)

	user := models.User{
		ID:                   1,
		FirstName:            "John",
		LastName:             "Doe",
		Username:             "johndoe",
		Password:             "password",
		Email:                "admin@email",
		PhoneNumberConfirmed: true,
		EmailConfirmed:       true,
		PhoneNumber:          "00000000000",
		Role:                 1,
		AuthenticationType:   2,
		AuthenticationKey:    "",
		RefreshToken:         "",
		Tasks:                []models.Task{},
	}

	u.data[userId] = &user

	return id
}

func (u *UserDatabase) GetUserById(id int) (models.User, error) {
	currentUser := u.data[uint(id)]
	if currentUser == nil {
		return models.User{}, errors.New("User not found")
	}
	return *currentUser, nil
}

func (u *UserDatabase) VerifyUserByUsername(username string) (bool, error) {
	for _, user := range u.data {
		if user.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func (u *UserDatabase) GetUserByUsername(username string) (models.User, error) {
	for _, user := range u.data {
		if user.Username == username {
			return *user, nil
		}
	}
	return models.User{}, errors.New("User not found")
}

func (u *UserDatabase) GetUserByEmail(email string) (models.User, error) {
	for _, user := range u.data {
		if user.Email == email {
			return *user, nil
		}
	}
	return models.User{}, errors.New("User not found")
}

func (u *UserDatabase) GetUserByRefreshToken(refreshToken string) (models.User, error) {
	for _, user := range u.data {
		if user.RefreshToken == refreshToken {
			return *user, nil
		}
	}
	return models.User{}, errors.New("User not found")
}

func (u *UserDatabase) CreateUser(user *models.User) error {
	id := uint(len(u.data) + 1)

	u.data[id] = user

	return nil
}

func (u *UserDatabase) UpdateUser(user *models.User) error {
	u.data[user.ID] = user
	return nil
}

func (u *UserDatabase) DeleteUser(user *models.User) error {
	delete(u.data, user.ID)
	return nil
}

func (u *UserDatabase) VerifyUserByEmail(email string) (bool, error) {
	for _, user := range u.data {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (u *UserDatabase) VerifyUserById(id int) (bool, error) {
	for _, user := range u.data {
		if user.ID == uint(id) {
			return true, nil
		}
	}
	return false, errors.New("Failed to verify user")
}

var _ repositories.UserRepository = (*UserDatabase)(nil)
