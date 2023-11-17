## 1. Travel Items

- GET /travel_items

```JSON
[
	{
	  "id":0, //integer
	  "ticketNumber":"", //string
	  "travelerName":"", //string
	  "itinerary":"". //string
	  "totalPrice": 0.00 //float

	}, ...
]
```

## 2. Countries

- GET /countries

```JSON
[
  {
	  "id":137,
	  "name":"MALI",
	  "code":"ML",
  }, ...
]
```

## 3. Customers

- GET /customers?fields=id,name ----> Get all customers account records with only id and customerName fields

```JSON
{
  "data": [
			  {
			      "Id": 0 // integer
				  "customerName":"",  //string
			  }, ...
		],
  "pageNumber" : 0, // default
  "pageSize":10, //default to totalNumberOfRecords
  "totalRowCount": 0, //totalNumberOfRecords
}
```

- GET /customers --> Get all customers with all fields
- GET /customers?page={pageNumber}&page-size={NumberOfRecords}. Get records with pagination

```json
{
  "data" : [
	   {
	      "id": 0 // integer
		  "customerName":"",  //string
		  "state":"",       //string ----> country Name
		  "accountNumber":"", //string
		  "idCountry":0,  //Integer
		  "alias": "",   //String
		  "tmcClientNumber":"", //String
		  "abKey":"", //string
		  "slug":23535, //integer
		  "isActive": true, //boolean
	  },...
  ],
  "pageNumber" : 0, // default
  "pageSize":10, //default
  "totalRowCount": 0, //totalNumberOfRecords
}
```

- POST /customers

```JSON
payload
{

	  "customerName":"",  //string
	  "state":"",       //string ----> country Name
	  "accountNumber":"", //string
	  "idCountry":0,  //Integer
	  "alias": "",   //String
	  "tmcClientNumber":"", //String
}
```

- PATCH /customers/id
  list of updatable fields : customerName, state, accountNumber ,idCountry, alias, tmcClientNumber \* Backend validation of updatable fields

```JSON
payload
{
      "id":0 // integer
	  "customerName":"",  //string
	  "state":"",       //string ----> country Name
	  "accountNumber":"", //string
	  "idCountry":0,  //Integer
	  "alias": "",   //String
	  "tmcClientNumber":"", //String
}
```

- DELETE /customers/id

## Invoice

- 1.1 GET /invoices ---> Get all the invoice records
- 1.2 GET /invoices?page={pageNumber}&page-size={numberOfRecords} ---> Get all the invoices with pagination.
  - Backend Validation page&page-size

```json

Response (1.1 & 1.2)

{
  "data" : [
		{
		  "id":0,
		  "invoice_number": "INV-001". //string (generated from backend)
		  "idCustomer":0, //integer
		  "creationDate":"2022-09-09", // date in this format
		  "dueDate":"2022-10-19" //date in this format,
		  "amount":10000.00 //float. ---> SUM of total price of all travel_item linked to invoice
		  "status":"", //string
		  "balance":0.00, //float
		  "credit_apply":0.00, //float
		  "travelItems":[
			  {
				  "id":0, //integer
				  "ticketNumber":"", //string
				  "travelerName":"", //string
				  "itinerary":"". //string
				  "totalPrice": 0.00 //float

			}, ...
		  ]
		},...
  ],
  "pageNumber" : 0, // default
  "pageSize":10, //default to total number of records if not set
  "totalRowCount": 0, //totalNumberOfRecords
}
```

- 1.3. GET /invoices?embed=customer ---> adds customer information to each invoice
- 1.4. GET /invoices?embed=customer&page={pageNumber}&page-size={numberOfRecords}

```json
response (1.3 & 1.4)

{
  "data" : [
		{
		  "id":0,
		  "invoice_number": "INV-001". //string (generated from backend)
		  "creationDate":"2022-09-09", // date in this format
		  "dueDate":"2022-10-19" //date in this format,
		  "amount":10000.00 //float. ---> SUM of total price of all travel_item linked to invoice
		  "status":"", //string
		  "balance":0.00, //float
		  "credit_apply":0.00, //float
		  "travelItems":[
			  {
				  "id":0, //integer
				  "ticketNumber":"", //string
				  "travelerName":"", //string
				  "itinerary":"". //string
				  "totalPrice": 0.00 //float

			}, ...
		  ],
		  "customer":{
		      "Id": 0 // integer
			  "customerName":"",  //string
			  "state":"",       //string ----> country Name
			  "accountNumber":"", //string
			  "idCountry":0,  //Integer
			  "alias": "",   //String
			  "tmcClientNumber":"", //String
			  "abKey":"", //string
			  "slug":23535, //integer
			  "isActive": true, //boolean
			},
		},...
  ],
  "pageNumber" : 0, // default
  "pageSize":10, //default to total number of records if not set
  "totalRowCount": 0, //totalNumberOfRecords
}
```

- POST /invoices

  ````JSON
  Payload

  	{
  	  "idCustomer":0, //integer
  	  "dueDate":"2022-10-19" //date in this format,
  	  "amount":10000.00 //float. ---> SUM of total price of all travel_item linked to invoice
  	  "travelItems":[
  		  {
  			  "id":0, //integer
  			  "totalPrice": 0.00 //float

  		}, ...
  	  ]
  	}

  	```

  ````

- PATCH /invoices/id

  - updatable fields : travelItems, dueDate, idCustomer
    - Backend Validation on updatable fields

  ````JSON
  Payload

  	{
  	  "id":1
  	  "idCustomer":0, //integer
  	  "dueDate":"2022-10-19" //date in this format,
  	  "travelItems":[
  		  {
  			  "id":0, //integer
  			  "totalPrice": 0.00 //float

  		}, ...
  	  ]
  	}

  	```

  ````

