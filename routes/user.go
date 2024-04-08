package routes

import (
	helpers "backend/Helpers"
	"backend/data"
	"bytes"
	"encoding/json"
	"fmt"
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if user.Nom == "" || user.Prenom == "" || user.Email == "" || user.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	jsonBody := []byte(`{"username":"` + user.Email + `", "password":"` + user.Password + `"}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", os.Getenv("GUARDIAN_ROUTE")+"register", bodyReader)

	if err != nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "unable to create request"})
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "unable to read response body"})
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), resBody)
		return
	}

	dataError := data.RegisterUser(user.Email, user.Nom, user.Prenom)
	if dataError != nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": dataError.Message})
		return
	}

	var responseBody interface{}
	if err := json.Unmarshal(resBody, &responseBody); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "unable to unmarshal response body"})
		return
	}

	c.IndentedJSON(http.StatusCreated, responseBody)
}

func ChangeAvatarRoute(c *gin.Context) {
	var form FormAvatar

	// Bind form data
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	// TODO - Appel du file manager
	uuid := "test"

	idUser, err := helpers.GetIdFromToken(c)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	errData := data.ChangeAvatar(uuid, idUser)
	if errData != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": errData.Message})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "avatar changed"})
	return
}
