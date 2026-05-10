package main

import (
	"console-rental/database"
	"console-rental/handlers"
	"console-rental/middleware"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
    log.Println("🚀 Starting ConsoleRent backend with PostgreSQL...")
    
    // Подключение к БД
    for i := 0; i < 30; i++ {
        err := database.InitDB()
        if err == nil {
            log.Println("✅ Database connected!")
            break
        }
        log.Printf("⚠️ DB attempt %d/30 failed: %v", i+1, err)
        if i == 29 {
            log.Fatal("❌ Could not connect to database")
        }
        time.Sleep(3 * time.Second)
    }
    
    r := gin.Default()
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
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
    
    // Публичные маршруты
    r.POST("/api/register", handlers.Register)
    r.POST("/api/login", handlers.Login)
    r.GET("/api/consoles", handlers.GetConsoles)
    r.GET("/api/consoles/:id", handlers.GetConsoleByID)
    
    // Защищенные маршруты
    auth := r.Group("/api")
    auth.Use(middleware.AuthMiddleware())
    {
        auth.GET("/user/profile", handlers.GetUserProfile)
        auth.PUT("/user/profile", handlers.UpdateUserProfile)
        auth.POST("/rentals", handlers.CreateRental)
        auth.GET("/my-rentals", handlers.GetUserRentals)
        auth.PUT("/rentals/:id/return", handlers.ReturnConsole)
    }
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("✅ Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}