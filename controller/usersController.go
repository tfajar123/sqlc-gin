package controller

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    db "database-migrate/db/sqlc"
    "database-migrate/pkg/dbconn"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
        return
    }

    queries := db.New(dbconn.DB)
    user, err := queries.GetUserByID(c, int32(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func GetListUsers(c *gin.Context) {
    queries := db.New(dbconn.DB)
    users, err := queries.GetListUsers(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "users not found"})
        return
    }

    c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
    queries := db.New(dbconn.DB)

    var input User
    if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}
    user, err := queries.CreateUser(c, db.CreateUserParams{
        Name: input.Name,
        Email: input.Email,
    }) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "user not created"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
    queries := db.New(dbconn.DB)
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
        return
    }
    var input User
    if err := c.BindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
        return
    }
    user, err := queries.UpdateUser(c, db.UpdateUserParams{
        ID: int32(id),
        Name: input.Name,
        Email: input.Email,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "user not updated"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
    queries := db.New(dbconn.DB)
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
        return
    }

    err = queries.DeleteUser(c, int32(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "user not deleted"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}