package controller

import (
	db "database-migrate/db/sqlc"
	"database-migrate/pkg/dbconn"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)


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

	title := c.PostForm("title")
	content := c.PostForm("content")
	userID := c.PostForm("user_id")

	uid, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user_id"})
		return
	}

	file, err := c.FormFile("image")
	var filename string
	if err == nil {
		filename = fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), file.Filename)

		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	image := sql.NullString{
		String: filename,
		Valid: filename != "",
	}

	post, err := queries.CreatePost(c, db.CreatePostParams{
		Title:   title,
		Content: content,
		Image:   image,
		UserID:  int32(uid),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "post not created"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func UpdatePost(c *gin.Context) {
	queries := db.New(dbconn.DB)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}


	oldPost, err := queries.GetPostByID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "post not found"})
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")
	// userID := c.PostForm("user_id")

	file, err := c.FormFile("image")
	var image sql.NullString
	if err == nil {
		filename := fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), file.Filename)

		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		image = sql.NullString{
			String: filename,
			Valid: true,
		}

		if oldPost.Image.Valid {
			if err := os.Remove(oldPost.Image.String); err != nil {
				log.Println("failed to remove old image:", err)
			}
		}
	} else {
		image = oldPost.Image
	}

	post, err := queries.UpdatePost(c, db.UpdatePostParams{
		ID:      int32(id),
		Title:   title,
		Content: content,
		Image:   image,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "post not updated"})
		return
	}

	c.JSON(http.StatusOK, post)
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

