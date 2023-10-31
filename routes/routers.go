package routes

import (
	"Go_curb/tableTypes"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CustomerRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/", func(c *gin.Context) {
		var customers []tableTypes.Customer
		if err := db.Find(&customers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, customers)
	})
	r.GET("/customer", func(c *gin.Context) {
		Host := c.Request.URL.RequestURI()
		myUrl, _ := url.Parse(Host)
		params, _ := url.ParseQuery(myUrl.RawQuery)

		var IdValue string
		for key := range params {
			IdValue = key
		}
		id, _ := strconv.Atoi(IdValue)

		var customerID = tableTypes.Customer{ID: id}
		if err := db.First(&customerID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, customerID)
	})

	r.PUT("/customer", func(c *gin.Context) {

		Host := c.Request.URL.RequestURI()
		myUrl, _ := url.Parse(Host)
		params, _ := url.ParseQuery(myUrl.RawQuery)

		var IdValue string
		for key := range params {
			IdValue = key
		}
		id, _ := strconv.Atoi(IdValue)

		var customerID = tableTypes.Customer{ID: id}
		if err := db.First(&customerID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := c.ShouldBindJSON(&customerID); err != nil {
			//log.Println(customerID)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Save(&customerID).Error; err != nil {
			//log.Println(customerID)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, customerID)
	})

	r.POST("/customers", func(c *gin.Context) {
		var customer tableTypes.Customer
		if err := c.ShouldBindJSON(&customer); err != nil {
			log.Println(customer)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, customer)
	})

}
