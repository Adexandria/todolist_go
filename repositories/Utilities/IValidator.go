package Utilities

type IValidator interface {
	IsNullOrEmpty(value string) *Validator
	CheckPasswordRule(password string) *Validator
	VerifyEmail(value string) *Validator
	VerifyEmailRequirement(isEmailConfirmed bool) *Validator
	VerifyPassword(value string, currentValue string) *Validator
	VerifyPhoneNumber(value string) *Validator
	Verify() *ValidatorResult
}
