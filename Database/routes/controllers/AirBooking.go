package controllers

import (
	"Go_curb/Database/components"
	"Go_curb/Database/initializers"
	"Go_curb/tableTypes"

	"github.com/gin-gonic/gin"
)

// Retrieve air bookings with server-side pagination
func GetAllItems(c *gin.Context) {

	id := c.Query("id")
	var airBookings []tableTypes.AirBooking
	query := initializers.DB.Model(&tableTypes.AirBooking{}).Where("id_invoice IS NULL AND product_type = 'flight' AND transaction_type = 'sales' AND status = 'pending'")
	embedField := c.Query("embed")
	components.Get(c, query, &airBookings, embedField, id)
}
