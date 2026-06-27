package Utilities

type IPasswordManager interface {
	VerifyPassword(password string, currentPassword string) bool
	HashPassword(password string) (string, error)
}
