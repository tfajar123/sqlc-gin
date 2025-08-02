package routes

import (
	"database-migrate/controller"
	"database-migrate/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/register", controller.RegisterUser )
	r.POST("/login", controller.LoginUser )
	
	r.GET("/posts/:id", controller.GetPostByID )
	r.GET("/posts", controller.GetListPosts )
	r.GET("/posts/user/:id", controller.GetListPostsByUserID )
	
	authenticated := r.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.PUT("/users/:id", controller.UpdateUser )
	authenticated.DELETE("/users/:id", controller.DeleteUser )

	authenticated.POST("/posts", controller.CreatePost )
	authenticated.PUT("/posts/:id", controller.UpdatePost )
	authenticated.DELETE("/posts/:id", controller.DeletePost )

	r.Static("/uploads", "./uploads")
}