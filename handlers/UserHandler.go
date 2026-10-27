package handlers

import (
	"SM/models"
	"SM/repositories/Utilities"
	"SM/services"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	IUserService           services.IUserService
	IAuthenticationService services.IAuthenticationService
	LockoutTrial           string
	IValidator             Utilities.IValidator
	IEmailService          services.IEmailService
}

func UserHandlerCon(userService *services.UserService, authenticationService *services.AuthenticationService,
	validator *Utilities.Validator, emailService *services.EmailService) *UserHandler {
	return &UserHandler{IUserService: userService,
		IAuthenticationService: authenticationService,
		LockoutTrial:           os.Getenv("LOCKOUT_TRIAL"),
		IValidator:             validator,
		IEmailService:          emailService}
}

// --- Authentication & Session Management ---

// SignUp handles new user registration
func (handler *UserHandler) SignUp(c *gin.Context) {
	var createUserDTO models.CreateUserDTO
	if err := c.ShouldBind(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to create user"))
		return
	}

	validatorResult := handler.IValidator.IsNullOrEmpty(createUserDTO.Username).
		IsNullOrEmpty(createUserDTO.LastName).IsNullOrEmpty(createUserDTO.FirstName).
		VerifyPassword(createUserDTO.Password, createUserDTO.RetypePassword).
		CheckPasswordRule(createUserDTO.Password).VerifyEmail(createUserDTO.Email).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	result := handler.IUserService.CreateUser(&createUserDTO)
	if !result.IsSuccess {
		c.JSON(http.StatusBadRequest, result)
		return
	}

	emailTokenResult := handler.IAuthenticationService.GenerateResetToken(createUserDTO.Email)
	if !emailTokenResult.IsSuccess {
		c.JSON(http.StatusBadRequest, emailTokenResult)
		return
	}
	data := map[string]any{
		"Verification_url": fmt.Sprintf("%s/api/auth/confirm-email?token=%s&email=%s", os.Getenv("BASE_URL"), emailTokenResult.Data, createUserDTO.Email),
		"Username":         createUserDTO.Username,
	}

	emailTemplate, err := handler.IEmailService.GenerateHTMLTemplate("Welcome", data)
	if err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to create user"))
		return
	}

	mailMessage := handler.IEmailService.GenerateMailgunTemplate(os.Getenv("FROM_MAIL"), createUserDTO.Email, "Welcome to Task Manager", emailTemplate)

	handler.IEmailService.SendEmailAsync(mailMessage)

	c.JSON(result.StatusCode, result)
}

// Login handles user authentication and session creation
func (handler *UserHandler) Login(c *gin.Context) {
	var loginDTO models.LoginDTO

	if err := c.ShouldBind(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Invalid username/password"))
		return
	}

	validatorResult := handler.IValidator.IsNullOrEmpty(loginDTO.Username).IsNullOrEmpty(loginDTO.Password).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	value := getLockoutCount(c, loginDTO.Username)
	if value == getInt(handler.LockoutTrial) {
		c.SetCookie(fmt.Sprintf("%slockout_", loginDTO.Username), "0", -1, "/", "", true, true)
		// when you create the ticket handler

		_ = handler.IUserService.EnableLockout(loginDTO.Username)

		c.JSON(http.StatusBadRequest, services.BadRequestResult("This user has been locked out"))
		return
	}

	result := handler.IAuthenticationService.Authenticate(loginDTO.Username, loginDTO.Password)
	if !result.IsSuccess && result.Lockout == 1 {

		if !result.IsLockedOut {
			retries := value + 1

			c.SetCookie(fmt.Sprintf("%slockout_", loginDTO.Username), strconv.Itoa(retries), 60*60, "/", "", true, true)

			retries = getInt(handler.LockoutTrial) - retries

			result.Error = append(result.Error, fmt.Sprintf("%d left", retries))
		}
		c.JSON(result.StatusCode, result.AccessResult)

		return

	}

	if result.IsTwoFactorEnabled {
		c.SetCookie("2fa-token", result.TwoFactorToken, 4*60, "/", "", true, true)
		c.JSON(http.StatusAccepted, result.AccessResult)
		return
	}

	c.SetCookie("refresh-token", result.RefreshToken, 3*60*60, "/", "", true, true)
	c.SetCookie(fmt.Sprintf("%slockout:", loginDTO.Username), "0", -1, "/", "", true, true)
	c.JSON(result.StatusCode, result.AccessResult)
}

