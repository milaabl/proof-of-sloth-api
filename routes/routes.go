package routes

import (
	"github.com/gin-gonic/gin"

    "github.com/milaabl/proof-of-sloth-api/controllers"
)

func InitRoutes(router *gin.Engine)  {
    router.GET("/albums", controllers.GetAlbums())

	router.GET("/albums/:id", controllers.GetAlbumByID())

	router.POST("/albums", controllers.CreateAlbum())

    router.PUT("/albums/:id", controllers.EditAlbum())

    router.DELETE("/albums/:id", controllers.DeleteAlbum())
}
