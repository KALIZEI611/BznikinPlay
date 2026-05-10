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
    // Инициализация базы данных с повторными попытками
    for i := 0; i < 30; i++ {
        err := database.InitDB()
        if err == nil {
            log.Println("Successfully connected to database")
            break
        }
        log.Printf("Failed to connect to database (attempt %d/30): %v", i+1, err)
        if i == 29 {
            log.Fatal("Could not connect to database after 30 attempts")
        }
        time.Sleep(3 * time.Second)
    }
    
    r := gin.Default()
    
    // Health check endpoint for Railway
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "message": "Server is running",
            "timestamp": time.Now().Unix(),
        })
    })
    
    // CORS настройки
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            "http://localhost:3000",
            "http://localhost:5173",
            "https://bznikin-play.vercel.app",
            "https://*.vercel.app",
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
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
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}