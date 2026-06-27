package Utilities

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type PasswordManager struct {
	Logger *slog.Logger
}

func PasswordManagerCons(handler *slog.JSONHandler) *PasswordManager {
	return &PasswordManager{
		Logger: slog.New(handler),
	}
}

var _ IPasswordManager = (*PasswordManager)(nil)

func (p *PasswordManager) VerifyPassword(password string, currentPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(password))
	return err == nil
}

func (p *PasswordManager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
