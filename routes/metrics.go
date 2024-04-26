package routes

import (
	"backend/data"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboardStatsRoute(c *gin.Context) {
	loginStats, err := services.GetLoginStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	registerStats, err := services.GetRegisterStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userCount, err := data.CountUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	creatorsCount, err := data.CountCreators()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	adminUserCount, disabledUserCount, err := services.GetUserMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mediaCount, err := data.CountMedias()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"login":             loginStats,
		"register":          registerStats,
		"mediaCount":        mediaCount,
		"userCount":         userCount,
		"creatorsCount":     creatorsCount,
		"adminUserCount":    adminUserCount,
		"disabledUserCount": disabledUserCount,
	})
}
