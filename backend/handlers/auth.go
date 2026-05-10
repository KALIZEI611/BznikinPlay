package handlers

import (
	"console-rental/database"
	"console-rental/middleware"
	"console-rental/models"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-super-secret-jwt-key-change-in-production")

func Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    var exists bool
    err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
    if err == nil && exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
        return
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    
    var userID int
    err = database.DB.QueryRow(
        "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
        req.Username, req.Email, string(hashedPassword),
    ).Scan(&userID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }
    
    token, err := generateToken(userID, req.Username, req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "user_id": userID,
        "token": token,
        "user": gin.H{
            "id": userID,
            "username": req.Username,
            "email": req.Email,
        },
    })
}

func Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    var user models.User
    var hashedPassword string
    err := database.DB.QueryRow(
        "SELECT id, username, email, password FROM users WHERE email = $1",
        req.Email,
    ).Scan(&user.ID, &user.Username, &user.Email, &hashedPassword)
    
    if err == sql.ErrNoRows {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    
    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    token, err := generateToken(user.ID, user.Username, user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "token": token,
        "user": gin.H{
            "id": user.ID,
            "username": user.Username,
            "email": user.Email,
        },
    })
}

func generateToken(userID int, username, email string) (string, error) {
    claims := models.Claims{
        UserID:   userID,
        Username: username,
        Email:    email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}
func GetUserProfile(c *gin.Context) {
    userID := middleware.GetUserID(c)
    if userID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    var user models.User
    err := database.DB.QueryRow(
        "SELECT id, username, email, created_at FROM users WHERE id = $1",
        userID,
    ).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
    
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "id": user.ID,
        "username": user.Username,
        "email": user.Email,
        "created_at": user.CreatedAt,
    })
}

func UpdateUserProfile(c *gin.Context) {
    userID := middleware.GetUserID(c)
    if userID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    var req struct {
        Username        string `json:"username"`
        Email           string `json:"email"`
        CurrentPassword string `json:"current_password"`
        NewPassword     string `json:"new_password"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }
    
    var hashedPassword string
    err := database.DB.QueryRow("SELECT password FROM users WHERE id = $1", userID).Scan(&hashedPassword)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    
    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.CurrentPassword))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
        return
    }
    
    tx, err := database.DB.Begin()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    defer tx.Rollback()
    
    if req.Username != "" {
        _, err = tx.Exec("UPDATE users SET username = $1 WHERE id = $2", req.Username, userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
            return
        }
    }
    
    if req.Email != "" {
        var exists bool
        tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)", req.Email, userID).Scan(&exists)
        if exists {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
            return
        }
        
        _, err = tx.Exec("UPDATE users SET email = $1 WHERE id = $2", req.Email, userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email"})
            return
        }
    }
    
    if req.NewPassword != "" {
        if len(req.NewPassword) < 6 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
            return
        }
        
        hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
            return
        }
        
        _, err = tx.Exec("UPDATE users SET password = $1 WHERE id = $2", string(hashedNewPassword), userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
            return
        }
    }
    
    var username, email string
    err = tx.QueryRow("SELECT username, email FROM users WHERE id = $1", userID).Scan(&username, &email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated user data"})
        return
    }
    
    err = tx.Commit()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }
    
    token, err := generateToken(userID, username, email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Profile updated successfully",
        "token": token,
        "user": gin.H{
            "id": userID,
            "username": username,
            "email": email,
        },
    })
}