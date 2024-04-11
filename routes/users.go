package routes

import (
	"backend/data"
	"backend/helpers"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type input struct {
	Email    string `json:"email"`
	Nom      string `json:"nom"`
	Prenom   string `json:"prenom"`
	Password string `json:"password"`
}

type FormAvatar struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func RegisterRoute(c *gin.Context) {
	var user input
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Nom == "" || user.Prenom == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	jsonBody := []byte(`{"username":"` + user.Email + `", "password":"` + user.Password + `", "rights": 1}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", os.Getenv("GUARDIAN_ROUTE")+"register", bodyReader)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), resBody)
		return
	}

	dataError := data.RegisterUser(user.Email, user.Nom, user.Prenom)
	if dataError != nil {
		c.JSON(http.StatusConflict, gin.H{"error": dataError.Message})
		return
	}

	var responseBody interface{}
	if err := json.Unmarshal(resBody, &responseBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to unmarshal response body"})
		return
	}

	c.JSON(http.StatusCreated, responseBody)
}

func ChangeAvatarRoute(c *gin.Context) {
	var form FormAvatar

	// Bind form data
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO - Appel du file manager
	uuid := "test"

	idUser, err := helpers.GetIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	errData := data.ChangeAvatar(uuid, idUser)
	if errData != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errData.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "avatar changed"})
	return
}
