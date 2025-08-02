package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	db "database-migrate/db/sqlc"
	"database-migrate/pkg/dbconn"
	"database-migrate/utils"
)

type User struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type UpdateUserInput struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func RegisterUser(c *gin.Context) {
    var input User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
        return
    }

    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "error hashing password"})
        return
    }

    queries := db.New(dbconn.DB)
    user, err := queries.CreateUser(c, db.CreateUserParams{
        Name: input.Name,
        Email: input.Email,
        Password: hashedPassword,
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Register Failed"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
    queries := db.New(dbconn.DB)

    userIDInterface, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"message": "user_id not found"})
        return
    }

    userIDFromToken, ok := userIDInterface.(int32)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user_id"})
        return
    }

    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
        return
    }

    if userIDFromToken != int32(id) {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not Allowed to update this user"})
        return
    }
    
    var input UpdateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    oldUser, err := queries.GetUserByID(c, int32(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
        return
    }
    updatedName := oldUser.Name
    if input.Name != nil && *input.Name != "" {
        updatedName = *input.Name
    }
    updatedEmail := oldUser.Email
    if input.Email != nil && *input.Email != "" {
        updatedEmail = *input.Email
    }

    updatedPassword := oldUser.Password
    if input.Password != nil && *input.Password != "" {
        hashedPassword, err := utils.HashPassword(*input.Password)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "error hashing password"})
            return
        }
        updatedPassword = hashedPassword
    }

    updatedUser, err := queries.UpdateUser(c, db.UpdateUserParams{
        ID: int32(id),
        Name: updatedName,
        Email: updatedEmail,
        Password: updatedPassword,
    })
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "user not updated"})
        return
    }

    c.JSON(http.StatusOK, updatedUser)
}

func DeleteUser(c *gin.Context) {
    queries := db.New(dbconn.DB)

    
    userIDInterface, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"message": "user_id not found"})
        return
    }

    userIDFromToken, ok := userIDInterface.(int32)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user_id"})
        return
    }


    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
        return
    }

     if userIDFromToken != int32(id) {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not Allowed to Delete this user"})
        return
    }

    err = queries.DeleteUser(c, int32(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func LoginUser(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
        return
    }

    queries := db.New(dbconn.DB)
    user, err := queries.GetUserByEmail(c, input.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
        return
    }

    if !utils.CheckPasswordHash(input.Password, user.Password) {
        c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect password"})
        return
    }

    token, err := utils.GenerateToken(user.Email, user.ID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "error generating token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message" : "login successful!", "token": token})
}