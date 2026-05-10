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
    }
    
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        host, port, user, password, dbname)
    
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("error opening database: %v", err)
    }
    
    DB.SetMaxOpenConns(25)
    DB.SetMaxIdleConns(5)
    
    err = DB.Ping()
    if err != nil {
        return fmt.Errorf("database ping failed: %v", err)
    }
    
    log.Println("✅ Database connected successfully!")
    createTables()
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
        if _, err := DB.Exec(query); err != nil {
            log.Printf("⚠️ Error creating table: %v", err)
        }
    }
    log.Println("✅ Tables ready")
}