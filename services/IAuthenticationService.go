package services

type IAuthenticationService interface {
	Authenticate(username string, password string) *AuthenticationResult
	GeneratePhoneNumberResetToken(id int) *ActionResult
	VerifyPhoneNumberResetToken(id int, token string) *ActionResult
	VerifyToken(email string, token string) *ActionResult
	GenerateResetToken(email string) *ActionResultModel[string]
	SendEmailConfirmationMail(userId int)
	DecodeToken(token string) *ActionResultModel[string]
	EnableTwoFactorAuthentication(id int, authenticationType string) *ActionResultModel[string]
	DisableTwoFactorAuthentication(id int) *ActionResult
	GenerateOTP(id int) *ActionResultModel[string]
	VerifyOTP(email string, otp string) *ActionResult
	GenerateAccessWithRefreshToken(token string) *AuthenticationResult
	GenerateAccessToken(email string) *AuthenticationResult
	DeleteRefreshToken(token string) *ActionResult
}
