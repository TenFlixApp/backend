package security

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ExtractFromToken(c *gin.Context, key string) (string, error) {
	// Récupérer le token d'authentification du header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	// Vérifier que le format du header est Bearer {token}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	tokenString := parts[1]

	// Parse le token JWT sans vérification de signature
	token, err := jwt.Parse(tokenString, nil)
	if err != nil {
		return "", fmt.Errorf("error parsing token: %w", err)
	}

	// Vérifier si le token est valide
	if !token.Valid {
		return "", errors.New("invalid token")
	}

	// Extraire l'email du token
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
