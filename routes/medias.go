package routes

import (
	"backend/data"
	"backend/helpers"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FormMedia struct {
	Title       string                `form:"title" binding:"required"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
	Cover       *multipart.FileHeader `form:"cover" binding:"required"`
	Description string                `form:"title" binding:"required"`
}

func DeleteMediaRoute(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur à supprimer depuis les paramètres de l'URL
	id := c.Param("id")
	param, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "incorrect id"})
		return
	}

	errData := data.DeleteMedia(param)
	if errData != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": errData.Message})
		return
	}

	// to do : delete in file manager

	c.IndentedJSON(http.StatusOK, gin.H{"message": "file deleted"})
	return
}

func UploadMediaRoute(c *gin.Context) {
	var form FormMedia

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

	errData := data.CreateMedia(idUser, form.Title, uuid, form.Description)
	if errData != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": errData.Message})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "media uploaded"})
	return
}

func GetUserChannelRoute(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis les paramètres de l'URL
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	user, err := data.GetInfoUser(userID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medias, err := data.GetMediaFromCreator(userID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"user": user, "medias": medias})
}

func SearchMediaRoute(c *gin.Context) {
	type Query struct {
		SearchTerm string `json:"q" binding:"required"`
	}

	var query Query
	if err := c.ShouldBind(&query); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	fmt.Println(query.SearchTerm)
	medias, err := data.SearchMedia(query.SearchTerm)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"medias": medias})
}
