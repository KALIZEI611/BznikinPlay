package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
        c.JSON(200, gin.H{
            "status": "ok",
            "message": "Server is running",
            "timestamp": time.Now().Unix(),
        })
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
    
    // ========== АВТОРИЗАЦИЯ ==========
    r.POST("/api/login", func(c *gin.Context) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        c.ShouldBindJSON(&req)
        
        token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20ifQ.test"
        
        c.JSON(200, gin.H{
            "message": "Login successful",
            "token": token,
            "user": gin.H{"id": 1, "username": "user", "email": req.Email},
        })
    })
    
    r.POST("/api/register", func(c *gin.Context) {
        var req struct {
            Username string `json:"username"`
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        c.ShouldBindJSON(&req)
        
        token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20ifQ.test"
        
        c.JSON(201, gin.H{
            "message": "User created successfully",
            "token": token,
            "user": gin.H{"id": 1, "username": req.Username, "email": req.Email},
        })
    })
    
    // ========== ПРОФИЛЬ ==========
    r.GET("/api/user/profile", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "id": 1,
            "username": "testuser",
            "email": "test@example.com",
            "created_at": time.Now().AddDate(-1, 0, 0).Format(time.RFC3339),
        })
    })
    
    r.PUT("/api/user/profile", func(c *gin.Context) {
        var req struct {
            Username        string `json:"username"`
            Email           string `json:"email"`
            CurrentPassword string `json:"current_password"`
            NewPassword     string `json:"new_password"`
        }
        c.ShouldBindJSON(&req)
        
        c.JSON(200, gin.H{
            "message": "Profile updated successfully",
            "token": "new-token",
            "user": gin.H{"id": 1, "username": req.Username, "email": req.Email},
        })
    })
    
    // ========== АРЕНДА ==========
    // Временное хранилище аренд в памяти
    var rentals = []gin.H{}
    var nextId = 1
    
    r.POST("/api/rentals", func(c *gin.Context) {
        var req struct {
            ConsoleID       int    `json:"console_id"`
            StartDate       string `json:"start_date"`
            EndDate         string `json:"end_date"`
            DeliveryAddress string `json:"delivery_address"`
            Phone           string `json:"phone"`
            Comment         string `json:"comment"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "Invalid request data"})
            return
        }
        
        // Находим консоль
        var consolePrice float64 = 800
        var consoleModel string = "Console"
        for _, c := range []gin.H{
            {"id": 1, "price_per_day": 800, "model": "PlayStation 5 Standard"},
            {"id": 2, "price_per_day": 800, "model": "PlayStation 5 Digital"},
            {"id": 3, "price_per_day": 500, "model": "PlayStation 4 Slim"},
            {"id": 4, "price_per_day": 800, "model": "Xbox Series X"},
            {"id": 5, "price_per_day": 800, "model": "Xbox Series S"},
            {"id": 6, "price_per_day": 500, "model": "Xbox One X"},
        } {
            if c["id"] == req.ConsoleID {
                consolePrice = c["price_per_day"].(float64)
                consoleModel = c["model"].(string)
                break
            }
        }
        
        // Рассчитываем стоимость
        totalPrice := consolePrice
        
        rental := gin.H{
            "id":               nextId,
            "console_id":       req.ConsoleID,
            "console_model":    consoleModel,
            "start_date":       req.StartDate,
            "end_date":         req.EndDate,
            "total_price":      totalPrice,
            "delivery_address": req.DeliveryAddress,
            "phone":            req.Phone,
            "comment":          req.Comment,
            "status":           "active",
            "created_at":       time.Now().Format(time.RFC3339),
        }
        
        rentals = append(rentals, rental)
        nextId++
        
        c.JSON(201, gin.H{
            "rental_id":   rental["id"],
            "total_price": totalPrice,
            "message":     "Rental created successfully",
        })
    })
    
    r.GET("/api/my-rentals", func(c *gin.Context) {
        // Возвращаем аренды текущего пользователя
        userRentals := []gin.H{}
        for _, r := range rentals {
            userRentals = append(userRentals, gin.H{
                "id": r["id"],
                "console": gin.H{
                    "type":       "PS5",
                    "model":      r["console_model"],
                    "image_url":  "https://avatars.mds.yandex.net/get-mpic/13230222/2a000001969fb44f0b2b0473bcfe73eb4de4/orig",
                },
                "start_date":        r["start_date"],
                "end_date":          r["end_date"],
                "total_price":       r["total_price"],
                "delivery_address":  r["delivery_address"],
                "status":            r["status"],
            })
        }
        c.JSON(200, userRentals)
    })
    
    r.PUT("/api/rentals/:id/return", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Console returned successfully"})
    })
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("✅ Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}