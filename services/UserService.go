package services

import (
	"SM/models"
	"SM/repositories"
	"SM/repositories/Utilities"
)

type UserService struct {
	userRepos       repositories.IUserRepository
	passwordManager Utilities.IPasswordManager
}

func (u *UserService) CreateUser(newUser *models.CreateUserDTO) *ActionResult {
	passwordHashed, err := u.passwordManager.HashPassword(newUser.Password)
	if err != nil {
		return BadRequestResult("Failed to create user")
	}
	user := &models.User{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Password:  passwordHashed,
		Role:      0,
	}

	err = u.userRepos.CreateUser(user)

	if err != nil {
		return ServerErrorResult("Failed to create user")
	}

	return SuccessResult()
}

func (u *UserService) CreateAdmin(id int, newUser *models.CreateUserDTO) *ActionResult {
	_, err := u.userRepos.VerifyUserById(id)
	if err != nil {
		return BadRequestResult("Failed to create user")
	}

	passwordHashed, err := u.passwordManager.HashPassword(newUser.Password)
	if err != nil {
		return BadRequestResult("Failed to create user")
	}
	user := &models.User{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Password:  passwordHashed,
		Role:      1,
	}

	err = u.userRepos.CreateUser(user)

	if err != nil {
		return ServerErrorResult("Failed to create user")
	}

	return SuccessResult()
}

func (u *UserService) GetUserById(id int) *ActionResultModel[*models.UserDTO] {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundDataResult[*models.UserDTO]("User not found")
	}

	userDTO := &models.UserDTO{
		FirstName:          currentUser.FirstName,
		LastName:           currentUser.LastName,
		Username:           currentUser.Username,
		Email:              currentUser.Email,
		PhoneNumber:        currentUser.PhoneNumber,
		IsLockoutEnabled:   currentUser.LockoutEnabled,
		IsTwoFactorEnabled: currentUser.TwoFactorEnabled,
	}

	return SuccessDataResult(userDTO)
}

func (u *UserService) DeleteUser(id int) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	err = u.userRepos.DeleteUser(&currentUser)
	if err != nil {
		return ServerErrorResult("Failed to delete user")
	}
	return SuccessResult()
}

func (u *UserService) SetEmail(id int, email string) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	currentUser.Email = email

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to update user")
	}
	return SuccessResult()
}

func (u *UserService) ConfirmEmail(id int) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	currentUser.EmailConfirmed = true

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to update user")
	}
	return SuccessResult()
}

func (u *UserService) SetPhoneNumber(id int, phone string) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	currentUser.PhoneNumber = phone

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to update user")
	}
	return SuccessResult()
}

func (u *UserService) ConfirmPhoneNumber(id int) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	currentUser.PhoneNumberConfirmed = true

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to update user")
	}
	return SuccessResult()
}

func (u *UserService) SetUsername(id int, username string) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	currentUser.Username = username

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to update user")
	}
	return SuccessResult()
}

func (u *UserService) SetFirstName(id int, firstName string) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	currentUser.FirstName = firstName

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to update user")
	}
	return SuccessResult()
}

func (u *UserService) SetLastName(id int, lastName string) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	currentUser.LastName = lastName

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to update user")
	}
	return SuccessResult()
}

func (u *UserService) ChangePassword(id int, oldPassword string, newPassword string) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("Failed to change password")
	}

	isValid := u.passwordManager.VerifyPassword(oldPassword, newPassword)

	if !isValid {
		return BadRequestResult("Failed to change password")
	}

	newHashedPassword, err := u.passwordManager.HashPassword(newPassword)

	currentUser.Password = newHashedPassword

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to change password")
	}
	return SuccessResult()
}

func (u *UserService) ChangeAnonPassword(email string, newPassword string) *ActionResult {
	currentUser, err := u.userRepos.GetUserByEmail(email)
	if err != nil {
		return NotFoundResult("User not found")
	}
	newHashedPassword, err := u.passwordManager.HashPassword(newPassword)

	currentUser.Password = newHashedPassword

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to change password")
	}

	return SuccessResult()
}

func (u *UserService) ChangeRole(id int, role string) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	if role == models.Roles[0] {
		currentUser.Role = 0
	} else if role == models.Roles[1] {
		currentUser.Role = 1
	}

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to update user")
	}
	return SuccessResult()
}

func (u *UserService) EnableLockout(id int) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	currentUser.LockoutEnabled = true

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to update user")
	}
	return SuccessResult()
}

func (u *UserService) DisableLockout(id int) *ActionResult {
	currentUser, err := u.userRepos.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	currentUser.LockoutEnabled = false

	err = u.userRepos.UpdateUser(&currentUser)

	if err != nil {
		return ServerErrorResult("Failed to update user")
	}
	return SuccessResult()
}

func UserServiceCon(userRepository *repositories.UserRepository, passwordManager *Utilities.PasswordManager) *UserService {
	return &UserService{
		userRepos:       userRepository,
		passwordManager: passwordManager}
}

var _ IUserService = (*UserService)(nil)
