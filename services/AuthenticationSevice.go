package services

import (
	"SM/models"
	"SM/repositories"
	"SM/repositories/Utilities"
	"os"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

type AuthenticationService struct {
	IUserRepository  repositories.IUserRepository
	IPasswordManager Utilities.IPasswordManager
	ITokenManager    Utilities.ITokenManager
	SecretKey        string
}

func (a *AuthenticationService) Authenticate(username string, password string) *AccessResult {
	currentUser, err := a.IUserRepository.GetUserByUsername(username)
	if err != nil {
		return FailAccessResult("Invalid username/password")
	}

	isValidPassword := a.IPasswordManager.VerifyPassword(currentUser.Password, password)

	if !isValidPassword {
		return FailAccessResult("Invalid username/password")
	}

	if currentUser.TwoFactorEnabled {
		token, _ := a.ITokenManager.GenerateEmailToken(currentUser.Email)
		return SuccessAccessTwoFactorResult(currentUser.TwoFactorEnabled, token)
	}

	claims := make(map[string]any)
	claims["username"] = currentUser.Username
	claims["id"] = currentUser.ID
	claims["role"] = currentUser.Role

	accessToken, err := a.ITokenManager.GenerateTokenWithClaims(claims)
	if err != nil {
		return FailAccessResult("Invalid username/password")
	}

	refreshToken, err := a.ITokenManager.GenerateRefreshToken(int(currentUser.ID))
	if err != nil {
		return FailAccessResult("Invalid username/password")
	}

	currentUser.RefreshToken = refreshToken

	err = a.IUserRepository.UpdateUser(&currentUser)
	if err != nil {
		return FailAccessResult("Invalid username/password")

	}

	return SuccessAccessResult(accessToken, refreshToken)
}

func (a *AuthenticationService) GenerateAccessToken(email string) *AccessResult {
	currentUser, err := a.IUserRepository.GetUserByEmail(email)
	if err != nil {
		return FailAccessResult("Failed to authenticate user")
	}
	claims := make(map[string]any)
	claims["username"] = currentUser.Username
	claims["id"] = currentUser.ID
	claims["role"] = currentUser.Role

	accessToken, err := a.ITokenManager.GenerateTokenWithClaims(claims)
	if err != nil {
		return FailAccessResult("Invalid username/password")
	}

	refreshToken, err := a.ITokenManager.GenerateRefreshToken(int(currentUser.ID))
	if err != nil {
		return FailAccessResult("Invalid username/password")
	}

	currentUser.RefreshToken = refreshToken

	err = a.IUserRepository.UpdateUser(&currentUser)
	if err != nil {
		return FailAccessResult("Invalid username/password")

	}

	return SuccessAccessResult(accessToken, refreshToken)
}

func (a *AuthenticationService) GeneratePhoneNumberResetToken(id int) *ActionResultModel[string] {
	currentUser, err := a.IUserRepository.GetUserById(id)
	if err != nil {
		return BadRequestDataResult[string]("User not found")
	}
	token, err := a.ITokenManager.GenerateEmailToken(currentUser.Email)
	if err != nil {
		return BadRequestDataResult[string]("Failed to generate email token")
	}

	return SuccessDataResult(token)
}

func (a *AuthenticationService) VerifyPhoneNumberResetToken(id int, token string) *ActionResult {
	currentUser, err := a.IUserRepository.GetUserById(id)
	if err != nil {
		return BadRequestResult("Failed to verify phone number")
	}
	isValid := a.ITokenManager.ValidateEmailToken(currentUser.Email, token)
	if !isValid {
		return BadRequestResult("Failed to verify phone number")
	}
	return SuccessResult()
}

func (a *AuthenticationService) VerifyEmailResetToken(id int, token string) *ActionResult {
	currentUser, err := a.IUserRepository.GetUserById(id)
	if err != nil {
		return BadRequestResult("Failed to verify email")
	}
	isValid := a.ITokenManager.ValidateEmailToken(currentUser.Email, token)
	if !isValid {
		return BadRequestResult[string]("Failed to verify email")
	}
	return SuccessResult()
}

func (a *AuthenticationService) DecodeToken(token string) *ActionResultModel[string] {
	email, err := a.ITokenManager.DecodeEmailToken(token)
	if err != nil {
		return BadRequestDataResult[string]("Invalid token")
	}
	// try to use validate instead of getting the user, saves performance
	result, err := a.IUserRepository.VerifyUserByEmail(email)
	if err != nil {
		return BadRequestDataResult[string]("Invalid token")
	}
	if result {
		return SuccessDataResult(email)
	}

	return BadRequestDataResult[string]("Failed to verify token")

}

func (a *AuthenticationService) GeneratePasswordResetToken(email string) *ActionResultModel[string] {
	currentUser, err := a.IUserRepository.GetUserByEmail(email)
	if err != nil {
		return BadRequestDataResult[string]("User not found")
	}

	token, err := a.ITokenManager.GenerateEmailToken(currentUser.Email)

	if err != nil {
		return BadRequestDataResult[string]("Failed to generate password reset token")
	}

	return SuccessDataResult(token)
}

func (a *AuthenticationService) GenerateEmailResetToken(email string) *ActionResultModel[string] {
	currentUser, err := a.IUserRepository.GetUserByEmail(email)

	if err != nil {
		return BadRequestDataResult[string]("User not found")
	}

	token, err := a.ITokenManager.GenerateEmailToken(currentUser.Email)

	if err != nil {
		return BadRequestDataResult[string]("Failed to generate email token")
	}

	return SuccessDataResult(token)
}

func (a *AuthenticationService) EnableTwoFactorAuthentication(id int, authenticationType string) *ActionResultModel[string] {
	currentUser, err := a.IUserRepository.GetUserById(id)
	if err != nil {
		return BadRequestDataResult[string]("User not found")
	}

	if currentUser.TwoFactorEnabled {
		return BadRequestDataResult[string]("TwoFactor authentication is already enabled")
	}

	if authenticationType == models.AuthenticationTypes[0] {
		currentUser.AuthenticationType = 0
	} else {
		currentUser.AuthenticationType = 1
		authenticationKey := gotp.NewDefaultTOTP(a.SecretKey).ProvisioningUri(currentUser.Email, "TasksManager")

		err = qrcode.WriteFile(authenticationKey, qrcode.Medium, 256, "qr.png")
		if err != nil {
			return ServerErrorDataResult[string]("Failed to generate OTP")
		}

		currentUser.AuthenticationKey = authenticationKey

		err = a.IUserRepository.UpdateUser(&currentUser)
		if err != nil {
			return ServerErrorDataResult[string]("Failed to generate OTP")
		}
	}

	currentUser.TwoFactorEnabled = true

	err = a.IUserRepository.UpdateUser(&currentUser)
	if err != nil {
		return BadRequestDataResult[string]("Failed to enable two-factor authentication")
	}
	return SuccessDataResult(currentUser.AuthenticationKey)
}

func (a *AuthenticationService) DisableTwoFactorAuthentication(id int) *ActionResult {
	currentUser, err := a.IUserRepository.GetUserById(id)
	if err != nil {
		return BadRequestResult("User not found")
	}

	if !currentUser.TwoFactorEnabled {
		return BadRequestResult("Two-factor authentication is already disabled")
	}
	currentUser.TwoFactorEnabled = false
	currentUser.AuthenticationType = 3

	err = a.IUserRepository.UpdateUser(&currentUser)
	if err != nil {
		return BadRequestResult("Failed to enable two-factor authentication")
	}
	return SuccessResult()
}

func (a *AuthenticationService) GenerateOTP(id int) *ActionResultModel[string] {

	currentUser, err := a.IUserRepository.GetUserById(id)
	if err != nil {
		return BadRequestDataResult[string]("User not found")
	}

	if currentUser.TwoFactorEnabled == false && currentUser.AuthenticationType == 0 {
		return BadRequestDataResult[string]("Failed to generate OTP")
	}
	uri := gotp.NewDefaultTOTP(a.SecretKey)

	otp := uri.Now()

	return SuccessDataResult(otp)
}

func (a *AuthenticationService) VerifyOTP(otp string) *ActionResult {
	totp := gotp.NewDefaultTOTP(a.SecretKey)

	if totp.Verify(otp, time.Now().Unix()) {
		return SuccessResult()
	}
	return BadRequestResult("Invalid otp")
}

func AuthenticationServiceCons(userRepository *repositories.UserRepository, passwordManager *Utilities.PasswordManager,
	tokenManager *Utilities.TokenManager) *AuthenticationService {
	return &AuthenticationService{
		IUserRepository:  userRepository,
		IPasswordManager: passwordManager,
		ITokenManager:    tokenManager,
		SecretKey:        os.Getenv("TWOFACTOR_AUTHENTICATION_kEY"),
	}
}

var _ IAuthenticationService = (*AuthenticationService)(nil)
