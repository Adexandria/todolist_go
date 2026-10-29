package Utilities

import (
	"SM/models"
	"fmt"
	"regexp"
	"strings"
)

type validator struct {
	AccessRule *AccessRule
	Errors     []string
}

var _ Validator = (*validator)(nil)

func ValidatorCon(rule *AccessRule) Validator {
	return &validator{
		AccessRule: rule,
		Errors:     []string{},
	}
}

func (r *validator) VerifyPhoneNumber(value string) *validator {
	if len(value) < 10 {
		r.Errors = append(r.Errors, "Invalid PhoneNumber")
	}
	return r
}

func (r *validator) CheckPasswordRule(password string) *validator {
	return r.verifyHasNumber(password).
		verifyHasCapitalLetter(password).
		verifyHasSpecialLetter(password).verifyHasSmallerLetter(password).
		verifyMinPasswordLength(password).verifyMaxPasswordLength(password)
}

func (r *validator) VerifyPassword(value string, currentValue string) *validator {
	if value != currentValue {
		r.Errors = append(r.Errors, "Passwords don't match")
	}
	return r
}

func (r *validator) IsNullOrEmpty(value string) *validator {
	result := value == "" || value == "null"
	if result {
		r.Errors = append(r.Errors, "Value must contain at least one character")
	}
	return r
}

func (r *validator) VerifyEmail(value string) *validator {
	re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
	result := re.MatchString(value)
	if !result {
		r.Errors = append(r.Errors, "Invalid email address")
	}
	return r
}

func (r *validator) Verify() *ValidatorResult {
	if len(r.Errors) > 0 {

		currentErrors := r.Errors
		r.Errors = []string{}

		return Fail(currentErrors)
	}
	return Success()
}

func (r *validator) verifyHasNumber(password string) *validator {
	if r.AccessRule.PasswordRule.HasNumber {
		re := regexp.MustCompile("[0-9]+")
		result := re.MatchString(password)
		if !result {
			r.Errors = append(r.Errors, "Password must contain at least one number")
		}
	}
	return r
}
func (r *validator) verifyHasCapitalLetter(password string) *validator {
	if r.AccessRule.PasswordRule.HasCapitalLetter {
		re := regexp.MustCompile("[A-Z]")
		result := re.MatchString(password)
		if !result {
			r.Errors = append(r.Errors, "Password must contain at least one uppercase letter")
		}
	}
	return r
}
func (r *validator) verifyHasSmallerLetter(password string) *validator {
	if r.AccessRule.PasswordRule.HasSpecialCharacter {
		re := regexp.MustCompile("[a-z]")
		result := re.MatchString(password)
		if !result {
			r.Errors = append(r.Errors, "Password must contain at least one lowercase letter")
		}
	}
	return r
}
func (r *validator) verifyHasSpecialLetter(password string) *validator {
	if r.AccessRule.PasswordRule.HasSpecialCharacter {
		re := regexp.MustCompile("[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>/?]")
		result := re.MatchString(password)
		if !result {
			r.Errors = append(r.Errors, "Password must contain at least one special character letter")
		}
	}
	return r
}

func (r *validator) verifyMinPasswordLength(password string) *validator {
	if r.AccessRule.PasswordRule.MinimumPasswordLength > len(password) {
		err := fmt.Sprintf("Invalid Password Length, Password must contain more than %d characters", r.AccessRule.PasswordRule.MinimumPasswordLength)
		r.Errors = append(r.Errors, err)
	}
	return r
}

func (r *validator) verifyMaxPasswordLength(password string) *validator {
	if r.AccessRule.PasswordRule.MaximumPasswordLength > len(password) {
		err := fmt.Sprintf("Invalid Password Length, Password must not contain more than %d characters", r.AccessRule.PasswordRule.MaximumPasswordLength)
		r.Errors = append(r.Errors, err)
	}
	return r
}

func (r *validator) VerifyAuthenticationType(value string) *validator {
	for _, name := range models.AuthenticationTypes {
		if name == strings.ToLower(value) {
			return r
		}
	}
	err := fmt.Sprintf("Invalid Authentication Type: %s Only sms and google are accepted", value)
	r.Errors = append(r.Errors, err)
	return r
}
