package utils

import (
	"fmt"
	"time"

	"upcycle_connect-api/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      string   `json:"user_id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, roles []string, permissions []string) (string, error) {
	claims := Claims{
		UserID:      userID,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecret())
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algorithme de signature inattendu: %v", token.Header["alg"])
		}
		return config.JWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token invalide")
	}

	return claims, nil
}

const emailVerifyPurpose = "email_verify"

type EmailVerificationClaims struct {
	UserID  string `json:"user_id"`
	Purpose string `json:"purpose"`
	jwt.RegisteredClaims
}

func GenerateEmailVerificationToken(userID string) (string, error) {
	claims := EmailVerificationClaims{
		UserID:  userID,
		Purpose: emailVerifyPurpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecret())
}

func ParseEmailVerificationToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &EmailVerificationClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algorithme de signature inattendu: %v", token.Header["alg"])
		}
		return config.JWTSecret(), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*EmailVerificationClaims)
	if !ok || !token.Valid || claims.Purpose != emailVerifyPurpose {
		return "", fmt.Errorf("token invalide")
	}

	return claims.UserID, nil
}
