package routes

import (
	"Go_curb/Database/routes/controllers"

	"github.com/gin-gonic/gin"
)

func Api(r *gin.Engine) {
	//--------------------- customer  ---------------------------||

	r.POST("/customers", controllers.CreateCustomer)
	r.GET("/customers", controllers.GetAllCustomer)
	// r.GET("/customers?id=:id", controllers.GetSpecificCustomer)
	r.PUT("/customers/", controllers.UpdateCustomer)
	r.DELETE("/customers?id=:id", controllers.DeleteCutomer)

	//--------------------- customer ---------------------------||

	//--------------------- Invoice ----------------------------||

	r.GET("/invoices", controllers.GetAllInvoices)
	// r.GET("/invoices/customers?id=:id", controllers.GetAllInvoices)
	r.POST("/invoices", controllers.CreateInvoice)
	r.GET("/invoices?id=:id", controllers.GetSpecificInvoice)
	r.GET("/invoices/customers?id=:id", controllers.GetSpecificInvoice)
	// r.PUT("/customers?id=:id", controllers.UpdatedInvoice)
	r.DELETE("/invoices?id=:id", controllers.DeleteInvoices)

	//--------------------- Invoice ----------------------------||

	//--------------------- Payments ---------------------------||

	r.GET("/payments", controllers.GetPayments)
	r.POST("/payments", controllers.CreatePayments)
	// r.PUT("/payments?id=:id", controllers.UpdateCustomer)
	r.DELETE("/payments?id=:id", controllers.DeletePayments)

	//--------------------- Payments ---------------------------||

	//--------------------- Currency ---------------------------||

	r.GET("/currencies", controllers.GetAllCurrency)

	//--------------------- Currency ---------------------------||

	//-------------------- AirBooking --------------------------||

	r.GET("/airbooking", controllers.GetAllItems)

	//-------------------- Imputation --------------------------||

	r.GET("/imputations", controllers.GetAllInvoicePayments)
	r.POST("/imputations", controllers.CreateInvoiceImputations)

	//-------------------- Imputation --------------------------||

	//-------------------- Transaction -------------------------||

	//-------------------- Transaction -------------------------||
}
