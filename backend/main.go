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
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "message": "Server is running",
            "timestamp": time.Now().Unix(),
        })
    })
    
    // API consoles
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
    
    // Login endpoint (временный)
    r.POST("/api/login", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Login successful",
            "token": "fake-token-for-testing",
            "user": gin.H{"id": 1, "username": "test", "email": "test@example.com"},
        })
    })
    
    // Register endpoint
    r.POST("/api/register", func(c *gin.Context) {
        c.JSON(201, gin.H{
            "message": "User created successfully",
            "token": "fake-token-for-testing",
            "user": gin.H{"id": 1, "username": "newuser", "email": "new@example.com"},
        })
    })
    
    // My rentals endpoint
    r.GET("/api/my-rentals", func(c *gin.Context) {
        c.JSON(200, []gin.H{})
    })
    
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
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("✅ Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}