// Logout terminates the user session
func (handler *UserHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh-token")
	if err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to log out"))
	}
	result := handler.IAuthenticationService.DeleteRefreshToken(refreshToken)

	c.JSON(result.StatusCode, result)
}

// GenerateAccessToken creates a new access token using a refresh token
func (handler *UserHandler) GenerateAccessToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh-token")
	if err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Invalid refresh token"))
		return
	}

	result := handler.IAuthenticationService.GenerateAccessWithRefreshToken(refreshToken)

	if result.IsSuccess {
		c.SetCookie("refresh-token", result.RefreshToken, 3*60*60, "/", "", true, true)

	}

	c.JSON(result.StatusCode, result.AccessResult)
}

// RequestOTP sends a One-Time Password to the user for 2FA
func (handler *UserHandler) RequestOTP(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	result := handler.IAuthenticationService.GenerateOTP(userId)
	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
		return
	}
	c.JSON(result.StatusCode, result)
}

// VerifyOTP validates the One-Time Password provided by the user
func (handler *UserHandler) VerifyOTP(c *gin.Context) {
	otp := c.Query("otp")

	validatorResult := handler.IValidator.IsNullOrEmpty(otp).Verify()
	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	token, err := c.Cookie("2fa-token")

	if err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to verify OTP"))
		return
	}

	decodedResult := handler.IAuthenticationService.DecodeToken(token)

	result := handler.IAuthenticationService.VerifyOTP(decodedResult.Data, otp)
	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
		return
	}

	tokenResult := handler.IAuthenticationService.GenerateAccessToken(decodedResult.Data)

	c.JSON(tokenResult.StatusCode, tokenResult.AccessResult)
}

// SetUpTwoFactorAuthentication enables 2FA for the user
func (handler *UserHandler) SetUpTwoFactorAuthentication(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	authenticationType := c.Query("authentication_type")

	validatorResult := handler.IValidator.IsNullOrEmpty(authenticationType).VerifyAuthenticationType(authenticationType).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	result := handler.IAuthenticationService.EnableTwoFactorAuthentication(userId, authenticationType)
	c.JSON(result.StatusCode, result)
}

// DisableTwoFactorAuthentication disables 2FA for the user
func (handler *UserHandler) DisableTwoFactorAuthentication(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	result := handler.IAuthenticationService.DisableTwoFactorAuthentication(userId)
	c.JSON(result.StatusCode, result)
}

// --- Account Management & Security ---

// ChangePassword updates the password for a logged-in user
func (handler *UserHandler) ChangePassword(c *gin.Context) {
	var changePasswordDTO models.ChangePasswordDTO
	if err := c.ShouldBind(&changePasswordDTO); err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to change password"))
		return
	}

	validatorResult := handler.IValidator.CheckPasswordRule(changePasswordDTO.NewPassword).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)

	result := handler.IUserService.ChangePassword(claims.Id, changePasswordDTO.OldPassword, changePasswordDTO.NewPassword)
	c.JSON(result.StatusCode, result)
}

