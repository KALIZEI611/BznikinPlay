package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
    // Пытаемся подключиться к БД, но не падаем, если не получилось
    dbErr := initDB()
    if dbErr != nil {
        log.Printf("⚠️ WARNING: Could not connect to database: %v", dbErr)
        log.Println("✅ Backend will continue to work with fake data (DB is optional)")
    } else {
        log.Println("🎉 Database connection is ready!")
    }

    r := gin.Default()

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

    // 1. Health check (для Railway)
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "message": "Server is running with DB support",
            "timestamp": time.Now().Unix(),
            "db_connected": dbErr == nil,
        })
    })

    // 2. НОВЫЙ ЭНДПОИНТ: проверка статуса БД
    r.GET("/api/db-status", func(c *gin.Context) {
        if db == nil {
            c.JSON(500, gin.H{"status": "disconnected", "error": "Database not initialized"})
            return
        }
        err := db.Ping()
        if err != nil {
            c.JSON(500, gin.H{"status": "disconnected", "error": err.Error()})
            return
        }
        c.JSON(200, gin.H{"status": "connected", "message": "PostgreSQL is alive!"})
    })

    // 3. API consoles (работает с фейковыми данными, как и раньше)
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

    // 4. Остальные эндпоинты (рабочие заглушки)
    // Login endpoint (исправленный)
r.POST("/api/login", func(c *gin.Context) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    // Создаём фейковый JWT токен для тестирования
    token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20ifQ.test"
    
    c.JSON(200, gin.H{
        "message": "Login successful",
        "token": token,
        "user": gin.H{
            "id": 1,
            "username": req.Email[:3],
            "email": req.Email,
        },
    })
})

// Register endpoint (исправленный)
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
    
    token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20ifQ.test"
    
    c.JSON(201, gin.H{
        "message": "User created successfully",
        "token": token,
        "user": gin.H{
            "id": 1,
            "username": req.Username,
            "email": req.Email,
        },
    })
})

    r.GET("/api/my-rentals", func(c *gin.Context) {
        c.JSON(200, []gin.H{})
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("✅ Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}