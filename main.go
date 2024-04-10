package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"backend/data"
	"backend/routes"
)

func main() {
	// Charchement .env
	_ = godotenv.Load()

	// Tentative de connexion à la base de donnée
	data.ConnectToDB()
	// Programmation de la fermeture de la base de données à la fermeture du programme
	defer data.CloseDB()

	// Création du routeur avec le framework GIN
	router := gin.Default()

	// Création de la configuration par défaut des cors du serveur
	configCors := cors.DefaultConfig()

	// Modification des paramètres
	configCors.AllowAllOrigins, configCors.AllowCredentials = true, true
	configCors.AddAllowHeaders("Authorization")
	configCors.AddAllowHeaders("creditential")

	// Application de la nouvelle configuration
	router.Use(cors.New(configCors))

	// users
	router.POST("/user/register", routes.RegisterRoute)
	router.PUT("/user/avatar", routes.ChangeAvatarRoute)
	router.GET("/user/:id", routes.GetUserChannelRoute)

	// medias
	router.DELETE("/media/:id", routes.DeleteMediaRoute)
	router.POST("/media/upload", routes.UploadMediaRoute)
	router.POST("/media/search", routes.SearchMediaRoute)
	router.GET("/media/random", routes.GetRandomMediaRoute)

	// commentaires

	// Lancement du serveur
	_ = router.Run("localhost:8088")
}
