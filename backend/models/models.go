package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`
    CreatedAt time.Time `json:"created_at"`
}

type Console struct {
    ID          int     `json:"id"`
    Type        string  `json:"type"`
    Model       string  `json:"model"`
    PricePerDay float64 `json:"price_per_day"`
    IsAvailable bool    `json:"is_available"`
    Description string  `json:"description"`
    ImageURL    string  `json:"image_url"`
}

type Rental struct {
    ID              int       `json:"id"`
    UserID          int       `json:"user_id"`
    ConsoleID       int       `json:"console_id"`
    StartDate       time.Time `json:"start_date"`
    EndDate         time.Time `json:"end_date"`
    TotalPrice      float64   `json:"total_price"`
    DeliveryAddress string    `json:"delivery_address"`
    Status          string    `json:"status"`
    CreatedAt       time.Time `json:"created_at"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type RentalRequest struct {
    ConsoleID       int    `json:"console_id" binding:"required"`
    StartDate       string `json:"start_date" binding:"required"`
    EndDate         string `json:"end_date" binding:"required"`
    DeliveryAddress string `json:"delivery_address" binding:"required"`
    Phone           string `json:"phone"`
    Comment         string `json:"comment"`
}

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}