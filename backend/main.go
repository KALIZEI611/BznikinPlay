package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("my-secret-key")

type User struct {
    ID       int
    Username string
    Email    string
    Password string
}

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

// Временное хранилище пользователей (in-memory)
var users = []User{
    {ID: 1, Username: "testuser", Email: "test@example.com", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMy.MrCvqKqKqKqKqKqKqKqKqKqKqKqKq"},
}
var nextUserID = 2

// Временное хранилище аренд
var rentals = []gin.H{}
var nextRentalID = 1

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
    r := gin.Default()
    
    log.Println("🚀 Starting ConsoleRent backend...")

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

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
    })

    // ========== КАТАЛОГ КОНСОЛЕЙ ==========
    r.GET("/api/consoles", func(c *gin.Context) {
        c.JSON(200, []gin.H{
            {"id": 1, "type": "PS5", "model": "PlayStation 5 Standard", "price_per_day": 800, "is_available": true, "description": "Полная комплектация", "image_url": "https://avatars.mds.yandex.net/get-mpic/13230222/2a000001969fb44f0b2b0473bcfe73eb4de4/orig"},
            {"id": 2, "type": "PS5", "model": "PlayStation 5 Digital", "price_per_day": 800, "is_available": true, "description": "Цифровая версия", "image_url": "https://avatars.mds.yandex.net/i?id=0b869e3e8145ca09fba9fa1e77702f95_l-4355007-images-thumbs&n=13"},
            {"id": 3, "type": "PS4", "model": "PlayStation 4 Slim", "price_per_day": 500, "is_available": true, "description": "500 GB HDD", "image_url": "https://avatars.mds.yandex.net/get-mpic/5173149/2a0000019180dbf1814ebb7ae678faa8667a/orig"},
            {"id": 4, "type": "XBOX", "model": "Xbox Series X", "price_per_day": 800, "is_available": true, "description": "1 ТБ SSD, 4K", "image_url": "https://hatiko.ru/wa-data/public/blog/img/photo_2024-01-30_12-36-25.jpg"},
            {"id": 5, "type": "XBOX", "model": "Xbox Series S", "price_per_day": 800, "is_available": true, "description": "512 ГБ SSD", "image_url": "https://api.2droida.ru/storage/products/b11679d3f73924628580ceea19c6e9eb/5153/7a33c4eca9aca748f6753e1cb3f90101.jpg"},
            {"id": 6, "type": "XBOX", "model": "Xbox One X", "price_per_day": 500, "is_available": true, "description": "1 ТБ HDD", "image_url": "https://gameshock174.ru/upload/iblock/bbb/bbbdbb58eb867f610801394b5ef15e3a.jpg"},
        })
    })

    // ========== РЕГИСТРАЦИЯ ==========
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
        for _, u := range users {
            if u.Email == req.Email {
                c.JSON(400, gin.H{"error": "User already exists"})
                return
            }
        }

        // Хешируем пароль
        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

        // Создаём пользователя
        user := User{
            ID:       nextUserID,
            Username: req.Username,
            Email:    req.Email,
            Password: string(hashedPassword),
        }
        users = append(users, user)
        nextUserID++

        // Генерируем токен
        token, _ := generateToken(user.ID, user.Username, user.Email)

        c.JSON(201, gin.H{
            "message": "User created successfully",
            "token":   token,
            "user": gin.H{
                "id": user.ID, "username": user.Username, "email": user.Email,
            },
        })
    })

    // ========== ВХОД ==========
    r.POST("/api/login", func(c *gin.Context) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "Invalid request"})
            return
        }

        // Ищем пользователя
        var foundUser *User
        for i, u := range users {
            if u.Email == req.Email {
                foundUser = &users[i]
                break
            }
        }

        if foundUser == nil {
            c.JSON(401, gin.H{"error": "Invalid credentials"})
            return
        }

        // Проверяем пароль
        err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid credentials"})
            return
        }

        // Генерируем токен
        token, _ := generateToken(foundUser.ID, foundUser.Username, foundUser.Email)

        c.JSON(200, gin.H{
            "message": "Login successful",
            "token":   token,
            "user": gin.H{
                "id": foundUser.ID, "username": foundUser.Username, "email": foundUser.Email,
            },
        })
    })

    // ========== ПРОФИЛЬ ==========
    r.GET("/api/user/profile", authMiddleware(), func(c *gin.Context) {
        userID := c.GetInt("user_id")
        
        var foundUser *User
        for i, u := range users {
            if u.ID == userID {
                foundUser = &users[i]
                break
            }
        }

        if foundUser == nil {
            c.JSON(404, gin.H{"error": "User not found"})
            return
        }

        c.JSON(200, gin.H{
            "id": foundUser.ID, "username": foundUser.Username, "email": foundUser.Email,
            "created_at": time.Now().AddDate(-1, 0, 0),
        })
    })

    // ========== АРЕНДА ==========
    r.POST("/api/rentals", authMiddleware(), func(c *gin.Context) {
        userID := c.GetInt("user_id")
        
        var req struct {
            ConsoleID       int    `json:"console_id"`
            StartDate       string `json:"start_date"`
            EndDate         string `json:"end_date"`
            DeliveryAddress string `json:"delivery_address"`
            Phone           string `json:"phone"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "Invalid request"})
            return
        }

        // Цены консолей
        prices := map[int]float64{1: 800, 2: 800, 3: 500, 4: 800, 5: 800, 6: 500}
        price := prices[req.ConsoleID]
        if price == 0 {
            price = 800
        }

        rental := gin.H{
            "id":               nextRentalID,
            "user_id":          userID,
            "console_id":       req.ConsoleID,
            "start_date":       req.StartDate,
            "end_date":         req.EndDate,
            "total_price":      price,
            "delivery_address": req.DeliveryAddress,
            "phone":            req.Phone,
            "status":           "active",
            "created_at":       time.Now(),
        }
        rentals = append(rentals, rental)
        nextRentalID++

        c.JSON(201, gin.H{
            "rental_id":   rental["id"],
            "total_price": price,
            "message":     "Rental created successfully",
        })
    })

    r.GET("/api/my-rentals", authMiddleware(), func(c *gin.Context) {
        userID := c.GetInt("user_id")
        
        var userRentals []gin.H
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