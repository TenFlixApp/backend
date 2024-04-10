package routes

import (
	"backend/data"
	"backend/helpers"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FormMedia struct {
	Title       string                `form:"title" binding:"required"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
	Cover       *multipart.FileHeader `form:"cover" binding:"required"`
	Description string                `form:"title" binding:"required"`
}

func rerouteFiles(form FormMedia) string {
	fileManagerRoute := os.Getenv("FILE_MANAGER_ROUTE") + "upload"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, _ := form.File.Open()
	cover, _ := form.Cover.Open()

	filePart, _ := writer.CreateFormFile("video", form.File.Filename)
	_, _ = io.Copy(filePart, file)

	coverPart, _ := writer.CreateFormFile("cover", form.Cover.Filename)
	_, _ = io.Copy(coverPart, cover)

	_ = writer.Close()

	req, err := http.NewRequest("POST", fileManagerRoute, body)
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(respBody)
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

	uuid := rerouteFiles(form)

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

func GetRandomMediaRoute(c *gin.Context) {
	count, success := c.GetQuery("count")

	if !success {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	_, err := strconv.Atoi(count)
	if err != nil {
		c.String(http.StatusBadRequest, "count must be a valid integer")
		return
	}

	// send get request to file manager on route /files/random?count=:count&type=media
	resp, err := http.Get(os.Getenv("FILE_MANAGER_ROUTE") + "/files/random?count=" + count + "&type=media")
	if err != nil {
		c.String(http.StatusInternalServerError, "error while fetching data")
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "error while reading data")
		return
	}

	c.JSON(http.StatusOK, gin.H{"medias": respBody})
}
