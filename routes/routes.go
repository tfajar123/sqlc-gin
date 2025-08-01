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

	r.GET("/posts/:id", controller.GetPostByID )
	r.GET("/posts", controller.GetListPosts )
	r.GET("/posts/user/:id", controller.GetListPostsByUserID )
	r.POST("/posts", controller.CreatePost )
	r.PUT("/posts/:id", controller.UpdatePost )
	r.DELETE("/posts/:id", controller.DeletePost )
}