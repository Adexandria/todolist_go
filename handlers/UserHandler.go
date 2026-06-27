package handlers

import (
	"SM/models"
	"SM/repositories/Utilities"
	"SM/services"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	IUserService           services.IUserService
	IAuthenticationService services.IAuthenticationService
	LockoutTrial           string
	IValidator             Utilities.IValidator
}

func UserHandlerCon(userService *services.UserService, authenticationService *services.AuthenticationService, validator *Utilities.Validator) *UserHandler {
	return &UserHandler{IUserService: userService,
		IAuthenticationService: authenticationService,
		LockoutTrial:           os.Getenv("LockoutTrial"),
		IValidator:             validator}
}

func (handler *UserHandler) SignUp(c *gin.Context) {
	var createUserDTO models.CreateUserDTO
	if err := c.ShouldBind(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to create user"))
	}

	validatorResult := handler.IValidator.IsNullOrEmpty(createUserDTO.Username).
		IsNullOrEmpty(createUserDTO.LastName).IsNullOrEmpty(createUserDTO.FirstName).
		VerifyPassword(createUserDTO.Password, createUserDTO.RetypePassword).
		CheckPasswordRule(createUserDTO.Password).VerifyEmail(createUserDTO.Email).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	result := handler.IUserService.CreateUser(&createUserDTO)
	c.JSON(result.StatusCode, result)
}

func (handler *UserHandler) Login(c *gin.Context) {
	var loginDTO models.LoginDTO

	if err := c.ShouldBind(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Invalid username/password"))
	}

	validatorResult := handler.IValidator.IsNullOrEmpty(loginDTO.Username).IsNullOrEmpty(loginDTO.Password).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	value := getLockoutCount(c)
	if value == getInt(handler.LockoutTrial) {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("This user has been locked out"))
	}

	result := handler.IAuthenticationService.Authenticate(loginDTO.Username, loginDTO.Password)
	if !result.IsSuccess {
		c.SetCookie("lockout", strconv.Itoa(value+1), -1, "/", "", true, true)
		c.JSON(result.StatusCode, result)
	}

	if result.IsTwoFactorEnabled {
		c.SetCookie("token", result.Token, 30, "/", "", true, true)

	}
	c.JSON(result.StatusCode, result)
}

func (handler *UserHandler) ChangePassword(c *gin.Context) {
	var changePasswordDTO models.ChangePasswordDTO
	if err := c.ShouldBind(&changePasswordDTO); err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to change password"))
	}

	validatorResult := handler.IValidator.CheckPasswordRule(changePasswordDTO.NewPassword).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	result := handler.IUserService.ChangePassword(userId, changePasswordDTO.OldPassword, changePasswordDTO.NewPassword)
	c.JSON(result.StatusCode, result)
}

func (handler *UserHandler) ChangeEmail(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	email := c.Query("email")

	validatorResult := handler.IValidator.IsNullOrEmpty(email).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	result := handler.IUserService.SetEmail(userId, email)

	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
	}

	tokenResult := handler.IAuthenticationService.GenerateEmailResetToken(email)

	c.JSON(result.StatusCode, tokenResult)
}

func (handler *UserHandler) ConfirmEmailToken(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	token := c.Query("token")

	validatorResult := handler.IValidator.IsNullOrEmpty(token).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	result := handler.IAuthenticationService.VerifyEmailResetToken(userId, token)

	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
	}
	confirmEmailToken := handler.IUserService.ConfirmEmail(userId)
	c.JSON(result.StatusCode, confirmEmailToken)
}

func (handler *UserHandler) ResetPassword(c *gin.Context) {
	email := c.Query("email")

	validatorResult := handler.IValidator.IsNullOrEmpty(email).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}
	result := handler.IAuthenticationService.GeneratePasswordResetToken(email)
	c.JSON(result.StatusCode, result)
}

func (handler *UserHandler) ResetPasswordByToken(c *gin.Context) {
	token := c.Query("token")
	var password models.AnonymousChangePasswordDTO
	if err := c.ShouldBind(&password); err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to change password"))
	}
	validatorResult := handler.IValidator.IsNullOrEmpty(password.NewPassword).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	result := handler.IAuthenticationService.DecodeToken(token)
	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
	}

	passwordChangedResult := handler.IUserService.ChangeAnonPassword(result.Data, password.NewPassword)

	c.JSON(result.StatusCode, passwordChangedResult)
}

func (handler *UserHandler) ChangePhonenumber(c *gin.Context) {
	phonenumber := c.Query("phonenumber")

	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	validatorResult := handler.IValidator.VerifyPhoneNumber(phonenumber).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	result := handler.IUserService.SetPhoneNumber(userId, phonenumber)

	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
	}

	tokenResult := handler.IAuthenticationService.GeneratePhoneNumberResetToken(userId)

	c.JSON(result.StatusCode, tokenResult)
}

