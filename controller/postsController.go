package controller

import (
	db "database-migrate/db/sqlc"
	"database-migrate/pkg/dbconn"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Image   sql.NullString `json:"image"`
	UserID  int32  `json:"user_id" binding:"required"`
}
type UpdatePostInput struct {
    Title   string  `json:"title" binding:"required"`
    Content string  `json:"content" binding:"required"`
    Image   *string `json:"image"` // opsional
}


func GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
        return
    }

    queries := db.New(dbconn.DB)
    user, err := queries.GetPostByID(c, int32(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "post not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func GetListPostsByUserID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	queries := db.New(dbconn.DB)
	users, err := queries.GetListPostsByUserID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "posts not found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetListPosts(c *gin.Context) {
	queries := db.New(dbconn.DB)
	users, err := queries.GetListPosts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "posts not found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func CreatePost(c *gin.Context) {
	queries := db.New(dbconn.DB)

	var input Post
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}
	user, err := queries.CreatePost(c, db.CreatePostParams{
		Title:   input.Title,
		Content: input.Content,
		Image:   input.Image,
		UserID:  input.UserID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "post not created"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdatePost(c *gin.Context) {
	queries := db.New(dbconn.DB)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}
	var input UpdatePostInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var image sql.NullString
	if input.Image != nil {
		image = sql.NullString{String: *input.Image, Valid: true}
	} else {
		image = sql.NullString{Valid: false}
	}

	user, err := queries.UpdatePost(c, db.UpdatePostParams{
		ID:      int32(id),
		Title:   input.Title,
		Content: input.Content,
		Image:   image,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "post not updated"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func DeletePost(c *gin.Context) {
	queries := db.New(dbconn.DB)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	err = queries.DeletePost(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "post not deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}

