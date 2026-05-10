package main

import (
	"log"
	"os"
	"strconv"
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

var users = []User{}
var nextUserID = 1

var rentals = []gin.H{}
var nextRentalID = 1

func init() {
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
    users = append(users, User{
        ID:       nextUserID,
        Username: "testuser",
        Email:    "test@example.com",
        Password: string(hashedPassword),
    })
    nextUserID++
    log.Println("✅ Тестовый пользователь создан: test@example.com / 123456")
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
            c.JSON(401, gin.H{"error": "❌ Ошибка авторизации: токен не предоставлен"})
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "❌ Ошибка авторизации: неверный формат токена"})
            c.Abort()
            return
        }

        tokenString := parts[1]
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "❌ Ошибка авторизации: токен истёк или недействителен"})
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
    
    log.Println("🚀 Запуск ConsoleRent backend...")

    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            "https://bznikin-play.vercel.app",
            "http://localhost:3000",
            "http://localhost:5173",
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
    })

    r.GET("/api/consoles", func(c *gin.Context) {
        c.JSON(200, []gin.H{
            {"id": 1, "type": "PS5", "model": "PlayStation 5 Standard", "price_per_day": 800, "is_available": true, "description": "Полная комплектация", "image_url": "https://avatars.mds.yandex.net/get-mpic/13230222/2a000001969fb44f0b2b0473bcfe73eb4de4/orig"},
            {"id": 2, "type": "PS5", "model": "PlayStation 5 Digital", "price_per_day": 800, "is_available": true, "description": "Цифровая версия", "image_url": "https://avatars.mds.yandex.net/i?id=0b869e3e8145ca09fba9fa1e77702f95_l-4355007-images-thumbs&n=13"},
            {"id": 3, "type": "PS4", "model": "PlayStation 4 Slim", "price_per_day": 500, "is_available": true, "description": "500 GB HDD", "image_url": "https://avatars.mds.yandex.net/get-mpic/5173149/2a0000019180dbf1814ebb7ae678faa8667a/orig"},
            {"id": 4, "type": "XBOX", "model": "Xbox Series X", "price_per_day": 800, "is_available": true, "description": "1 ТБ SSD, 4K", "image_url": "https://avatars.mds.yandex.net/get-mpic/4956986/2a0000018da2f1bf854292750979fe770113/orig"},
            {"id": 5, "type": "XBOX", "model": "Xbox Series S", "price_per_day": 800, "is_available": true, "description": "512 ГБ SSD", "image_url": "https://api.2droida.ru/storage/products/b11679d3f73924628580ceea19c6e9eb/5153/7a33c4eca9aca748f6753e1cb3f90101.jpg"},
            {"id": 6, "type": "XBOX", "model": "Xbox One X", "price_per_day": 500, "is_available": true, "description": "1 ТБ HDD", "image_url": "https://gameshock174.ru/upload/iblock/bbb/bbbdbb58eb867f610801394b5ef15e3a.jpg"},
        })
    })

    r.POST("/api/register", func(c *gin.Context) {
        var req struct {
            Username string `json:"username"`
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "❌ Ошибка: неверный формат данных"})
            return
        }

        if req.Username == "" {
            c.JSON(400, gin.H{"error": "❌ Ошибка: имя пользователя не может быть пустым"})
            return
        }
        
        if req.Email == "" {
            c.JSON(400, gin.H{"error": "❌ Ошибка: email не может быть пустым"})
            return
        }
        
        if req.Password == "" {
            c.JSON(400, gin.H{"error": "❌ Ошибка: пароль не может быть пустым"})
            return
        }
        
        if len(req.Password) < 6 {
            c.JSON(400, gin.H{"error": "❌ Ошибка: пароль должен содержать минимум 6 символов"})
            return
        }

        for _, u := range users {
            if u.Email == req.Email {
                c.JSON(400, gin.H{"error": "❌ Ошибка: пользователь с таким email уже существует"})
                return
            }
        }

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(500, gin.H{"error": "❌ Ошибка сервера: не удалось зашифровать пароль"})
            return
        }

        user := User{
            ID:       nextUserID,
            Username: req.Username,
            Email:    req.Email,
            Password: string(hashedPassword),
        }
        users = append(users, user)
        nextUserID++

        token, _ := generateToken(user.ID, user.Username, user.Email)

        c.JSON(201, gin.H{
            "message": "✅ Регистрация прошла успешно! Добро пожаловать!",
            "token":   token,
            "user": gin.H{
                "id": user.ID, 
                "username": user.Username, 
                "email": user.Email,
            },
        })
    })

    r.POST("/api/login", func(c *gin.Context) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "❌ Ошибка: неверный формат данных"})
            return
        }

        if req.Email == "" {
            c.JSON(400, gin.H{"error": "❌ Ошибка: введите email"})
            return
        }
        
        if req.Password == "" {
            c.JSON(400, gin.H{"error": "❌ Ошибка: введите пароль"})
            return
        }

        var foundUser *User
        for i, u := range users {
            if u.Email == req.Email {
                foundUser = &users[i]
                break
            }
        }

        if foundUser == nil {
            c.JSON(401, gin.H{"error": "❌ Ошибка: пользователь с таким email не найден"})
            return
        }

        err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
        if err != nil {
            c.JSON(401, gin.H{"error": "❌ Ошибка: неверный пароль"})
            return
        }

        token, _ := generateToken(foundUser.ID, foundUser.Username, foundUser.Email)

        c.JSON(200, gin.H{
            "message": "✅ Вход выполнен успешно! Добро пожаловать, " + foundUser.Username + "!",
            "token":   token,
            "user": gin.H{
                "id": foundUser.ID, 
                "username": foundUser.Username, 
                "email": foundUser.Email,
            },
        })
    })

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
            c.JSON(404, gin.H{"error": "❌ Ошибка: пользователь не найден"})
            return
        }

        c.JSON(200, gin.H{
            "id": foundUser.ID, 
            "username": foundUser.Username, 
            "email": foundUser.Email,
            "created_at": time.Now().AddDate(-1, 0, 0),
        })
    })

    r.PUT("/api/user/profile", authMiddleware(), func(c *gin.Context) {
        userID := c.GetInt("user_id")
        
        var req struct {
            Username        string `json:"username"`
            Email           string `json:"email"`
            CurrentPassword string `json:"current_password"`
            NewPassword     string `json:"new_password"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": "❌ Ошибка: неверный формат данных"})
            return
        }

        var foundUser *User
        var userIndex int
        for i, u := range users {
            if u.ID == userID {
                foundUser = &users[i]
                userIndex = i
                break
            }
        }

        if foundUser == nil {
            c.JSON(404, gin.H{"error": "❌ Ошибка: пользователь не найден"})
            return
        }

        err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.CurrentPassword))
        if err != nil {
            c.JSON(401, gin.H{"error": "❌ Ошибка: текущий пароль введён неверно"})
            return
        }

        changes := []string{}

        if req.Username != "" && req.Username != foundUser.Username {
            foundUser.Username = req.Username
            changes = append(changes, "имя пользователя")
        }
        
        if req.Email != "" && req.Email != foundUser.Email {
            for i, u := range users {
                if u.Email == req.Email && i != userIndex {
                    c.JSON(400, gin.H{"error": "❌ Ошибка: этот email уже используется другим пользователем"})
                    return
                }
            }
            foundUser.Email = req.Email
            changes = append(changes, "email")
        }
        
        if req.NewPassword != "" {
            if len(req.NewPassword) < 6 {
                c.JSON(400, gin.H{"error": "❌ Ошибка: новый пароль должен содержать минимум 6 символов"})
                return
            }
            hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
            foundUser.Password = string(hashedPassword)
            changes = append(changes, "пароль")
        }

        message := "✅ Профиль успешно обновлён"
        if len(changes) > 0 {
            message += " (изменено: " + strings.Join(changes, ", ") + ")"
        }

        token, _ := generateToken(foundUser.ID, foundUser.Username, foundUser.Email)

        c.JSON(200, gin.H{
            "message": message,
            "token":   token,
            "user": gin.H{
                "id": foundUser.ID, 
                "username": foundUser.Username, 
                "email": foundUser.Email,
            },
        })
    })

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
            c.JSON(400, gin.H{"error": "❌ Ошибка: неверный формат данных"})
            return
        }

        if req.DeliveryAddress == "" {
            c.JSON(400, gin.H{"error": "❌ Ошибка: укажите адрес доставки"})
            return
        }
        
        if req.Phone == "" {
            c.JSON(400, gin.H{"error": "❌ Ошибка: укажите номер телефона"})
            return
        }

        prices := map[int]float64{1: 800, 2: 800, 3: 500, 4: 800, 5: 800, 6: 500}
        price := prices[req.ConsoleID]
        if price == 0 {
            price = 800
        }

        images := map[int]string{
            1: "https://avatars.mds.yandex.net/get-mpic/13230222/2a000001969fb44f0b2b0473bcfe73eb4de4/orig",
            2: "https://avatars.mds.yandex.net/i?id=0b869e3e8145ca09fba9fa1e77702f95_l-4355007-images-thumbs&n=13",
            3: "https://avatars.mds.yandex.net/get-mpic/5173149/2a0000019180dbf1814ebb7ae678faa8667a/orig",
            4: "https://avatars.mds.yandex.net/get-mpic/4956986/2a0000018da2f1bf854292750979fe770113/orig",
            5: "https://api.2droida.ru/storage/products/b11679d3f73924628580ceea19c6e9eb/5153/7a33c4eca9aca748f6753e1cb3f90101.jpg",
            6: "https://gameshock174.ru/upload/iblock/bbb/bbbdbb58eb867f610801394b5ef15e3a.jpg",
        }

        models := map[int]string{
            1: "PlayStation 5 Standard",
            2: "PlayStation 5 Digital",
            3: "PlayStation 4 Slim",
            4: "Xbox Series X",
            5: "Xbox Series S",
            6: "Xbox One X",
        }

        types := map[int]string{
            1: "PS5", 2: "PS5", 3: "PS4", 4: "XBOX", 5: "XBOX", 6: "XBOX",
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
            "console": gin.H{
                "type":      types[req.ConsoleID],
                "model":     models[req.ConsoleID],
                "image_url": images[req.ConsoleID],
            },
        }
        rentals = append(rentals, rental)
        nextRentalID++

        days := 1
        if req.StartDate != "" && req.EndDate != "" {
            start, _ := time.Parse(time.RFC3339, req.StartDate)
            end, _ := time.Parse(time.RFC3339, req.EndDate)
            days = int(end.Sub(start).Hours()/24) + 1
            if days < 1 {
                days = 1
            }
        }
        totalPrice := price * float64(days)

        c.JSON(201, gin.H{
            "rental_id":   rental["id"],
            "total_price": totalPrice,
            "message":     "✅ Аренда успешно оформлена! Консоль будет доставлена по указанному адресу.",
        })
    })

    r.GET("/api/my-rentals", authMiddleware(), func(c *gin.Context) {
        userID := c.GetInt("user_id")
        
        var userRentals []gin.H
        for _, r := range rentals {
            if r["user_id"] == userID {
                status := r["status"].(string)
                statusText := "Активна"
                if status == "returned" {
                    statusText = "Возвращена"
                } else if status == "cancelled" {
                    statusText = "Отменена"
                }
                
                userRentals = append(userRentals, gin.H{
                    "id":               r["id"],
                    "console":          r["console"],
                    "start_date":       r["start_date"],
                    "end_date":         r["end_date"],
                    "total_price":      r["total_price"],
                    "delivery_address": r["delivery_address"],
                    "status":           status,
                    "status_text":      statusText,
                })
            }
        }
        c.JSON(200, userRentals)
    })

    r.PUT("/api/rentals/:id/return", authMiddleware(), func(c *gin.Context) {
        rentalID, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(400, gin.H{"error": "❌ Ошибка: неверный идентификатор аренды"})
            return
        }

        userID := c.GetInt("user_id")
        
        var rentalIndex = -1
        for i, r := range rentals {
            if r["id"] == rentalID {
                rentalIndex = i
                break
            }
        }

        if rentalIndex == -1 {
            c.JSON(404, gin.H{"error": "❌ Ошибка: аренда не найдена"})
            return
        }

        if rentals[rentalIndex]["user_id"] != userID {
            c.JSON(403, gin.H{"error": "❌ Ошибка: вы можете вернуть только свою аренду"})
            return
        }

        if rentals[rentalIndex]["status"] != "active" {
            c.JSON(400, gin.H{"error": "❌ Ошибка: эта аренда уже возвращена"})
            return
        }

        rentals[rentalIndex]["status"] = "returned"

        c.JSON(200, gin.H{"message": "✅ Консоль успешно возвращена! Спасибо за аренду!"})
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("✅ Сервер запущен на порту %s", port)
    log.Printf("📊 Пользователей: %d", len(users))
    log.Printf("📊 Активных аренд: %d", len(rentals))
    log.Fatal(r.Run(":" + port))
}