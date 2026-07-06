package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Use a secure environment variable for this in production!
var secretKey = []byte("your-highly-secure-secret-key")

type MyCustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, role string) (string, error) {
	// 1. Create custom claims payload
	claims := MyCustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Short lifetime
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "my-go-api",
		},
	}

	// 2. Choose signing method (HS256 for symmetric, RS256 for asymmetric)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Sign the token with your secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
