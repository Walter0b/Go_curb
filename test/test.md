Customers

1. Get All Customers
   Route: GET /customers
   Description: Retrieves all customers.
   JSON Response Example:

```JSON

{
  "data": [
        {
            "ID": 87,
            "Customer_name": "SAFCOL",
            "Street": "Michael Brink Place",
            "City": "Capetown",
            "State": "Western Cape",
            "Zip_code": "ZA565420",
            "Notes": "",
            "Terms": 30,
            "Account_number": "",
            "Tax_id": "ZA45419-PTY",
            "Balance": "$3,958,147.86",
            "Is_active": true,
            "Is_sub_agency": false,
            "Language": "",
            "Slug": 110000000091338,
            "Id_currency": 455,
            "Id_country": 205,
            "Irs_share_key": "",
            "Currency_rate": 2.90241,
            "Agency": "AB",
            "Avoid_deletion": false,
            "Is_editable": true,
            "Alias": "",
            "Already_used": 12,
            "Ab_key": "",
            "Tmc_client_number": ""
        }
        ],
    "totalRowCount": 81
}

```

2. Create Customer
   Route: POST /customers
   Description: Creates a new customer.
   JSON Request Example:

```JSON

  {
  "Ab_key": "AB456",
  "Account_number": "987654321",
  "Agency": "NewAgency",
  "Alias": "NewAlias123",
  "Already_used": 1,
  "Avoid_deletion": true,
  "Balance": "$2,000.00",
  "City": "Los Angeles",
  "Currency_rate": 1.2,
  "Customer_name": "Alice",
  "ID": 22,
  "Id_country": 205,
  "Id_currency": 550,
  "Irs_share_key": "NewIRSShare",
  "Is_active": false,
  "Is_editable": false,
  "Is_sub_agency": true,
  "Language": "fr",
  "Notes": "Another new customer",
  "Slug": 4,
  "State": "CA",
  "Street": "456 Elm St",
  "Tax_id": "987654321",
  "Terms": 45,
  "Tmc_client_number": "TMC5678",
  "Zip_code": "90001",
  "isEditing": true
}

```

JSON Response Example:

```JSON

  {
  {
  "Ab_key": "AB456",
  "Account_number": "987654321",
  "Agency": "NewAgency",
  "Alias": "NewAlias123",
  // ... other customer fields
  }
  }
```

3. Get Specific Customer
   Route: GET /customers/:id
   Description: Retrieves details of a specific customer.
   JSON Response Example:

```JSON

  {
  "id": 2,
  "name": "Jane Doe",
  "email": "jane@example.com",
  "phone": "987-654-3210"
  // ... other customer fields
  }
```

4. Update Customer
   Route: PUT /customers/:id
   Description: Updates details of a specific customer.
   JSON Request Example:

```JSON

 {
    "ID": 22,
    "customer_name": "Alice Smith",
    "street": "456 Elm St",
    "city": "Los Angeless",
    "state": "CA",
    "zip_code": "90001",
    "notes": "Another new customer",
    "terms": 45,
    "account_number": "987654321",
    "tax_id": "987654321",
    "balance": "2000.0",
    "credit_limit": 8000.0,
    "is_active": false,
    "is_sub_agency": true,
    "opening_balance": 800.0,
    "language": "fr",
    "slug": 4,
    "id_currency": 550,
    "id_country": 205,
    "irs_share_key": "NewIRSShare",
    "currency_rate": 1.2,
    "agency": "NewAgency",
    "opening_balance_date": "2023-11-01",
    "avoid_deletion": true,
    "is_editable": false,
    "alias": "NewAlias123",
    "already_used": 1,
    "ab_key": "AB456",
    "tmc_client_number": "TMC5678"
}

```

JSON Response Example:

```JSON

  {
  "message": "Customer updated successfully",
  "data": {
  "id": 22,
    "customer_name": "Alice Smith",
    "street": "456 Elm St",
    "city": "Los Angeless",
    "state": "CA",
    "zip_code": "90001",
  // ... other updated fields
  }
  }
```

5. Delete Customer
   Route: DELETE /customers/:id
   Description: Deletes a specific customer.

JSON Response Example:

```JSON

  {
  "message": "Customer deleted successfully"
  }

```

Invoices 6. Get All Invoices
Route: GET /invoices
Description: Retrieves all invoices.
JSON Response Example:

```JSON

  {
  "message": "List of invoices",
  "data": [
  {
  "id": 1,
  "customer_id": 1,
  "amount": 500.00,
  "status": "unpaid"
  // ... other invoice fields
  },
  // ... other invoices
  ]
  }
```

7. Create Invoice
   Route: POST /invoices
   Description: Creates a new invoice.
   JSON Request Example:

```JSON

  {
  "customer_id": 2,
  "amount": 750.00
  // ... other invoice fields
  }
```

JSON Response Example:

```JSON

  {
  "message": "Invoice created successfully",
  "data": {
  "id": 2,
  "customer_id": 2,
  "amount": 750.00,
  "status": "unpaid"
  // ... other invoice fields
  }
  }
```

8. Get Specific Invoice
   Route: GET /invoices/:id
   Description: Retrieves details of a specific invoice.
   JSON Response Example:

```JSON

  {
  "id": 2,
  "customer_id": 2,
  "amount": 750.00,
  "status": "unpaid"
  // ... other invoice fields
  }
```

9. Delete Invoices
   Route: DELETE /invoices/:id
   Description: Deletes a specific invoice.
   JSON Response Example:

```JSON

  {
  "message": "Invoice deleted successfully"
  }
```

Payments 10. Get All Payments
Route: GET /payments
Description: Retrieves all payments.
JSON Response Example:

```JSON

  {
  "message": "List of payments",
  "data": [
  {
  "id": 1,
  "invoice_id": 1,
  "amount": 500.00,
  "status": "processed"
  // ... other payment fields
  },
  // ... other payments
  ]
  }
```

11. Create Payments
    Route: POST /payments
    Description: Creates a new payment.
    JSON Request Example:

```JSON

  {
  "invoice_id": 2,
  "amount": 750.00
  // ... other payment fields
  }
```

JSON Response Example:

```JSON

  {
  "message": "Payment created successfully",
  "data": {
  "id": 2,
  "invoice_id": 2,
  "amount": 750.00,
  "status": "processed"
  // ... other payment fields
  }
  }
```

12. Get Specific Payments
    Route: GET /payments/:id
    Description: Retrieves details of a specific payment.
    JSON Response Example:

```JSON

  {
  "id": 2,
  "invoice_id": 2,
  "amount": 750.00,
  "status": "processed"
  // ... other payment fields
  }
```

13. Delete Payments
    Route: DELETE /payments/:id
    Description: Deletes a specific payment.
    JSON Response Example:

```JSON

  {
  "message": "Payment deleted successfully"
  }
```

Currency 14. Get All Currency
Route: GET /currencies
Description: Retrieves all currencies.
JSON Response Example:

```JSON

  {
  "message": "List of currencies",
  "data": [
  {
  "id": 1,
  "code": "USD",
  "name": "US Dollar"
  // ... other currency fields
  },
  // ... other currencies
  ]
  }
```

AirBooking 15. Get All Items
Route: GET /airbooking
Description: Retrieves all items related to AirBooking.
JSON Response Example:

```JSON

  {
  "message": "List of air booking items",
  "data": [
  {
  "id": 1,
  "name": "Item 1",
  // ... other item fields
  },
  // ... other items
  ]
  }
```

Imputations 16. Get All Invoice Payments
Route: GET /imputations
Description: Retrieves all invoice payments (imputations).
JSON Response Example:

```JSON

  {
  "message": "List of invoice payments",
  "data": [
  {
  "id": 1,
  "invoice_id": 1,
  "payment_id": 1,
  "amount_apply": 200.00
  // ... other imputation fields
  },
  // ... other invoice payments
  ]
  }
```

17. Create Invoice Imputations
    Route: POST /imputations
    Description: Creates a new invoice imputation.
    JSON Request Example:

```JSON

  [
  {
  "invoice_id": 2,
  "payment_id": 2,
  "amount_apply": 300.00
  // ... other imputation fields
  },
  // ... other imputations
  ]
```

JSON Response Example:

```JSON

  {
  "message": "Invoice imputations created successfully",
  "data": [
  {
  "id": 2,
  "invoice_id": 2,
  "payment_id": 2,
  "amount_apply": 300.00
  // ... other imputation fields
  },
  // ... other created imputations
  ]
  }
```