- DELETE /invoices/id

## Payments

- 1.1. GET /payments ---> get all the payment records (no pagination)
- 1.2. GET /payments?page={pageNumber}&page-size={numberOfRecords} ---> gets paginated payment records

```json

response (1.1 & 1.2):

{
  "data" : [
	    {
		  "id":1
		  "paymentNumber":"PER-001", //string (generated from backend)
		  "paymentDate":"2022-10-19", //date in this format,
		  "paymentMode": "cash", //string, --> enum values(cash,check,bank_tranfer)
		  "idCustomer":0,
		  "amount":0.00, //float
		  "balance":0.00, //float
		  "usedAmount": 0.00, //float
		  "status": "open", //string ---> enum values(open,used)
		}, ...

  ],
  "pageNumber" : 0, // default
  "pageSize":10, //default
  "totalRowCount": 0, //totalNumberOfRecords
}
```

- 1.3. GET /payments?embed=customer ---> adds customer info to each payment
- 1.4. GET /payments?embed=customer&page={pageNumber}&page-size={numberOfRecords}

```json

response (1.3 & 1.4):

{
  "data" : [
	    {
		  "id":1
		  "paymentNumber":"PER-001", //string (generated from backend)
		  "paymentDate":"2022-10-19", //date in this format,
		  "paymentMode": "cash", //string, --> enum values(cash,check,bank_tranfer)
		  "amount":0.00, //float
		  "balance":0.00, //float
		  "usedAmount": 0.00, //float
		  "status": "open", //string ---> enum values(open,used)
		  "customer":{
		      "Id": 0 // integer
			  "customerName":"",  //string
			  "state":"",       //string ----> country Name
			  "accountNumber":"", //string
			  "idCountry":0,  //Integer
			  "alias": "",   //String
			  "tmcClientNumber":"", //String
			  "abKey":"", //string
			  "slug":23535, //integer
			  "isActive": true, //boolean
			},
		}, ...

  ],
  "pageNumber" : 0, // default
  "pageSize":10, //default
  "totalRowCount": 0, //totalNumberOfRecords
}
```

- 1.5. GET /customers/id/payments ---> get all payments linked to a particular customer.
- 1.6. GET /customers/id/payments?page={pageNumber}&page-size={numberOfRecords} ----> get a set of customer payments
- 1.7. GET /customers/id/payments?used=false ---> get list of all customer payments whose balance != 0

```json

response (1.5 & 1.6 & 1.7):

{
  "data" : {
         "idCustomer": 0,
         "payments": [
	         {
			  "id":1
			  "paymentNumber":"PER-001", //string (generated from backend)
			  "paymentDate":"2022-10-19", //date in this format,
			  "paymentMode": "cash", //string, --> enum values(cash,check,bank_tranfer)
			  "amount":0.00, //float
			  "balance":0.00, //float
			  "usedAmount": 0.00, //float
			  "status": "open", //string ---> enum values(open,used)
			}, ...
		]
  },
  "pageNumber" : 0, // default
  "pageSize":10, //default
  "totalRowCount": 0, //totalNumberOfRecords
}
```

- POST /payments

  ```JSON
  payload

  	{
  		"IdCustomer":"", //string
  		"amount": 0.00, // float
  		"paymentMode" : "cash"  // string --> enum values(cash,check,bank_tranfer)
  	}
  ```

- PATCH /payments/id
  - Updatable fields : idCustomer, amount, paymentMode. Backend Validation of these fields.

```JSON
	payload

		{
			"IdCustomer":"", //string
			"amount": 0.00, // float
			"paymentMode" : "cash"  // string --> enum values(cash,check,bank_tranfer)
		}

```

- DELETE /payments/id

## Invoice imputations

- GET /invoices/id/imputations ---> get all imputations linked to the invoice
- GET /invoices/id/imputations?page={pageNumber}&page-size={totalNumberOfRecords}

```json

response
{
 "data" : [
			  {
				  "idInvoice": 0,
				  "imputations":[
					{
						"payment":{
						  "id":0,
						  "paymentNumber":"PER-001",
						  "amount":100.00, // float
						  "balance": 100.00 // float
						}
						"amountApplied":0.00 // float

					}, ...
			     ]
			 }
	]


 },
 "pageNumber" : 0, // default
 "pageSize":10, //default
 "totalRowCount": 0, //totalNumberOfRecords
}
```

- POST /invoices/id/imputations

```JSON
 payload
		{

			"idInvoice": 0, //integer
			"imputations":[
				{
					"idPayment":0 //integer,
					"amountApplied":0.00 // float

				}, ...
			]
		}
```

- PATCH /invoice/id/imputations
  - Updatable fields : imputations.amountApplied

````JSON
 payload
		{

			"idInvoice": 0, //integer
			"imputations":[
				{
					"idPayment":0 //integer, //backend validation here, check whether link exist
					"amountApplied":0.00 // float

				}, ...
			]
		}
	```

````
