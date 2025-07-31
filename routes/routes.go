package routes

import (
	"github.com/gin-gonic/gin"
	"database-migrate/controller"

)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/users/:id", controller.GetUser )
	r.GET("/users", controller.GetListUsers )
	r.POST("/users", controller.CreateUser )
	r.PUT("/users/:id", controller.UpdateUser )
	r.DELETE("/users/:id", controller.DeleteUser )
}