package Utilities

import "github.com/golang-jwt/jwt/v5"

type ITokenManager interface {
	GenerateTokenWithClaims(claims map[string]any) (string, error)
	GenerateRefreshToken(id int) (string, error)
	GenerateEmailToken(email string) (string, error)
	ValidateEmailToken(email string, token string) bool
	DecodeEmailToken(token string) (string, error)
	VerifyToken(token string) bool
	DecodeToken(token string) (jwt.MapClaims, error)
}
