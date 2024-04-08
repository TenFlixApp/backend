package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"backend/data"
)

func main() {
	// Charchement .env
	err := godotenv.Load()
	// Si erreur, on plente
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
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

	// Lancement du serveur
	router.Run("localhost:8088")
}
