package security

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("GO_SECRET_KEY_ACCESS_TOKEN")), nil
}

func ExtractFromToken(c *gin.Context, key string) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return "", fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	claim, ok := claims[key].(string)
	if !ok || claim == "" {
		return "", errors.New("invalid token payload")
	}

	return claim, nil
}
