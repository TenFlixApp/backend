package routes

import (
	"backend/data"
	"backend/helpers"
	"bytes"
	"encoding/json"
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

type FileManagerResponse struct {
	Links struct {
		Cover  string `json:"cover"`
		Stream string `json:"stream"`
	} `json:"_links"`
	Id   string `json:"id"`
	Type string `json:"type"`
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
	id := c.Param("id")
	param, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id"})
		return
	}

	errData := data.DeleteMedia(param)
	if errData != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errData.Message})
		return
	}

	// TODO : delete in file manager

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}

func UploadMediaRoute(c *gin.Context) {
	var form FormMedia
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	uuid := rerouteFiles(form)

	idUser, err := helpers.GetIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	errData := data.CreateMedia(idUser, form.Title, uuid, form.Description)
	if errData != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errData.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "media uploaded"})
}

func GetUserChannelRoute(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	user, err := data.GetUserInfo(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medias, err := data.GetMediaFromCreator(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "medias": medias})
}

func SearchMediaRoute(c *gin.Context) {
	type Query struct {
		SearchTerm string `json:"q" binding:"required"`
	}

	var query Query
	if err := c.ShouldBind(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(query.SearchTerm)
	medias, err := data.SearchMedia(query.SearchTerm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"medias": medias})
}

func GetRandomMediaRoute(c *gin.Context) {
	count, success := c.GetQuery("count")

	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	_, err := strconv.Atoi(count)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "count must be a valid integer"})
		return
	}

	resp, err := http.Get(os.Getenv("FILE_MANAGER_ROUTE") + "files/random?count=" + count + "&type=media")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching data"})
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while reading data"})
		return
	}

	type BodyType struct {
		Medias []FileManagerResponse `json:"medias"`
	}

	var responseBody BodyType
	if err := json.Unmarshal(respBody, &responseBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var mediaIds = make([]string, len(responseBody.Medias))
	for i, media := range responseBody.Medias {
		mediaIds[i] = media.Id
	}

	medias, err := data.GetMedias(mediaIds)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"medias": medias})
}
