package Utilities

import (
	"SM/models"

	"github.com/golang-jwt/jwt/v5"
)

type IdentityClaim struct {
	Id       int
	Role     string
	Username string
}

func ExtractIdentity(claim jwt.MapClaims) *IdentityClaim {
	return &IdentityClaim{
		Id:       int(claim["id"].(float64)),
		Username: claim["username"].(string),
		Role:     extractRole(claim["role"].(float64)),
	}
}

func extractRole(role float64) string {
	if name, ok := models.Roles[models.Role(role)]; ok {
		return name
	}
	return "none"
}
