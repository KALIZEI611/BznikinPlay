package main

import (
	"console-rental/database"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-super-secret-jwt-key-change-in-production")

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

func generateToken(userID int, username, email string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Email:    email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "No authorization header"})
            c.Abort()
            return
        }
        
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        claims := &Claims{}
        
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        
        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("email", claims.Email)
        c.Next()
    }
}

func main() {
    // Подключение к БД
    if err := database.InitDB(); err != nil {
        log.Printf("⚠️ WARNING: Could not connect to database: %v", err)
        log.Println("Server will continue with limited functionality")
    }
    
    r := gin.Default()
    log.Println("🚀 Starting ConsoleRent backend with JWT auth...")
    
    // CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            "https://bznikin-play.vercel.app",
            "http://localhost:3000",
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))
    
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
    })
    
    // Каталог консолей
    r.GET("/api/consoles", func(c *gin.Context) {
        c.JSON(200, []gin.H{
            {"id": 1, "type": "PS5", "model": "PlayStation 5 Standard", "price_per_day": 800, "is_available": true, "image_url": "https://avatars.mds.yandex.net/get-mpic/13230222/2a000001969fb44f0b2b0473bcfe73eb4de4/orig"},
            {"id": 2, "type": "PS5", "model": "PlayStation 5 Digital", "price_per_day": 800, "is_available": true, "image_url": "https://avatars.mds.yandex.net/i?id=0b869e3e8145ca09fba9fa1e77702f95_l-4355007-images-thumbs&n=13"},
            {"id": 3, "type": "PS4", "model": "PlayStation 4 Slim", "price_per_day": 500, "is_available": true, "image_url": "https://avatars.mds.yandex.net/get-mpic/5173149/2a0000019180dbf1814ebb7ae678faa8667a/orig"},
            {"id": 4, "type": "XBOX", "model": "Xbox Series X", "price_per_day": 800, "is_available": true, "image_url": "https://hatiko.ru/wa-data/public/blog/img/photo_2024-01-30_12-36-25.jpg"},
            {"id": 5, "type": "XBOX", "model": "Xbox Series S", "price_per_day": 800, "is_available": true, "image_url": "https://api.2droida.ru/storage/products/b11679d3f73924628580ceea19c6e9eb/5153/7a33c4eca9aca748f6753e1cb3f90101.jpg"},
            {"id": 6, "type": "XBOX", "model": "Xbox One X", "price_per_day": 500, "is_available": true, "image_url": "https://gameshock174.ru/upload/iblock/bbb/bbbdbb58eb867f610801394b5ef15e3a.jpg"},
        })
    })
    
    // Регистрация
    r.POST("/api/register", func(c *gin.Context) {
        var req struct {
            Username string `json:"username"`
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "Invalid request"})
            return
        }
        
        // Проверяем, существует ли пользователь
        var exists bool
        database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
        if exists {
            c.JSON(400, gin.H{"error": "User with this email already exists"})
            return
        }
        
        // Хешируем пароль
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to hash password"})
            return
        }
        
        // Создаём пользователя
        var userID int
        err = database.DB.QueryRow(
            "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
            req.Username, req.Email, string(hashedPassword),
        ).Scan(&userID)
        
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to create user: " + err.Error()})
            return
        }
        
        token, err := generateToken(userID, req.Username, req.Email)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to generate token"})
            return
        }
        
        c.JSON(201, gin.H{
            "message": "User created successfully",
            "token": token,
            "user": gin.H{
                "id": userID,
                "username": req.Username,
                "email": req.Email,
            },
        })
    })
    
    // Вход
    r.POST("/api/login", func(c *gin.Context) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "Invalid request"})
            return
        }
        
        var userID int
        var username, hashedPassword string
        err := database.DB.QueryRow(
            "SELECT id, username, password FROM users WHERE email = $1",
            req.Email,
        ).Scan(&userID, &username, &hashedPassword)
        
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid email or password"})
            return
        }
        
        err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid email or password"})
            return
        }
        
        token, err := generateToken(userID, username, req.Email)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to generate token"})
            return
        }
        
        c.JSON(200, gin.H{
            "message": "Login successful",
            "token": token,
            "user": gin.H{
                "id": userID,
                "username": username,
                "email": req.Email,
            },
        })
    })
    
    // Профиль
    r.GET("/api/user/profile", authMiddleware(), func(c *gin.Context) {
        userID := c.GetInt("user_id")
        
        var username, email string
        var createdAt time.Time
        err := database.DB.QueryRow(
            "SELECT username, email, created_at FROM users WHERE id = $1",
            userID,
        ).Scan(&username, &email, &createdAt)
        
        if err != nil {
            c.JSON(404, gin.H{"error": "User not found"})
            return
        }
        
        c.JSON(200, gin.H{
            "id": userID,
            "username": username,
            "email": email,
            "created_at": createdAt,
        })
    })
    
    r.PUT("/api/user/profile", authMiddleware(), func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Profile updated"})
    })
    
    // Аренда
    var rentals = []gin.H{}
    var nextId = 1
    
    r.POST("/api/rentals", authMiddleware(), func(c *gin.Context) {
        userID := c.GetInt("user_id")
        
        var req struct {
            ConsoleID       int    `json:"console_id"`
            StartDate       string `json:"start_date"`
            EndDate         string `json:"end_date"`
            DeliveryAddress string `json:"delivery_address"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "Invalid request"})
            return
        }
        
        var consolePrice float64
        
        switch req.ConsoleID {
        case 1, 2, 4, 5:
            consolePrice = 800
        case 3, 6:
            consolePrice = 500
        default:
            consolePrice = 800
        }
        
        rental := gin.H{
            "id":               nextId,
            "user_id":          userID,
            "total_price":      consolePrice,
            "delivery_address": req.DeliveryAddress,
            "status":           "active",
        }
        rentals = append(rentals, rental)
        nextId++
        
        c.JSON(201, gin.H{
            "rental_id":   rental["id"],
            "total_price": consolePrice,
            "message":     "Rental created successfully",
        })
    })
    
    r.GET("/api/my-rentals", authMiddleware(), func(c *gin.Context) {
        userID := c.GetInt("user_id")
        userRentals := []gin.H{}
        for _, r := range rentals {
            if r["user_id"] == userID {
                userRentals = append(userRentals, r)
            }
        }
        c.JSON(200, userRentals)
    })
    
    r.PUT("/api/rentals/:id/return", authMiddleware(), func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Console returned successfully"})
    })
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("✅ Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}