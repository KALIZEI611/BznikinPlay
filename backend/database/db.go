package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
    // Приоритет у переменных Railway
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    
    // Если переменные Railway не заданы, используем локальные
    if host == "" {
        host = os.Getenv("PGHOST")
        if host == "" {
            host = "localhost"
        }
    }
    if port == "" {
        port = os.Getenv("PGPORT")
        if port == "" {
            port = "5432"
        }
    }
    if user == "" {
        user = os.Getenv("PGUSER")
        if user == "" {
            user = "rental_user"
        }
    }
    if password == "" {
        password = os.Getenv("PGPASSWORD")
        if password == "" {
            password = "rental_pass"
        }
    }
    if dbname == "" {
        dbname = os.Getenv("PGDATABASE")
        if dbname == "" {
            dbname = "console_rental"
        }
    }
    
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        host, port, user, password, dbname)
    
    log.Printf("Connecting to database with connection string: host=%s port=%s dbname=%s", host, port, dbname)
    
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("error opening database: %v", err)
    }
    
    err = DB.Ping()
    if err != nil {
        return fmt.Errorf("database ping failed: %v", err)
    }
    
    createTables()
    insertSampleData()
    return nil
}

func createTables() {
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
    
    for _, query := range queries {
        _, err := DB.Exec(query)
        if err != nil {
            log.Printf("Error creating table: %v", err)
        }
    }
}

func insertSampleData() {
    var count int
    DB.QueryRow("SELECT COUNT(*) FROM consoles").Scan(&count)
    
    if count == 0 {
        consoles := []struct {
            typeVal     string
            model       string
            pricePerDay float64
            description string
            imageURL    string
        }{
            {"PS5", "PlayStation 5 Standard", 800, "Полная комплектация: 2 контроллера DualSense, зарядная станция, подписка PlayStation Plus на 1 месяц", "https://avatars.mds.yandex.net/get-mpic/13230222/2a000001969fb44f0b2b0473bcfe73eb4de4/orig"},
            {"PS5", "PlayStation 5 Digital", 800, "Цифровая версия без дисковода, 1 контроллер DualSense, 825 GB SSD, подписка PlayStation Plus Essential", "https://avatars.mds.yandex.net/i?id=0b869e3e8145ca09fba9fa1e77702f95_l-4355007-images-thumbs&n=13"},
            {"PS4", "PlayStation 4 Slim", 500, "500 GB HDD, 1 контроллер DualShock 4, компактный дизайн, подборка лучших игр", "https://avatars.mds.yandex.net/get-mpic/5173149/2a0000019180dbf1814ebb7ae678faa8667a/orig"},
            {"XBOX", "Xbox Series X", 800, "1 ТБ SSD, 4K Gaming, контроллер Xbox Wireless, Game Pass Ultimate на 1 месяц", "https://hatiko.ru/wa-data/public/blog/img/photo_2024-01-30_12-36-25.jpg"},
            {"XBOX", "Xbox Series S", 800, "512 ГБ SSD, цифровая версия, компактный дизайн, Game Pass Ultimate на 1 месяц", "https://api.2droida.ru/storage/products/b11679d3f73924628580ceea19c6e9eb/5153/7a33c4eca9aca748f6753e1cb3f90101.jpg"},
            {"XBOX", "Xbox One X", 500, "1 ТБ HDD, 4K Gaming, поддержка HDR, контроллер Xbox Wireless", "https://gameshock174.ru/upload/iblock/bbb/bbbdbb58eb867f610801394b5ef15e3a.jpg"},
        }
        
        for _, c := range consoles {
            _, err := DB.Exec(
                "INSERT INTO consoles (type, model, price_per_day, description, image_url) VALUES ($1, $2, $3, $4, $5)",
                c.typeVal, c.model, c.pricePerDay, c.description, c.imageURL,
            )
            if err != nil {
                log.Println("Error inserting console:", err)
            }
        }
        log.Println("Sample consoles inserted successfully")
    }
}