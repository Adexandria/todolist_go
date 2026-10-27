package Utilities

type IValidator interface {
	IsNullOrEmpty(value string) *Validator
	CheckPasswordRule(password string) *Validator
	VerifyEmail(value string) *Validator
	VerifyPassword(value string, currentValue string) *Validator
	VerifyPhoneNumber(value string) *Validator
	Verify() *ValidatorResult
	VerifyAuthenticationType(value string) *Validator
}
