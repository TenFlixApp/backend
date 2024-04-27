package routes

import (
	"backend/data"
	"backend/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePlaylistRoute(c *gin.Context) {
	type Playlist struct {
		Titre string `json:"titre" binding:"required"`
	}

	var playlist Playlist
	if err := c.ShouldBind(&playlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idUser, err := helpers.GetIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	errData := data.CreatePlaylist(idUser, playlist.Titre)
	if errData != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errData.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Playlist created"})
}

func DeletePlaylistRoute(c *gin.Context) {
	id := c.Param("id")
	param, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id"})
		return
	}

	errData := data.DeletePlaylist(param)
	if errData != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errData.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted"})
}

func AddMediaToPlaylistRoute(c *gin.Context) {
	type PlaylistMedia struct {
		IdPlaylist int    `json:"id_playlist" binding:"required"`
		UUIDMedia  string `json:"uuid_media" binding:"required"`
	}

	var playlistMedia PlaylistMedia
	if err := c.ShouldBind(&playlistMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errData := data.AddMediaToPlaylist(playlistMedia.IdPlaylist, playlistMedia.UUIDMedia)
	if errData != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errData.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media added to playlist"})
}

func DeleteMediaFromPlaylistRoute(c *gin.Context) {
	type PlaylistMedia struct {
		IdPlaylist int    `json:"id_playlist" binding:"required"`
		UUIDMedia  string `json:"uuid_media" binding:"required"`
	}

	var playlistMedia PlaylistMedia
	if err := c.ShouldBind(&playlistMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errData := data.DeleteMediaFromPlaylist(playlistMedia.IdPlaylist, playlistMedia.UUIDMedia)
	if errData != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errData.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted from playlist"})
}

func GetPlaylistsFromUserRoute(c *gin.Context) {
	idUser, err := helpers.GetIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	playlists, errData := data.GetPlaylistsFromUser(idUser)
	if errData != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errData.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"playlists": playlists})
}

func GetPlaylistRoute(c *gin.Context) {
	id := c.Param("id")
	param, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id"})
		return
	}

	playlist, err := data.GetMediaFromPlaylist(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"medias": playlist})
}
