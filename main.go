package main

import (
	"gallery-api/db"
	"gallery-api/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	conn, _ := db.GetDBConnection().DB()
	db.InitMigration()

	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	router := gin.Default()

	router.POST("/upload", handler.UploadPhoto)
	router.GET("/list", handler.ListPhotos)
	router.GET("/take/image/:id", handler.GetImage)
	router.GET("/take/preview/:id", handler.GetPreview)
	router.DELETE("/delete/:id", handler.DeletePhoto)

	router.Run(":8080")
}
