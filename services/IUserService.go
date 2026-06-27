package services

import "SM/models"

type IUserService interface {
	CreateUser(newUser *models.CreateUserDTO) *ActionResult

	CreateAdmin(id int, newUser *models.CreateUserDTO) *ActionResult

	GetUserById(id int) *ActionResultModel[*models.UserDTO]

	DeleteUser(id int) *ActionResult

	SetEmail(id int, email string) *ActionResult

	ConfirmEmail(id int) *ActionResult

	SetPhoneNumber(id int, phone string) *ActionResult

	ConfirmPhoneNumber(id int) *ActionResult

	SetUsername(id int, username string) *ActionResult

	SetFirstName(id int, firstName string) *ActionResult

	SetLastName(id int, lastName string) *ActionResult

	ChangePassword(id int, oldPassword string, newPassword string) *ActionResult

	ChangeAnonPassword(email string, newPassword string) *ActionResult

	ChangeRole(id int, role string) *ActionResult

	EnableLockout(id int) *ActionResult

	DisableLockout(id int) *ActionResult
}
