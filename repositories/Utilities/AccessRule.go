package Utilities

type AccessRule struct {
	IsRequireEmailConfirmation bool
	PasswordRule               PasswordRule
}
type PasswordRule struct {
	HasNumber             bool
	HasCapitalLetter      bool
	HasSmallLetter        bool
	HasSpecialCharacter   bool
	MinimumPasswordLength int
	MaximumPasswordLength int
}

func (r PasswordRule) EnableNumber() PasswordRule {
	r.HasNumber = true
	return r
}

func (r PasswordRule) EnableCapitalLetter() PasswordRule {
	r.HasCapitalLetter = true
	return r
}

func (r PasswordRule) EnableSmallLetter() PasswordRule {
	r.HasSmallLetter = true
	return r
}

func (r PasswordRule) EnableSpecialCharacter() PasswordRule {
	r.HasSpecialCharacter = true
	return r
}

func (r PasswordRule) SetMinimumPasswordLength(length ...int) PasswordRule {
	minLen := 6
	if len(length) > 0 {
		minLen = length[0]
	}
	r.MinimumPasswordLength = minLen
	return r
}

func (r PasswordRule) SetMaximumPasswordLength(length ...int) PasswordRule {
	maxLen := 3

	if len(length) > 0 {
		maxLen = length[0]
	}

	r.MinimumPasswordLength = maxLen
	return r
}

func (e AccessRule) EnableEmailConfirmation() AccessRule {
	e.IsRequireEmailConfirmation = true
	return e
}

func (e AccessRule) EnablePasswordValidation(p PasswordRule) AccessRule {
	e.PasswordRule = p
	return e
}
