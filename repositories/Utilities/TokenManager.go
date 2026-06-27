package Utilities

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	b64 "encoding/base64"
)

type TokenManager struct {
	Logger          *slog.Logger
	SecretKey       string
	DurationMins    string
	DurationHrs     string
	RefreshTokenKey string
	TokenSecretKey  string
}

func (t *TokenManager) GenerateEmailToken(email string) (string, error) {
	i, err := strconv.ParseInt(t.DurationMins, 10, 64)
	if err != nil {
		t.Logger.Error(err.Error())
		return "", err
	}

	duration := time.Duration(i) * time.Hour

	TokenString := email + "_" + duration.String() + "_" + t.SecretKey

	h := b64.StdEncoding.EncodeToString([]byte(TokenString))

	return h, nil
}

func (t *TokenManager) ValidateEmailToken(email string, token string) bool {
	h, err := b64.StdEncoding.DecodeString(token)

	if err != nil {
		t.Logger.Error(err.Error())
		return false
	}

	decodedToken := strings.Split(string(h), "_")

	if len(decodedToken) < 3 {
		return false
	}

	if decodedToken[0] == email && decodedToken[2] == t.SecretKey {
		return true
	}
	return false
}

func (t *TokenManager) DecodeEmailToken(token string) (string, error) {
	h, err := b64.StdEncoding.DecodeString(token)
	if err != nil {
		t.Logger.Error(err.Error())
		return "", err
	}
	decodedToken := strings.Split(string(h), "_")
	if len(decodedToken) < 3 {
		return "", errors.New("invalid token")
	}
	if decodedToken[2] != t.SecretKey {
		return "", errors.New("invalid token")
	}

	return decodedToken[0], nil

}

func (t *TokenManager) GenerateTokenWithClaims(claims map[string]any) (string, error) {

	newclaims := jwt.MapClaims{}

	for key, value := range claims {
		newclaims[key] = value
	}

	i, err := strconv.ParseInt(t.DurationMins, 10, 64)

	if err != nil {
		t.Logger.Error(err.Error())
		return "", err
	}

	newclaims["exp"] = time.Now().Add(time.Minute * time.Duration(i)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newclaims)

	tokenString, err := token.SignedString([]byte(t.TokenSecretKey))

	if err != nil {
		t.Logger.Error(err.Error())
		return "", err
	}
	return tokenString, nil
}

func (t *TokenManager) GenerateRefreshToken(id int) (string, error) {
	i, err := strconv.ParseInt(t.DurationHrs, 10, 64)
	if err != nil {
		t.Logger.Error(err.Error())
		return "", err
	}

	duration := time.Duration(i) * time.Hour

	userRefreshTokenString := strconv.Itoa(id) + "_" + duration.String() + "_" + t.RefreshTokenKey

	h := b64.StdEncoding.EncodeToString([]byte(userRefreshTokenString))

	return h, nil
}

func (t *TokenManager) VerifyToken(token string) bool {
	validToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.TokenSecretKey), nil
	})
	if err != nil {
		t.Logger.Error(err.Error())
	}

	return validToken.Valid
}

func (t *TokenManager) DecodeToken(token string) (jwt.MapClaims, error) {
	validToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(t.TokenSecretKey), nil
	})
	if err != nil {
		t.Logger.Error(err.Error())
		return nil, err
	}

	claims, ok := validToken.Claims.(jwt.MapClaims)
	if !ok || !validToken.Valid {
		t.Logger.Error(err.Error())
		return nil, err
	}

	return claims, nil
}
func TokenManagerCons(handler *slog.JSONHandler) TokenManager {
	return TokenManager{
		Logger:          slog.New(handler),
		SecretKey:       os.Getenv("TOKEN_SECRET"),
		TokenSecretKey:  os.Getenv("SECRET_KEY"),
		RefreshTokenKey: os.Getenv("REFRESH_TOKEN_SECRET"),
		DurationHrs:     os.Getenv("REFRESH_EXPIRATION_HR"),
		DurationMins:    os.Getenv("TOKEN_EXPIRATION_MINS"),
	}
}

var _ ITokenManager = (*TokenManager)(nil)
