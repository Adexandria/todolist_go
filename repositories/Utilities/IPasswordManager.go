package Utilities

type PasswordManager interface {
	VerifyPassword(currentPassword string, password string) bool
	HashPassword(password string) (string, error)
}
