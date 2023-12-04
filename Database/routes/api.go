package routes

import (
	"Go_curb/Database/routes/controllers"

	"github.com/gin-gonic/gin"
)

func Api(r *gin.Engine) {
	//--------------------- customer  ---------------------------||

	r.POST("/customers", controllers.CreateCustomer)
	r.GET("/customers", controllers.GetAllCustomer)
	r.PUT("/customers/", controllers.UpdateCustomer)
	r.DELETE("/customers", controllers.DeleteCutomer)

	//--------------------- customer ---------------------------||

	//--------------------- Invoice ----------------------------||

	r.GET("/invoices", controllers.GetAllInvoices)
	r.POST("/invoices", controllers.CreateInvoice)
	// r.PUT("/customers", controllers.UpdatedInvoice)
	r.DELETE("/invoices", controllers.DeleteInvoices)

	//--------------------- Invoice ----------------------------||

	//--------------------- Payments ---------------------------||

	r.GET("/payments", controllers.GetPayments)
	r.POST("/payments", controllers.CreatePayments)
	// r.PUT("/payments", controllers.UpdateCustomer)
	r.DELETE("/payments", controllers.DeletePayments)

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
