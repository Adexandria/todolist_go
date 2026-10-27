package services

import (
	"SM/models"
	"SM/repositories"
	"SM/repositories/Utilities"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

type AuthenticationService struct {
	IUserRepository  repositories.IUserRepository
	IPasswordManager Utilities.IPasswordManager
	ITokenManager    Utilities.ITokenManager
	AccessRule       *Utilities.AccessRule
	IEmailService    IEmailService
	ISmsService      ISmsService
	Log              *slog.Logger
}

func (a *AuthenticationService) Authenticate(username string, password string) *AuthenticationResult {
	currentUser, err := a.IUserRepository.GetUserByUsername(username)
	if err != nil {
		return FailAuthenticationResult("Invalid username/password")
	}

	if currentUser.LockoutEnabled {
		return SuccessLockoutAuthenticationResult(1)
	}

	isValidPassword := a.IPasswordManager.VerifyPassword(currentUser.Password, password)

	if !isValidPassword {
		return FailLockoutAuthenticationResult(1, "Invalid username/password")
	}

	if a.AccessRule.IsRequireEmailConfirmation && !currentUser.EmailConfirmed {
		emailToken := a.GenerateResetToken(currentUser.Email)
		sendEmailConfirmationMail(a, currentUser.Username, currentUser.Email, emailToken.Data)
		return FailAuthenticationResult("Require email confirmation")
	}

	if currentUser.TwoFactorEnabled {
		_2faToken, _ := a.ITokenManager.GenerateToken(currentUser.Email)
		return SuccessAuthenticateTwoFactorResult(currentUser.TwoFactorEnabled, _2faToken)
	}

	claims := make(map[string]any)
	claims["username"] = currentUser.Username
	claims["id"] = currentUser.ID
	claims["role"] = currentUser.Role

	accessToken, err := a.ITokenManager.GenerateJWTTokenWithClaims(claims)
	if err != nil {
		return FailAuthenticationResult("Invalid username/password")
	}

	refreshToken, err := a.ITokenManager.GenerateRefreshToken(int(currentUser.ID))
	if err != nil {
		return FailAuthenticationResult("Invalid username/password")
	}

	currentUser.RefreshToken = refreshToken

	err = a.IUserRepository.UpdateUser(&currentUser)
	if err != nil {
		return FailAuthenticationResult("Invalid username/password")

	}

	sendLoginMail(a, currentUser.Username, currentUser.Email)

	return SuccessAuthenticationResult(accessToken, refreshToken)
}

func sendLoginMail(a *AuthenticationService, username string, email string) {
	data := map[string]any{
		"Username":    username,
		"Login_time":  time.Now().UTC().Format("2006-01-02 15:04:05"),
		"Secure_link": path.Join(os.Getenv("BASE_URL"), ""),
	}

	emailTemplate, err := a.IEmailService.GenerateHTMLTemplate("Login", data)
	if err != nil {
		a.Log.Error(err.Error())
	}

	mailMessage := a.IEmailService.GenerateMailgunTemplate(os.Getenv("FROM_MAIL"), email, "Login Notification", emailTemplate)

	a.IEmailService.SendEmailAsync(mailMessage)
}

func sendEmailConfirmationMail(a *AuthenticationService, username string, email string, token string) {
	data := map[string]any{
		"Verification_url": fmt.Sprintf("%s/api/auth/confirm-email?token=%s&email=%s", os.Getenv("BASE_URL"), token, email),
		"Username":         username,
	}

	emailTemplate, err := a.IEmailService.GenerateHTMLTemplate("Welcome", data)
	if err != nil {
		a.Log.Error(err.Error())
	}

	mailMessage := a.IEmailService.GenerateMailgunTemplate(os.Getenv("FROM_MAIL"), email, "Confirm Email", emailTemplate)

	a.IEmailService.SendEmailAsync(mailMessage)
}

func (a *AuthenticationService) SendEmailConfirmationMail(userId int) {
	currentUser, err := a.IUserRepository.GetUserById(userId)
	if err != nil {
		a.Log.Error(err.Error())
	}

	tokenResult := a.GenerateResetToken(currentUser.Email)
	if tokenResult != nil {
		a.Log.Error(err.Error())
	}

	sendEmailConfirmationMail(a, currentUser.Username, currentUser.Email, tokenResult.Data)
}

func (a *AuthenticationService) GenerateAccessWithRefreshToken(token string) *AuthenticationResult {
	currentUser, err := a.IUserRepository.GetUserByRefreshToken(token)

	verifiedToken := a.ITokenManager.VerifyRefreshToken(int(currentUser.ID), token)

	if verifiedToken {
		return FailAuthenticationResult("Invalid refresh token")
	}

	if err != nil {
		return FailAuthenticationResult("Invalid refresh token")
	}

	claims := make(map[string]any)
	claims["username"] = currentUser.Username
	claims["id"] = currentUser.ID
	claims["role"] = currentUser.Role

	accessToken, err := a.ITokenManager.GenerateJWTTokenWithClaims(claims)
	if err != nil {
		return FailAuthenticationResult("Invalid refresh token")
	}

	refreshToken, err := a.ITokenManager.GenerateRefreshToken(int(currentUser.ID))
	if err != nil {
		return FailAuthenticationResult("Invalid refresh token")
	}

	currentUser.RefreshToken = refreshToken

	err = a.IUserRepository.UpdateUser(&currentUser)
	if err != nil {
		return FailAuthenticationResult("Invalid refresh token")

	}

	return SuccessAuthenticationResult(accessToken, refreshToken)
}

func (a *AuthenticationService) GenerateAccessToken(email string) *AuthenticationResult {
	currentUser, err := a.IUserRepository.GetUserByEmail(email)

	if err != nil {
		return FailAuthenticationResult("Invalid OTP")
	}

	claims := make(map[string]any)
	claims["username"] = currentUser.Username
	claims["id"] = currentUser.ID
	claims["role"] = currentUser.Role

	accessToken, err := a.ITokenManager.GenerateJWTTokenWithClaims(claims)
	if err != nil {
		return FailAuthenticationResult("Invalid OTP")
	}

	refreshToken, err := a.ITokenManager.GenerateRefreshToken(int(currentUser.ID))
	if err != nil {
		return FailAuthenticationResult("Invalid OTP")
	}

	currentUser.RefreshToken = refreshToken

	err = a.IUserRepository.UpdateUser(&currentUser)
	if err != nil {
		return FailAuthenticationResult("Invalid OTP")

	}

	return SuccessAuthenticationResult(accessToken, refreshToken)
}

func (a *AuthenticationService) DeleteRefreshToken(token string) *ActionResult {
	currentUser, err := a.IUserRepository.GetUserByRefreshToken(token)

	if err != nil {
		return BadRequestResult("Failed to log out")
	}

	verifiedToken := a.ITokenManager.VerifyRefreshToken(int(currentUser.ID), token)

	if verifiedToken {
		return BadRequestResult("Failed to log out")
	}

	currentUser.RefreshToken = ""
	err = a.IUserRepository.UpdateUser(&currentUser)
	if err != nil {
		return BadRequestResult("Failed to log out")
	}

	return SuccessResult()

}

func (a *AuthenticationService) GeneratePhoneNumberResetToken(id int) *ActionResult {
	currentUser, err := a.IUserRepository.GetUserById(id)

	if err != nil {
		return NotFoundResult("User not found")
	}

	token, err := a.ITokenManager.GenerateToken(currentUser.Email)

	if err != nil {
		return BadRequestResult("Failed to generate email token")
	}

	data := map[string]any{
		"Verification_url": fmt.Sprintf("%s/api/account/verify-phonenumber?token=%s", os.Getenv("BASE_URL"), token),
		"Username":         currentUser.Username,
	}

	emailTemplate, err := a.IEmailService.GenerateHTMLTemplate("Phonenumber", data)
	if err != nil {
		return BadRequestResult("Failed to create user")
	}

	mailMessage := a.IEmailService.GenerateMailgunTemplate(os.Getenv("FROM_MAIL"), currentUser.Email, "Confirm Phonenumber", emailTemplate)

	a.IEmailService.SendEmailAsync(mailMessage)

	return SuccessResult()
}

func (a *AuthenticationService) VerifyPhoneNumberResetToken(id int, token string) *ActionResult {
	currentUser, err := a.IUserRepository.GetUserById(id)
	if err != nil {
		return BadRequestResult("Failed to verify phone number")
	}
	isValid := a.ITokenManager.ValidateToken(currentUser.Email, token)
	if !isValid {
		return BadRequestResult("Failed to verify phone number")
	}
	return SuccessResult()
}

func (a *AuthenticationService) VerifyToken(email string, token string) *ActionResult {
	_, err := a.IUserRepository.VerifyUserByEmail(email)

	if err != nil {
		return BadRequestResult("Failed to verify email")
	}

	isValid := a.ITokenManager.ValidateToken(email, token)

	if !isValid {
		return BadRequestResult("Failed to verify email")
	}

	return SuccessResult()
}

func (a *AuthenticationService) DecodeToken(token string) *ActionResultModel[string] {
	email, err := a.ITokenManager.DecodeToken(token)
	if err != nil {
		return BadRequestDataResult[string]("Invalid token")
	}

	result, err := a.IUserRepository.VerifyUserByEmail(email)

	if err != nil {
		return BadRequestDataResult[string]("Invalid token")
	}

	if result {
		return SuccessDataResult(email)
	}

	return BadRequestDataResult[string]("Failed to verify token")

}

func (a *AuthenticationService) GenerateResetToken(email string) *ActionResultModel[string] {
	currentUser, err := a.IUserRepository.GetUserByEmail(email)

	if err != nil {
		return NotFoundDataResult[string]("User not found")
	}

	token, err := a.ITokenManager.GenerateToken(currentUser.Email)

	if err != nil {
		return BadRequestDataResult[string]("Failed to generate email token")
	}

	return SuccessDataResult(token)
}

func (a *AuthenticationService) EnableTwoFactorAuthentication(id int, authenticationType string) *ActionResultModel[string] {
	currentUser, err := a.IUserRepository.GetUserById(id)
	if err != nil {
		return NotFoundDataResult[string]("User not found")
	}

	if currentUser.TwoFactorEnabled {
		return BadRequestDataResult[string]("TwoFactor authentication is already enabled")
	}
	data := ""

	if authenticationType == models.AuthenticationTypes[0] {

		if currentUser.PhoneNumber == "" || !currentUser.PhoneNumberConfirmed {
			return BadRequestDataResult[string]("Phonenumber has to be set or confirmed")
		}

		currentUser.AuthenticationType = 0
	} else {
		currentUser.AuthenticationType = 1
		secret := gotp.RandomSecret(32)
		url := gotp.NewDefaultTOTP(secret).ProvisioningUri(currentUser.Email, "TasksManager")

		png, err := qrcode.Encode(url, qrcode.Medium, 256)
		if err != nil {
			return ServerErrorDataResult[string]("Failed to generate OTP")
		}

		currentUser.AuthenticationKey = secret

		err = a.IUserRepository.UpdateUser(&currentUser)
		if err != nil {
			return ServerErrorDataResult[string]("Failed to generate OTP")
		}

		data = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	}

	currentUser.TwoFactorEnabled = true

	err = a.IUserRepository.UpdateUser(&currentUser)

	if err != nil {
		return BadRequestDataResult[string]("Failed to enable two-factor authentication")
	}

	return SuccessDataResult(data)
}

func (a *AuthenticationService) DisableTwoFactorAuthentication(id int) *ActionResult {
	currentUser, err := a.IUserRepository.GetUserById(id)
	if err != nil {
		return NotFoundResult("User not found")
	}

	if !currentUser.TwoFactorEnabled {
		return BadRequestResult("Two-factor authentication is already disabled")
	}
	currentUser.TwoFactorEnabled = false
	currentUser.AuthenticationType = 2
	currentUser.AuthenticationKey = ""

	err = a.IUserRepository.UpdateUser(&currentUser)
	if err != nil {
		return BadRequestResult("Failed to enable two-factor authentication")
	}
	return SuccessResult()
}

func (a *AuthenticationService) GenerateOTP(id int) *ActionResultModel[string] {
	currentUser, err := a.IUserRepository.GetUserById(id)

	if err != nil {
		return NotFoundDataResult[string]("User not found")
	}

	if currentUser.TwoFactorEnabled == false && currentUser.AuthenticationType != 0 {
		return BadRequestDataResult[string]("Failed to generate OTP")
	}
	totp := gotp.NewTOTP(currentUser.AuthenticationKey, 6, 180, nil)

	otp := totp.Now()
	// store the otp so it can't be reuse after a new one has been requested.

	a.ISmsService.SendOTP(currentUser.PhoneNumber, otp)

	return SuccessDataResult(otp)
}

func (a *AuthenticationService) VerifyOTP(email string, otp string) *ActionResult {

	currentUser, err := a.IUserRepository.GetUserByEmail(email)
	if err != nil {
		return NotFoundResult("User not found")
	}

	totp := gotp.NewDefaultTOTP(currentUser.AuthenticationKey)

	now := time.Now().Unix()
	for _, offset := range []int64{-30, 0, 30} {
		if totp.Verify(otp, now+offset) {
			return SuccessResult()
		}
	}
	return BadRequestResult("Invalid otp")
}

func AuthenticationServiceCons(userRepository *repositories.UserRepository, passwordManager *Utilities.PasswordManager,
	tokenManager *Utilities.TokenManager, smsService *SmsService, accessrule *Utilities.AccessRule, emailService *EmailService, handler *slog.JSONHandler) *AuthenticationService {
	return &AuthenticationService{
		IUserRepository:  userRepository,
		IPasswordManager: passwordManager,
		ITokenManager:    tokenManager,
		AccessRule:       accessrule,
		IEmailService:    emailService,
		ISmsService:      smsService,
		Log:              slog.New(handler),
	}
}

var _ IAuthenticationService = (*AuthenticationService)(nil)
