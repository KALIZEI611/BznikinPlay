package handlers

import (
	"console-rental/database"
	"console-rental/middleware"
	"console-rental/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateRental(c *gin.Context) {
    var req models.RentalRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Printf("Error binding JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
        return
    }
    
    log.Printf("Rental request: console_id=%d, start=%s, end=%s, address=%s, phone=%s", 
        req.ConsoleID, req.StartDate, req.EndDate, req.DeliveryAddress, req.Phone)
    
    userID := middleware.GetUserID(c)
    if userID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - please login"})
        return
    }
    
    log.Printf("User ID: %d", userID)
    
    startDate, err := time.Parse(time.RFC3339, req.StartDate)
    if err != nil {
        log.Printf("Error parsing start date: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
        return
    }
    
    endDate, err := time.Parse(time.RFC3339, req.EndDate)
    if err != nil {
        log.Printf("Error parsing end date: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
        return
    }
    
    log.Printf("Parsed dates: start=%v, end=%v", startDate, endDate)
    
    var pricePerDay float64
    var isAvailable bool
    var consoleType string
    err = database.DB.QueryRow("SELECT price_per_day, is_available, type FROM consoles WHERE id = $1", req.ConsoleID).Scan(&pricePerDay, &isAvailable, &consoleType)
    if err != nil {
        log.Printf("Console not found: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Console not found"})
        return
    }
    
    if !isAvailable {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Console is not available for rent"})
        return
    }
    
    if endDate.Before(startDate) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "End date cannot be before start date"})
        return
    }
    
    days := endDate.Sub(startDate).Hours() / 24
    if days < 1 {
        days = 1
    }
    totalPrice := pricePerDay * days
    
    log.Printf("Calculated: %.2f * %.0f = %.2f", pricePerDay, days, totalPrice)
    
    deliveryInfo := req.DeliveryAddress
    if req.Phone != "" {
        deliveryInfo += " | Телефон: " + req.Phone
    }
    if req.Comment != "" {
        deliveryInfo += " | Комментарий: " + req.Comment
    }
    
    tx, err := database.DB.Begin()
    if err != nil {
        log.Printf("Transaction begin error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    defer tx.Rollback()
    
    var rentalID int
    err = tx.QueryRow(
        "INSERT INTO rentals (user_id, console_id, start_date, end_date, total_price, delivery_address, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
        userID, req.ConsoleID, startDate, endDate, totalPrice, deliveryInfo, "active",
    ).Scan(&rentalID)
    
    if err != nil {
        log.Printf("Error creating rental: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rental: " + err.Error()})
        return
    }
    
    _, err = tx.Exec("UPDATE consoles SET is_available = false WHERE id = $1", req.ConsoleID)
    if err != nil {
        log.Printf("Error updating console: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update console status"})
        return
    }
    
    err = tx.Commit()
    if err != nil {
        log.Printf("Transaction commit error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }
    
    log.Printf("Rental created successfully: ID=%d", rentalID)
    
    c.JSON(http.StatusCreated, gin.H{
        "rental_id": rentalID,
        "total_price": totalPrice,
        "message": "Rental created successfully",
    })
}

func GetUserRentals(c *gin.Context) {
    userID := middleware.GetUserID(c)
    if userID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    rows, err := database.DB.Query(`
        SELECT r.id, r.start_date, r.end_date, r.total_price, r.status, r.delivery_address,
               c.type, c.model, c.image_url
        FROM rentals r
        JOIN consoles c ON r.console_id = c.id
        WHERE r.user_id = $1
        ORDER BY r.created_at DESC
    `, userID)
    
    if err != nil {
        log.Printf("Error fetching rentals: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    defer rows.Close()
    
    var rentals []gin.H
    for rows.Next() {
        var id int
        var startDate, endDate time.Time
        var totalPrice float64
        var status, deliveryAddress, consoleType, model, imageURL string
        
        err := rows.Scan(&id, &startDate, &endDate, &totalPrice, &status, &deliveryAddress, &consoleType, &model, &imageURL)
        if err != nil {
            log.Printf("Error scanning row: %v", err)
            continue
        }
        
        rentals = append(rentals, gin.H{
            "id": id,
            "console": gin.H{
                "type": consoleType,
                "model": model,
                "image_url": imageURL,
            },
            "start_date": startDate,
            "end_date": endDate,
            "total_price": totalPrice,
            "delivery_address": deliveryAddress,
            "status": status,
        })
    }
    
    c.JSON(http.StatusOK, rentals)
}

func ReturnConsole(c *gin.Context) {
    rentalID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rental ID"})
        return
    }
    
    userID := middleware.GetUserID(c)
    if userID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    var dbUserID int
    var consoleID int
    var status string
    err = database.DB.QueryRow("SELECT user_id, console_id, status FROM rentals WHERE id = $1", rentalID).Scan(&dbUserID, &consoleID, &status)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Rental not found"})
        return
    }
    
    if dbUserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You can only return your own rentals"})
        return
    }
    
    if status != "active" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "This rental is already returned or cancelled"})
        return
    }
    
    tx, err := database.DB.Begin()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    defer tx.Rollback()
    
    _, err = tx.Exec("UPDATE rentals SET status = 'returned' WHERE id = $1", rentalID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rental"})
        return
    }
    
    _, err = tx.Exec("UPDATE consoles SET is_available = true WHERE id = $1", consoleID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update console"})
        return
    }
    
    err = tx.Commit()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Console returned successfully"})
}