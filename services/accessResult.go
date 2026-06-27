package services

import "net/http"

type AccessResult struct {
	StatusCode         int
	AccessToken        string
	RefreshToken       string
	IsSuccess          bool
	error              []string
	IsTwoFactorEnabled bool
	Token              string
}

func SuccessAccessResult(accessToken string, refreshToken string) *AccessResult {
	return &AccessResult{
		StatusCode:   http.StatusOK,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsSuccess:    true,
	}
}

func SuccessAccessTwoFactorResult(isTwoFactorEnabled bool, token string) *AccessResult {
	return &AccessResult{
		StatusCode:         http.StatusOK,
		IsSuccess:          true,
		IsTwoFactorEnabled: isTwoFactorEnabled,
		Token:              token,
	}
}

func FailAccessResult(error ...string) *AccessResult {
	return &AccessResult{
		StatusCode: http.StatusBadRequest,
		IsSuccess:  false,
		error:      error,
	}
}
