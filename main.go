package main

import (
	"database-migrate/pkg/dbconn"
	"database-migrate/routes"

	"github.com/gin-gonic/gin"
)

func main() {
    dbconn.Connect()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")


}
