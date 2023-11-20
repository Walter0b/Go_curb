package routes

import (
	"Go_curb/Database/routes/controllers"

	"github.com/gin-gonic/gin"
)

func Api(r *gin.Engine) {
	//--------------------- customer  ---------------------------||

	r.GET("/customers", controllers.GetAllCustomer)
	r.POST("/customers", controllers.CreateCustomer)
	r.GET("/customers/:id", controllers.GetSpecificCustomer)
	r.PUT("/customers/:id", controllers.UpdateCustomer)
	r.DELETE("/customers/:id", controllers.DeleteCutomer)

	//--------------------- customer ---------------------------||

	//--------------------- Invoice ----------------------------||

	r.GET("/invoices", controllers.GetAllInvoices)
	r.GET("/invoices/customers/:id", controllers.GetAllInvoices)
	r.POST("/invoices", controllers.CreateInvoice)
	r.GET("/invoices/:id", controllers.GetSpecificInvoice)
	// r.PUT("/customers/:id", controllers.UpdatedInvoice)
	r.DELETE("/invoices/:id", controllers.DeleteInvoices)

	//--------------------- Invoice ----------------------------||

	//--------------------- Payments ---------------------------||

	r.GET("/payments", controllers.GetAllCustomer)
	r.POST("/payments", controllers.CreatePayments)
	r.GET("/payments/:id", controllers.GetSpecifiPayments)
	// r.PUT("/payments/:id", controllers.UpdateCustomer)
	r.DELETE("/payments/:id", controllers.DeletePayments)

	//--------------------- Payments ---------------------------||

	//--------------------- Currency ---------------------------||

	r.GET("/currencies", controllers.GetAllCurrency)

	//--------------------- Currency ---------------------------||

	//-------------------- AirBooking --------------------------||

	r.GET("/airbooking", controllers.GetAllItems)

	//-------------------- AirBooking --------------------------||

	r.GET("/patch", controllers.Imputations)
	
}
