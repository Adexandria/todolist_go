package mocks

import (
	"SM/repositories/Utilities"
	"log/slog"
)

type passwordHandler struct {
	Logger *slog.Logger
}

func (p passwordHandler) VerifyPassword(currentPassword string, password string) bool {
	return currentPassword == password
}

func (p passwordHandler) HashPassword(password string) (string, error) {
	return password, nil
}

func NewPasswordHandler(handler *slog.JSONHandler) Utilities.PasswordManager {
	return &passwordHandler{
		Logger: slog.New(handler),
	}
}

var _ Utilities.PasswordManager = (*passwordHandler)(nil)