func (handler *UserHandler) VerifyPhonenumberToken(c *gin.Context) {
	token := c.Query("token")
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	validatorResult := handler.IValidator.IsNullOrEmpty(token).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	result := handler.IAuthenticationService.VerifyPhoneNumberResetToken(userId, token)
	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
	}
	confirmPhonenumber := handler.IUserService.ConfirmPhoneNumber(userId)
	c.JSON(result.StatusCode, confirmPhonenumber)
}

// update user
func (handler *UserHandler) UpdateFirstName(c *gin.Context) {
	firstName := c.Query("firstname")
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	validatorResult := handler.IValidator.IsNullOrEmpty(firstName).Verify()
	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}
	result := handler.IUserService.SetFirstName(userId, firstName)
	c.JSON(result.StatusCode, result)
}

func (handler *UserHandler) UpdateLastName(c *gin.Context) {
	lastName := c.Query("lastname")
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	validatorResult := handler.IValidator.IsNullOrEmpty(lastName).Verify()
	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	result := handler.IUserService.SetLastName(userId, lastName)
	c.JSON(result.StatusCode, result)
}

func (handler *UserHandler) UpdateUsername(c *gin.Context) {
	username := c.Query("username")
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	validatorResult := handler.IValidator.IsNullOrEmpty(username).Verify()
	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	result := handler.IUserService.SetUsername(userId, username)
	c.JSON(result.StatusCode, result)
}

// delete user -user

func (handler *UserHandler) DeleteUserAccount(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)

	userId := claims["user_id"].(int)
	result := handler.IUserService.DeleteUser(userId)

	c.JSON(result.StatusCode, result)
}

// delete user - user
func (handler *UserHandler) DeleteUserById(c *gin.Context) {
	userId := c.Param("user_id")
	id, _ := strconv.Atoi(userId)

	result := handler.IUserService.DeleteUser(id)
	c.JSON(result.StatusCode, result)
}

// disable lockout user
func (handler *UserHandler) DisableUserLockout(c *gin.Context) {
	userId := c.Param("user_id")

	validatorResult := handler.IValidator.IsNullOrEmpty(userId).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	id, _ := strconv.Atoi(userId)

	result := handler.IUserService.DisableLockout(id)
	c.JSON(result.StatusCode, result)
}

// set up two factor authentication
func (handler *UserHandler) SetUpTwoFactorAuthentication(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	authenticationType := claims["authentication_type"].(string)

	result := handler.IAuthenticationService.EnableTwoFactorAuthentication(userId, authenticationType)
	c.JSON(result.StatusCode, result)
}
func (handler *UserHandler) DisableTwoFactorAuthentication(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	result := handler.IAuthenticationService.DisableTwoFactorAuthentication(userId)
	c.JSON(result.StatusCode, result)
}

func (handler *UserHandler) GenerateOTP(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	result := handler.IAuthenticationService.GenerateOTP(userId)
	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
	}
	c.JSON(result.StatusCode, result)
}

func (handler *UserHandler) VerifyOTP(c *gin.Context) {
	otp := c.Query("otp")

	validatorResult := handler.IValidator.IsNullOrEmpty(otp).Verify()
	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
	}

	result := handler.IAuthenticationService.VerifyOTP(otp)
	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
	}
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to verify OTP"))
	}
	emailResult := handler.IAuthenticationService.DecodeToken(token)

	tokenResult := handler.IAuthenticationService.GenerateAccessToken(emailResult.Data)

	c.JSON(tokenResult.StatusCode, tokenResult)
}

// disable two factor authentication
func (handler *UserHandler) DisableTwoFactorVerification(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	result := handler.IAuthenticationService.DisableTwoFactorAuthentication(userId)
	c.JSON(result.StatusCode, result)
}

// create admin
func (handler *UserHandler) CreateAdmin(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	var createUserDTO models.CreateUserDTO
	if err := c.ShouldBind(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to create user"))
	}

	result := handler.IUserService.CreateAdmin(userId, &createUserDTO)
	c.JSON(result.StatusCode, result)
}

// Change user to admin
func (handler *UserHandler) ChangeRole(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	role := c.Param("role")

	result := handler.IUserService.ChangeRole(userId, role)

	c.JSON(result.StatusCode, result)
}

// user profile - user
func (handler *UserHandler) GetUserProfile(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(jwt.MapClaims)
	userId := claims["user_id"].(int)

	result := handler.IUserService.GetUserById(userId)

	c.JSON(result.StatusCode, result)
}

// user profile - admin
func (handler *UserHandler) GetUserById(c *gin.Context) {
	id := c.Query("user_id")

	verificationResult := handler.IAuthenticationService.VerifyOTP(id)
	if !verificationResult.IsSuccess {
		c.JSON(verificationResult.StatusCode, verificationResult)
	}

	idValue, _ := strconv.Atoi(id)

	result := handler.IUserService.GetUserById(idValue)

	c.JSON(result.StatusCode, result)
}

func getLockoutCount(c *gin.Context) int {
	val, err := c.Cookie("lockout")
	if err != nil {
		return 0
	}
	count, _ := strconv.Atoi(val)
	return count
}

func getInt(key string) int {
	val, err := strconv.Atoi(key)
	if err != nil {
		return 0
	}
	return val
}
