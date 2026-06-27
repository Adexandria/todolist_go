package repositories

import "SM/models"

type IUserRepository interface {
	GetUserById(id int) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)

	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error

	VerifyUserByEmail(email string) (bool, error)
	VerifyUserById(id int) (bool, error)
}
