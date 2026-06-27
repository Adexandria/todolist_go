package services

type IAuthenticationService interface {
	Authenticate(username string, password string) *AccessResult
	GeneratePhoneNumberResetToken(id int) *ActionResultModel[string]
	VerifyPhoneNumberResetToken(id int, token string) *ActionResult
	VerifyEmailResetToken(id int, token string) *ActionResult
	GenerateEmailResetToken(email string) *ActionResultModel[string]
	GeneratePasswordResetToken(email string) *ActionResultModel[string]
	DecodeToken(token string) *ActionResultModel[string]
	EnableTwoFactorAuthentication(id int, authenticationType string) *ActionResultModel[string]
	DisableTwoFactorAuthentication(id int) *ActionResult
	GenerateOTP(id int) *ActionResultModel[string]
	VerifyOTP(otp string) *ActionResult
	GenerateAccessToken(email string) *AccessResult
}
