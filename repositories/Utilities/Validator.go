package Utilities

import (
	"SM/models"
	"fmt"
	"regexp"
	"strings"
)

type Validator struct {
	AccessRule *AccessRule
	Errors     []string
}

var _ IValidator = (*Validator)(nil)

func ValidatorCon(rule *AccessRule) *Validator {
	return &Validator{
		AccessRule: rule,
		Errors:     []string{},
	}
}

func (r *Validator) VerifyPhoneNumber(value string) *Validator {
	if len(value) < 10 {
		r.Errors = append(r.Errors, "Invalid PhoneNumber")
	}
	return r
}

func (r *Validator) CheckPasswordRule(password string) *Validator {
	return r.verifyHasNumber(password).
		verifyHasCapitalLetter(password).
		verifyHasSpecialLetter(password).verifyHasSmallerLetter(password).
		verifyMinPasswordLength(password).verifyMaxPasswordLength(password)
}

func (r *Validator) VerifyPassword(value string, currentValue string) *Validator {
	if value != currentValue {
		r.Errors = append(r.Errors, fmt.Sprintf("Passwords don't match"))
	}
	return r
}

func (r *Validator) IsNullOrEmpty(value string) *Validator {
	result := value == "" || value == "null"
	if result {
		r.Errors = append(r.Errors, "Value must contain at least one character")
	}
	return r
}

func (r *Validator) VerifyEmail(value string) *Validator {
	re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
	result := re.MatchString(value)
	if !result {
		r.Errors = append(r.Errors, "Invalid email address")
	}
	return r
}

func (r *Validator) Verify() *ValidatorResult {
	if len(r.Errors) > 0 {

		currentErrors := r.Errors
		r.Errors = []string{}

		return Fail(currentErrors)
	}
	return Success()
}

func (r *Validator) verifyHasNumber(password string) *Validator {
	if r.AccessRule.PasswordRule.HasNumber {
		re := regexp.MustCompile("[0-9]+")
		result := re.MatchString(password)
		if !result {
			r.Errors = append(r.Errors, "Password must contain at least one number")
		}
	}
	return r
}
func (r *Validator) verifyHasCapitalLetter(password string) *Validator {
	if r.AccessRule.PasswordRule.HasCapitalLetter {
		re := regexp.MustCompile("[A-Z]")
		result := re.MatchString(password)
		if !result {
			r.Errors = append(r.Errors, "Password must contain at least one uppercase letter")
		}
	}
	return r
}
func (r *Validator) verifyHasSmallerLetter(password string) *Validator {
	if r.AccessRule.PasswordRule.HasSpecialCharacter {
		re := regexp.MustCompile("[a-z]")
		result := re.MatchString(password)
		if !result {
			r.Errors = append(r.Errors, "Password must contain at least one lowercase letter")
		}
	}
	return r
}
func (r *Validator) verifyHasSpecialLetter(password string) *Validator {
	if r.AccessRule.PasswordRule.HasSpecialCharacter {
		re := regexp.MustCompile("[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>/?]")
		result := re.MatchString(password)
		if !result {
			r.Errors = append(r.Errors, "Password must contain at least one special character letter")
		}
	}
	return r
}

func (r *Validator) verifyMinPasswordLength(password string) *Validator {
	if r.AccessRule.PasswordRule.MinimumPasswordLength > len(password) {
		err := fmt.Sprintf("Invalid Password Length, Password must contain more than %d characters", r.AccessRule.PasswordRule.MinimumPasswordLength)
		r.Errors = append(r.Errors, err)
	}
	return r
}

func (r *Validator) verifyMaxPasswordLength(password string) *Validator {
	if r.AccessRule.PasswordRule.MaximumPasswordLength > len(password) {
		err := fmt.Sprintf("Invalid Password Length, Password must not contain more than %d characters", r.AccessRule.PasswordRule.MaximumPasswordLength)
		r.Errors = append(r.Errors, err)
	}
	return r
}

func (r *Validator) VerifyAuthenticationType(value string) *Validator {
	for _, name := range models.AuthenticationTypes {
		if name == strings.ToLower(value) {
			return r
		}
	}
	err := fmt.Sprintf("Invalid Authentication Type: %s Only sms and google are accepted", value)
	r.Errors = append(r.Errors, err)
	return r
}
