package Utilities

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type passwordManager struct {
	Logger *slog.Logger
}

func NewPasswordManager(handler *slog.JSONHandler) PasswordManager {
	return &passwordManager{
		Logger: slog.New(handler),
	}
}

var _ PasswordManager = (*passwordManager)(nil)

func (p *passwordManager) VerifyPassword(currentPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(password))
	return err == nil
}

func (p *passwordManager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
