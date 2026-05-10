package handlers

import (
	"console-rental/database"
	"console-rental/models"
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
    
    var consoles []models.Console
    for rows.Next() {
        var console models.Console
        err := rows.Scan(&console.ID, &console.Type, &console.Model, &console.PricePerDay, &console.IsAvailable, &console.Description, &console.ImageURL)
        if err != nil {
            continue
        }
        consoles = append(consoles, console)
    }
    
    c.JSON(http.StatusOK, consoles)
}

func GetConsoleByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    
    var console models.Console
    err := database.DB.QueryRow(
        "SELECT id, type, model, price_per_day, is_available, description, image_url FROM consoles WHERE id = $1",
        id,
    ).Scan(&console.ID, &console.Type, &console.Model, &console.PricePerDay, &console.IsAvailable, &console.Description, &console.ImageURL)
    
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Console not found"})
        return
    }
    
    c.JSON(http.StatusOK, console)
}