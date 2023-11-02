package routes

import (
    "Go_curb/tableTypes"
    "net/http"
    "net/url"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func CustomerRoutes(r *gin.Engine, db *gorm.DB) {
    // GET / - Retrieve all customers
	// http://your-api-domain/
    r.GET("/", func(c *gin.Context) {
        var customers []tableTypes.Customer
        if err := db.Find(&customers).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, customers)
    })

    // GET /customer - Retrieve a specific customer by ID
	//  http://your-api-domain/customer?id=1

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
	 // GET / - Retrieve all Id from currencys
	// http://your-api-domain/
	r.GET("/Currencies", func(c *gin.Context) {
		var currencies []tableTypes.Currency
		if err := db.Find(&currencies).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, currencies)
	})
	
 // GET / - Retrieve all Id from countries
	// http://your-api-domain/
	r.GET("/Countries", func(c *gin.Context) {
		var Countries []tableTypes.Country
		if err := db.Find(&Countries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, Countries)
	})
	
    // PUT /customer - Update a specific customer by ID
	// http://your-api-domain/customer?id=1
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
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := db.Save(&customerID).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, customerID)
    })

    // POST /customers - Create a new customer
	// http://your-api-domain/customers
    r.POST("/customers", func(c *gin.Context) {
        var customer tableTypes.Customer
        if err := c.ShouldBindJSON(&customer); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := db.Create(&customer).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, customer)
    })
    
    // DELETE /customers/:id - Delete a specific customer by ID
	// http://your-api-domain/customers/1
    r.DELETE("/customers/:id", func(c *gin.Context) {
        customerID := c.Param("id")

        if err := db.Where("id = ?", customerID).Delete(&tableTypes.Customer{}).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    
        c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
    })
}
