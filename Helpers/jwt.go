package helpers

import (
	"backend/data"
	"backend/security"
	"errors"

	"github.com/gin-gonic/gin"
)

func GetIdFromToken(c *gin.Context) (int, error) {
	email, err := security.ExtractFromToken(c, "username")
	if err != nil {
		return 0, errors.New("unable to get user from token")
	}

	idUser, errData := data.GetUserIDByEmail(email)
	if errData != nil {
		return 0, errors.New("unable to get user from token")
	}

	return idUser, nil
}
