package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-super-secret-jwt-key-change-in-production")
var DB *sql.DB

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

func initDB() error {
    host := os.Getenv("PGHOST")
    port := os.Getenv("PGPORT")
    user := os.Getenv("PGUSER")
    password := os.Getenv("PGPASSWORD")
    dbname := os.Getenv("PGDATABASE")

    if host == "" {
        host = "postgres-production-096b.up.railway.app"
        port = "5432"
        user = "postgres"
        password = "pagfWwxnPVmpgnsLeXzxhyBAaWgkpmEO"
        dbname = "railway"
        log.Println("Using default database settings")
    }

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        host, port, user, password, dbname)

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }

    DB.SetMaxOpenConns(10)
    DB.SetMaxIdleConns(5)

    if err = DB.Ping(); err != nil {
        return err
    }

    log.Println("✅ Database connected!")

    // Создаём таблицы
    queries := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(100) NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
        `CREATE TABLE IF NOT EXISTS consoles (
            id SERIAL PRIMARY KEY,
            type VARCHAR(50) NOT NULL,
            model VARCHAR(100) NOT NULL,
            price_per_day DECIMAL(10,2) NOT NULL,
            is_available BOOLEAN DEFAULT TRUE,
            description TEXT,
            image_url VARCHAR(500)
        )`,
        `CREATE TABLE IF NOT EXISTS rentals (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id),
            console_id INTEGER REFERENCES consoles(id),
            start_date TIMESTAMP NOT NULL,
            end_date TIMESTAMP NOT NULL,
            total_price DECIMAL(10,2) NOT NULL,
            delivery_address TEXT NOT NULL,
            status VARCHAR(50) DEFAULT 'active',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
    }

    for _, q := range queries {
        if _, err := DB.Exec(q); err != nil {
            log.Printf("Table creation warning: %v", err)
        }
    }

    // Вставляем тестовые консоли если их нет
    var count int
    DB.QueryRow("SELECT COUNT(*) FROM consoles").Scan(&count)
    if count == 0 {
        consoles := []string{
            `INSERT INTO consoles (type, model, price_per_day, description, image_url) VALUES 
            ('PS5', 'PlayStation 5 Standard', 800, 'Полная комплектация', 'https://avatars.mds.yandex.net/get-mpic/13230222/2a000001969fb44f0b2b0473bcfe73eb4de4/orig')`,
            `INSERT INTO consoles (type, model, price_per_day, description, image_url) VALUES 
            ('PS5', 'PlayStation 5 Digital', 800, 'Цифровая версия', 'https://avatars.mds.yandex.net/i?id=0b869e3e8145ca09fba9fa1e77702f95_l-4355007-images-thumbs&n=13')`,
            `INSERT INTO consoles (type, model, price_per_day, description, image_url) VALUES 
            ('PS4', 'PlayStation 4 Slim', 500, '500 GB HDD', 'https://avatars.mds.yandex.net/get-mpic/5173149/2a0000019180dbf1814ebb7ae678faa8667a/orig')`,
            `INSERT INTO consoles (type, model, price_per_day, description, image_url) VALUES 
            ('XBOX', 'Xbox Series X', 800, '1 ТБ SSD', 'https://hatiko.ru/wa-data/public/blog/img/photo_2024-01-30_12-36-25.jpg')`,
            `INSERT INTO consoles (type, model, price_per_day, description, image_url) VALUES 
            ('XBOX', 'Xbox Series S', 800, '512 ГБ SSD', 'https://api.2droida.ru/storage/products/b11679d3f73924628580ceea19c6e9eb/5153/7a33c4eca9aca748f6753e1cb3f90101.jpg')`,
            `INSERT INTO consoles (type, model, price_per_day, description, image_url) VALUES 
            ('XBOX', 'Xbox One X', 500, '1 ТБ HDD', 'https://gameshock174.ru/upload/iblock/bbb/bbbdbb58eb867f610801394b5ef15e3a.jpg')`,
        }
        for _, c := range consoles {
            DB.Exec(c)
        }
        log.Println("Sample consoles inserted")
    }

    return nil
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
    if err := initDB(); err != nil {
        log.Printf("Database connection warning: %v", err)
    }

    r := gin.Default()

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

    // Каталог консолей (из БД)
    r.GET("/api/consoles", func(c *gin.Context) {
        rows, err := DB.Query("SELECT id, type, model, price_per_day, is_available, description, image_url FROM consoles")
        if err != nil {
            c.JSON(500, gin.H{"error": "Database error"})
            return
        }
        defer rows.Close()

        var consoles []gin.H
        for rows.Next() {
            var id int
            var typeVal, model, description, imageURL string
            var pricePerDay float64
            var isAvailable bool
            rows.Scan(&id, &typeVal, &model, &pricePerDay, &isAvailable, &description, &imageURL)
            consoles = append(consoles, gin.H{
                "id": id, "type": typeVal, "model": model,
                "price_per_day": pricePerDay, "is_available": isAvailable,
                "description": description, "image_url": imageURL,
            })
        }
        c.JSON(200, consoles)
    })

    // РЕГИСТРАЦИЯ
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

        // Проверяем существует ли пользователь
        var exists bool
        DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
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
        err = DB.QueryRow(
            "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
            req.Username, req.Email, string(hashedPassword),
        ).Scan(&userID)

        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to create user"})
            return
        }

        // Генерируем токен
        token, err := generateToken(userID, req.Username, req.Email)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to generate token"})
            return
        }

        c.JSON(201, gin.H{
            "message": "User created successfully",
            "token":   token,
            "user": gin.H{
                "id": userID, "username": req.Username, "email": req.Email,
            },
        })
    })

    // ВХОД
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
        err := DB.QueryRow("SELECT id, username, password FROM users WHERE email = $1", req.Email).Scan(&userID, &username, &hashedPassword)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid email or password"})
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
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
            "token":   token,
            "user": gin.H{
                "id": userID, "username": username, "email": req.Email,
            },
        })
    })

    // ПРОФИЛЬ (защищённый)
    r.GET("/api/user/profile", authMiddleware(), func(c *gin.Context) {
        userID := c.GetInt("user_id")
        var username, email string
        var createdAt time.Time
        err := DB.QueryRow("SELECT username, email, created_at FROM users WHERE id = $1", userID).Scan(&username, &email, &createdAt)
        if err != nil {
            c.JSON(404, gin.H{"error": "User not found"})
            return
        }
        c.JSON(200, gin.H{
            "id": userID, "username": username, "email": email, "created_at": createdAt,
        })
    })

    // АРЕНДА (защищённая)
    var rentals = []gin.H{}
    var nextID = 1

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

        var price float64
        DB.QueryRow("SELECT price_per_day FROM consoles WHERE id = $1", req.ConsoleID).Scan(&price)

        rental := gin.H{
            "id":               nextID,
            "user_id":          userID,
            "total_price":      price,
            "delivery_address": req.DeliveryAddress,
            "status":           "active",
        }
        rentals = append(rentals, rental)
        nextID++

        c.JSON(201, gin.H{"rental_id": rental["id"], "total_price": price, "message": "Rental created"})
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
        c.JSON(200, gin.H{"message": "Console returned"})
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("✅ Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}