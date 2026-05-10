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
        time.Sleep(2 * time.Second)
    }
    
    r := gin.Default()
    
    // Получаем URL фронтенда из переменных окружения
    frontendURL := os.Getenv("FRONTEND_URL")
    if frontendURL == "" {
        frontendURL = "http://localhost:3000"
    }
    
    // CORS настройки для Vercel и Railway
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            frontendURL,
            "http://localhost:3000",
            "https://console-rental.vercel.app",
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
        // Профиль
        auth.GET("/user/profile", handlers.GetUserProfile)
        auth.PUT("/user/profile", handlers.UpdateUserProfile)
        
        // Аренды
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