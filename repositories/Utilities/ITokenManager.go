package Utilities

import "github.com/golang-jwt/jwt/v5"

type ITokenManager interface {
	GenerateJWTTokenWithClaims(claims map[string]any) (string, error)
	VerifyJWTToken(token string) bool
	DecodeJWTToken(token string) (jwt.MapClaims, error)

	GenerateRefreshToken(id int) (string, error)
	VerifyRefreshToken(id int, token string) bool

	GenerateToken(email string) (string, error)
	ValidateToken(email string, token string) bool
	DecodeToken(token string) (string, error)
}
