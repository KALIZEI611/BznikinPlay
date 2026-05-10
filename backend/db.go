package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

var db *sql.DB

func initDB() error {
    // Берём переменные, которые вы уже установили в Railway
    host := os.Getenv("PGHOST")
    port := os.Getenv("PGPORT")
    user := os.Getenv("PGUSER")
    password := os.Getenv("PGPASSWORD")
    dbname := os.Getenv("PGDATABASE")

    if host == "" {
        // Если нет PGHOST, используем то, что вы видели в интерфейсе
        host = "postgres-production-096b.up.railway.app"
        port = "5432"
        user = "postgres"
        password = "pagfWwxnPVmpgnsLeXzxhyBAaWgkpmEO"
        dbname = "railway"
    }

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        host, port, user, password, dbname)

    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }

    err = db.Ping()
    if err != nil {
        return err
    }

    log.Println("✅ SUCCESS: Connected to PostgreSQL database!")
    return nil
}