// ResetPassword initiates the password reset process by sending an email
func (handler *UserHandler) ResetPassword(c *gin.Context) {
	email := c.Query("email")

	validatorResult := handler.IValidator.IsNullOrEmpty(email).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	result := handler.IAuthenticationService.GenerateResetToken(email)

	data := map[string]any{
		"Company_Name":    "Todo List API",
		"Expiration_Time": time.Now().Add(30 * 60 * time.Second).String(),
		"Reset_Link":      fmt.Sprintf("%s/api//auth/reset-password?token=%s&email=%s", os.Getenv("BASE_URL"), result.Data, email),
		"Support_Link":    "mailto:adeolaaderibigbe09@gmail.com",
	}

	emailTemplate, err := handler.IEmailService.GenerateHTMLTemplate("Forgot", data)
	if err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to generate token"))
		return
	}

	mailMessage := handler.IEmailService.GenerateMailgunTemplate(os.Getenv("FROM_MAIL"), email, "Password Reset", emailTemplate)

	handler.IEmailService.SendEmailAsync(mailMessage)

	c.JSON(result.StatusCode, result)
}

// ResetPasswordByToken completes the password reset process using a token
func (handler *UserHandler) ResetPasswordByToken(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")

	var password models.AnonymousChangePasswordDTO
	if err := c.ShouldBind(&password); err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to change password"))
		return
	}
	validatorResult := handler.IValidator.IsNullOrEmpty(password.NewPassword).
		IsNullOrEmpty(token).IsNullOrEmpty(email).CheckPasswordRule(password.NewPassword).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	result := handler.IAuthenticationService.VerifyToken(email, token)
	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
		return
	}

	passwordChangedResult := handler.IUserService.ChangeAnonPassword(email, password.NewPassword)

	c.JSON(result.StatusCode, passwordChangedResult)
}

// ChangeEmail updates the user's email address and sends a confirmation mail
func (handler *UserHandler) ChangeEmail(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	email := c.Query("email")

	validatorResult := handler.IValidator.IsNullOrEmpty(email).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	result := handler.IUserService.SetEmail(userId, email)

	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
		return
	}

	handler.IAuthenticationService.SendEmailConfirmationMail(userId)

	c.JSON(result.StatusCode, result)
}

// ConfirmEmailToken verifies the email confirmation token
func (handler *UserHandler) ConfirmEmailToken(c *gin.Context) {
	email := c.Query("email")
	token := c.Query("token")

	validatorResult := handler.IValidator.IsNullOrEmpty(token).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	result := handler.IAuthenticationService.VerifyToken(email, token)

	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
		return
	}
	confirmEmailResult := handler.IUserService.ConfirmEmail(email)
	c.JSON(result.StatusCode, confirmEmailResult)
}

// ChangePhonenumber updates the user's phone number and generates a verification token
func (handler *UserHandler) ChangePhonenumber(c *gin.Context) {
	phonenumber := c.Query("phonenumber")

	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	validatorResult := handler.IValidator.VerifyPhoneNumber(phonenumber).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	result := handler.IUserService.SetPhoneNumber(userId, phonenumber)

	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
		return
	}

	tokenResult := handler.IAuthenticationService.GeneratePhoneNumberResetToken(claims.Id)

	c.JSON(result.StatusCode, tokenResult)
}

// VerifyPhonenumberToken verifies the phone number using a token
func (handler *UserHandler) VerifyPhonenumberToken(c *gin.Context) {
	token := c.Query("token")
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	validatorResult := handler.IValidator.IsNullOrEmpty(token).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	result := handler.IAuthenticationService.VerifyPhoneNumberResetToken(userId, token)
	if !result.IsSuccess {
		c.JSON(result.StatusCode, result)
		return
	}
	confirmPhonenumber := handler.IUserService.ConfirmPhoneNumber(userId)
	c.JSON(result.StatusCode, confirmPhonenumber)
}

// --- User Profile & Information ---

// GetUserProfile retrieves the profile information of the currently logged-in user
func (handler *UserHandler) GetUserProfile(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	result := handler.IUserService.GetUserById(userId)

	c.JSON(result.StatusCode, result)
}

