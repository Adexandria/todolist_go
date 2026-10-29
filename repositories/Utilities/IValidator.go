package Utilities

type Validator interface {
	IsNullOrEmpty(value string) *validator
	CheckPasswordRule(password string) *validator
	VerifyEmail(value string) *validator
	VerifyPassword(value string, currentValue string) *validator
	VerifyPhoneNumber(value string) *validator
	Verify() *ValidatorResult
	VerifyAuthenticationType(value string) *validator
}
