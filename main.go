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
	_ = godotenv.Load()

	data.ConnectToDB()
	defer data.CloseDB()

	router := gin.Default()

	configCors := cors.DefaultConfig()

	configCors.AllowAllOrigins, configCors.AllowCredentials = true, true
	configCors.AddAllowHeaders("Authorization")
	configCors.AddAllowHeaders("creditential")

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

	// playlists
	router.POST("/playlist", routes.CreatePlaylistRoute)
	router.DELETE("/playlist/:id", routes.DeletePlaylistRoute)
	router.POST("/playlist/media", routes.AddMediaToPlaylistRoute)
	router.DELETE("/playlist/media", routes.DeleteMediaFromPlaylistRoute)
	router.GET("/playlist/:id", routes.GetPlaylistRoute)
	router.GET("playlists", routes.GetPlaylistsFromUserRoute)

	// commentaires

	_ = router.Run(":8088")
}