// UpdateFirstName updates the first name of the user
func (handler *UserHandler) UpdateFirstName(c *gin.Context) {
	firstName := c.Query("firstname")
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	validatorResult := handler.IValidator.IsNullOrEmpty(firstName).Verify()
	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}
	result := handler.IUserService.SetFirstName(userId, firstName)
	c.JSON(result.StatusCode, result)
}

// UpdateLastName updates the last name of the user
func (handler *UserHandler) UpdateLastName(c *gin.Context) {
	lastName := c.Query("lastname")
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	validatorResult := handler.IValidator.IsNullOrEmpty(lastName).Verify()
	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	result := handler.IUserService.SetLastName(userId, lastName)
	c.JSON(result.StatusCode, result)
}

// UpdateUsername updates the username of the user
func (handler *UserHandler) UpdateUsername(c *gin.Context) {
	username := c.Query("username")
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	validatorResult := handler.IValidator.IsNullOrEmpty(username).Verify()
	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	result := handler.IUserService.SetUsername(userId, username)
	c.JSON(result.StatusCode, result)
}

// --- Administrative Actions ---

// GetUserById retrieves user profile by ID (Admin only)
func (handler *UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("id")

	verificationResult := handler.IValidator.IsNullOrEmpty(id).Verify()
	if !verificationResult.IsSuccess {
		c.JSON(http.StatusBadRequest, verificationResult)
		return
	}

	idValue, _ := strconv.Atoi(id)

	result := handler.IUserService.GetUserById(idValue)

	c.JSON(result.StatusCode, result)
}

// CreateAdmin creates a new administrative user
func (handler *UserHandler) CreateAdmin(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	var createUserDTO models.CreateUserDTO
	if err := c.ShouldBind(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, services.BadRequestResult("Failed to create user"))
		return
	}

	result := handler.IUserService.CreateAdmin(userId, &createUserDTO)
	c.JSON(result.StatusCode, result)
}

// ChangeRole changes the role of a user
func (handler *UserHandler) ChangeRole(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id

	//Add id as a query

	role := c.Query("role")

	result := handler.IUserService.ChangeRole(userId, role)

	c.JSON(result.StatusCode, result)
}

// DisableUserLockout removes the lockout from a user account (Admin only)
func (handler *UserHandler) DisableUserLockout(c *gin.Context) {
	userId := c.Param("id")

	validatorResult := handler.IValidator.IsNullOrEmpty(userId).Verify()

	if !validatorResult.IsSuccess {
		c.JSON(http.StatusBadRequest, validatorResult)
		return
	}

	id, _ := strconv.Atoi(userId)

	result := handler.IUserService.DisableLockout(id)
	c.JSON(result.StatusCode, result)
}

// DisableUserTwoFactorAuthentication disables 2FA for a user using admin privileges
func (handler *UserHandler) DisableUserTwoFactorAuthentication(c *gin.Context) {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)

	result := handler.IAuthenticationService.DisableTwoFactorAuthentication(userId)
	c.JSON(result.StatusCode, result)
}

// DeleteUserById deletes a user account by ID (Admin only)
func (handler *UserHandler) DeleteUserById(c *gin.Context) {
	userId := c.Param("id")
	id, _ := strconv.Atoi(userId)

	result := handler.IUserService.DeleteUser(id)
	c.JSON(result.StatusCode, result)
}

// --- Account Lifecycle ---

// DeleteUserAccount allows a user to delete their own account
func (handler *UserHandler) DeleteUserAccount(c *gin.Context) {
	claimsString, _ := c.Get("claims")
	claims := claimsString.(*Utilities.IdentityClaim)
	userId := claims.Id
	result := handler.IUserService.DeleteUser(userId)

	c.JSON(result.StatusCode, result)
}

// --- Utility Functions ---

func getLockoutCount(c *gin.Context, username string) int {
	val, err := c.Cookie(fmt.Sprintf("%slockout_", username))
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
