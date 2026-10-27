package Utilities

type IPasswordManager interface {
	VerifyPassword(currentPassword string, password string) bool
	HashPassword(password string) (string, error)
}
