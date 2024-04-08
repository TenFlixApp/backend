package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var input struct {
	Email    string `json:"email"`
	Nom      string `json:"nom"`
	Prenom   string `json:"prenom"`
	Password string `json:"password"`
}

func RegisterRoute(c *gin.Context) {
	var user input
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	jsonBody := []byte(`{"username":` + user.Email + `, "password"}`)

	// resp, err := http.Post(os.Getenv("GUARDIAN_ROOT"), {

	// })
}
