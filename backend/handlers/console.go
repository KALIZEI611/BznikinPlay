package handlers

import (
	"console-rental/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetConsoles(c *gin.Context) {
    rows, err := database.DB.Query("SELECT id, type, model, price_per_day, is_available, description, image_url FROM consoles")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    defer rows.Close()
    
    var consoles []gin.H
    for rows.Next() {
        var id int
        var typeVal, model, description, imageURL string
        var pricePerDay float64
        var isAvailable bool
        
        err := rows.Scan(&id, &typeVal, &model, &pricePerDay, &isAvailable, &description, &imageURL)
        if err != nil {
            continue
        }
        
        consoles = append(consoles, gin.H{
            "id": id,
            "type": typeVal,
            "model": model,
            "price_per_day": pricePerDay,
            "is_available": isAvailable,
            "description": description,
            "image_url": imageURL,
        })
    }
    
    c.JSON(http.StatusOK, consoles)
}

func GetConsoleByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    
    var console gin.H
    err := database.DB.QueryRow(
        "SELECT id, type, model, price_per_day, is_available, description, image_url FROM consoles WHERE id = $1",
        id,
    ).Scan(&console)
    
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Console not found"})
        return
    }
    
    c.JSON(http.StatusOK, console)
}