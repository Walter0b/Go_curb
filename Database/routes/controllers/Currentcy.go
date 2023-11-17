package controllers

import (
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Retrieve all Id from countries
func GetAllCurrency(c *gin.Context) {
	var currencies []tableTypes.Currency
	if err := initializers.DB.Find(&currencies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, currencies)
}
