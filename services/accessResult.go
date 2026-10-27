package services

import "net/http"

type AccessResult struct {
	StatusCode         int
	AccessToken        string
	IsSuccess          bool
	Error              []string
	IsTwoFactorEnabled bool
	IsLockedOut        bool
}

type AuthenticationResult struct {
	Lockout        int
	RefreshToken   string
	TwoFactorToken string
	AccessResult
}

func SuccessAccessResult(accessToken string) *AccessResult {
	return &AccessResult{
		StatusCode:  http.StatusOK,
		AccessToken: accessToken,
		IsSuccess:   true,
	}
}
func SuccessAuthenticationResult(accessToken string, refreshToken string) *AuthenticationResult {
	return &AuthenticationResult{
		RefreshToken: refreshToken,
		AccessResult: AccessResult{
			StatusCode:  http.StatusOK,
			AccessToken: accessToken,
			IsSuccess:   true,
		},
	}
}

func FailLockoutAuthenticationResult(lockOut int, errors ...string) *AuthenticationResult {
	return &AuthenticationResult{
		Lockout: lockOut,
		AccessResult: AccessResult{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Error:      errors,
		},
	}
}

func SuccessLockoutAuthenticationResult(lockout int) *AuthenticationResult {
	return &AuthenticationResult{
		Lockout: lockout,
		AccessResult: AccessResult{
			StatusCode:  http.StatusLocked,
			IsSuccess:   false,
			IsLockedOut: true,
		},
	}
}

func FailAuthenticationResult(errors ...string) *AuthenticationResult {
	return &AuthenticationResult{
		AccessResult: AccessResult{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Error:      errors,
		},
	}
}
func SuccessAccessTwoFactorResult(isTwoFactorEnabled bool, token string) *AccessResult {
	return &AccessResult{
		StatusCode:         http.StatusOK,
		IsSuccess:          true,
		IsTwoFactorEnabled: isTwoFactorEnabled,
	}
}

func SuccessAuthenticateTwoFactorResult(isTwoFactorEnabled bool, token string) *AuthenticationResult {
	return &AuthenticationResult{
		TwoFactorToken: token,
		AccessResult: AccessResult{
			StatusCode:         http.StatusOK,
			IsSuccess:          true,
			IsTwoFactorEnabled: isTwoFactorEnabled,
		},
	}
}

func FailAccessResult(error ...string) *AccessResult {
	return &AccessResult{
		StatusCode: http.StatusBadRequest,
		IsSuccess:  false,
		Error:      error,
	}
}
