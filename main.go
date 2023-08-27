package main

import (
    "github.com/gin-gonic/gin"
    "github.com/milaabl/proof-of-sloth-api/db"
    "github.com/milaabl/proof-of-sloth-api/routes"
)


func main() {
    router := gin.Default()

    db.ConnectDB()

    routes.InitRoutes(router)

    router.Run("localhost:8080")
}
