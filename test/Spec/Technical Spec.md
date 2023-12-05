## Tables

- payment_received
- customer
- invoice
- invoice_payment_received

## Common

Add paginations to invoice, customers and payments.

Paginated Routes structure : GET /{resource}?page={pageNumber}&pageSize={numberOfRows}

NB : {resource} : invoices | customers | payments

Response :

```json
{
  "data": [],
  "pageNumber": 0, // default
  "pageSize": 10, //default
  "totalRowCount": 0 //totalNumberOfRecords
}
```

## Travel Items

GET /travel_items
Query for travel items:

```SQL
 SELECT total_price,itinerary,traveler_name,ticket_number,conjunction_number,status
	 FROM air_booking
	 WHERE id_invoice IS NULL
		 AND product_type='flight'
		 AND transaction_type='sales'
		 AND status = 'pending'

```

## Countries List

GET /countries

- Query :

```sql
SELECT * FROM country WHERE currency_code='XOF'
```

## Customers

- DB table name : 'customer'
- Constants (Backend) :

  - currency Id (XOF) --> 550;

- Create ---> _**/POST customers**_

  - payload (Frontend)

    ```JSON
    {
    	"customerName":"",  //string
    	"state":"",       //string ----> country Name
    	"accountNumber":"", //string
    	"idCountry":0,  //Integer
    	"alias": "",   //String
    	"tmcClientNumber":"" //String

    }
    ```

- Required fields (Backend) :

  - Id (integer)
  - customer_name(string)
  - state(String) ---> Optional
  - account_number (String) - has unique constraint
  - slug (bigInt) - has unique constraint
  - Id_currency(integer)
  - id_country(integer)
  - is_active(boolean) ---> Optional , default to true
  - ab_key (string)
  - tmc_client_number(string) - has unique constraint

- Update ----> _**/PUT customers/id**_

  - Constraints : none;

- Delete -----> _**/DELETE customers/Id**_

  - Constraints : Can't delete a customer that has invoices.

## Invoice

- DB table name : 'invoice'
- **Creation Process**

  - Choose a customer account from a list of customer account
    Task :Send a list of customers with their Id and names --->(route): GET /customers?fields=customer.id,customer.name
  - Enter due date
  - Select billable Travel items to add to invoice
  - Save ---->(route) POST /invoices
  - Payload (frontend) for post creation

    ```JSON
     new row for relation "invoice" violates check constraint "balance_ck"{
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

  - when a new invoice is created if it has travel items, the id of the invoice is added to the corresponding air_booking records.

- Required fields for Invoice Creation

  - creation_date (timestamp without timezone)
  - invoice_number (string) , must be unique (generate)
  - status (ab_invoice_status)
  - due_date (date)
  - amount (money) ---> net_amount + tax_amount. Tax_amount default to Zero (just save invoice amount directly in this field).
  - balance (money) ---> (Not Required)
  - net_amount (money) - (Not Used) ---> Total net Amount of travel Items within the invoice >> (might be used in case of constraints on amount)
  - Credit apply (money) - (Not required) --->
  - base_amount (money) ---> Amount converted in the base currency, if invoice is in base currency then same value from amount field >>
  - slug (bigInt) , must be unique ---> should be a random value
  - already_used(integer).
  - tags (group data tags)
  - id_customer (integer) --> ForeignKey
  - void_summary (json) ---> used to keep the list of travel_item linked to invoice

- DB constraints:
  - "status_ck" CHECK (status <> 'invoiced'::ab_booking_status AND id_invoice IS NULL AND id_credit_note IS NULL OR status = 'invoiced'::ab_booking_status AND (id_invoice IS NOT NULL OR id_credit_note IS NOT NULL))
  - "balance_ck" CHECK (balance = (amount - credit_apply - credit_used) AND balance >= 0::money)
  - "credit_apply_ck" CHECK (credit_apply >= 0::money)
  - "due_date_later" CHECK (due_date >= creation_date::date)
  - "id_check" CHECK (id > 0)
  - "invoice_status_ck" CHECK (status = 'paid'::ab_invoice_status AND balance = 0::money OR status = 'unpaid'::ab_invoice_status AND balance > 0::money OR status = 'void'::ab_invoice_status AND balance = 0::money)
- Constraints on Update & Delete Operations

  - An invoice can only be Updated / deleted if the invoice amount hasn't been charged i.e there's is no record in "invoice_payment_received" and " invoice.credit_apply == 0 " && "invoice.balance == invoice.amount"
  - Update Operations :
    - Change due Date
    - Add or remove travel_items. ----> affects the amount of the invoice
    - Credit_apply (in case the amount charged is refunded or changed). If changed Balance is recalculated accordingly. (formula ---> balance = amount -credit_apply)
  - Routes :
    - PUT invoices/id (full update), payload similar to that of POST
    - DELETE invoices/id

- Validations & Constraints (Application Level) :
  - Balance should be positive
  - Credit apply should be positive
  - Amount == Balance + credit apply

## Payments

- DB Table name : payment_received
- Creation Process

  - Enter customer account name ----> GET /customers?fields=id,name

    - When a customer is selected on the UI, a the list of customer unpaid invoices is shown in a table ----> (route):GET /customers/id/invoices?page=0&pageSize=10&paid=false

  - Choose deposit Account ---> Deposit Account Id : 39, Account Name = PettyCash
  - Chose payment Method : (Bank Transfer, Cash, Cheque, Mobile, Pos)
  - Enter Amount
  - Save ---->(route): POST /payments

    payload (frontend)

    ```JSON
    {
    	"IdCustomer":"", //string
    	"amount": 0.00, // float
    	"paymentMode" : "". // string
    }
    ```

- Task : When we chose a customer, send the list of invoice not paid related to the customer
- Required fields for Payment Creation

  - number(string) ---> Payment Number, unique (generate)
  - date (date) ---> Payment date , # on save default to creation date.
  - balance (money)
  - amount (money)
  - fop (enum --> ab_payment_fop) (string)
  - used_amount (money) ---> default to 0
  - status (enum --> ab_payment_status) --> Payment status
  - base_amount (money) ----> Amount converted in base currency. same as amount
  - slug (bigInt) ---> random value generated before save
  - id_customer (integer) --> FK
  - id_charts_of_accounts (integer) ---> FK, == deposit account Id.
  - id_currency(integer) ---> Fk
  - tags (group data tags)

- DB Check constraints:

  - "id_check" CHECK (id > 0)
  - "payment_received_balance_ck" CHECK (balance = (amount - used_amount))
  - "payment_received_status_ck" CHECK (status = 'used'::ab_payment_status AND balance = 0::money OR status = 'open'::ab_payment_status AND balance <> 0::money OR status = 'void'::ab_payment_status AND balance = 0::money)

- Constraints on Update & Delete Operations

  - Only payment records that haven't been used can be updated on deleted manually by the user.
  - When an invoice payment is created, payment.used_amount is updated. When updated the payment balance is recalculated
  - When an invoice payment is updated,
    - if amount charged was increased, payment.used_amount is increased accordingly and the balance is recalculated
    - if amount charged was decreased, payment.used_amount is decreased accordingly and the balance is recalculated
  - Routes :
    - PUT /payments/id
    - DELETE /payments/id

- GET /customers/id/payments ---> get all payments linked to a particular customer
- GET /customers/id/payments?used=false ---> get list of all customer payments whose balance != 0

## Invoice Payment Received

- DB table name : invoice_payment_received
- Creation Process
  Payload (frontend) Creation process ---> POST /invoice/id/imputations

  ```JSON
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

  - When records are created :
    - invoice.credit_apply (invoice table) is updated
      - all the amount charged (in invoice_payment_record) are summed up and the result is placed in credit_apply for the given invoice. (see update on invoice)

- Required fields :

  - id_invoice (integer) FK
  - id_payment_received(integer)
  - amount_apply(money) << Amount applied >>
  - payment_amount(money) << Amount applied from payment in base currency >> ?
  - invoice_amount(money) << Amount applied from the invoice in base currency >> ?
  - id(integer)
  - tag(tags)

- Update and Delete Operations

  - Routes :
    - PATCH /invoice/id/imputations
  - Rules :

    - If the new value for amount charged is zero (and the previous value is > 0), then :

      - 1. amount applied(previous value) is deducted from invoice.credit_apply (in invoice table), and balance recalculated
      - 2. amount applied(previous value) is deducted from payment.used_amount(in payment_received table), and balance recalculated
      - 3. The record is deleted

    - if the new value for amount charged is different from zero, and different from saved value, invoice.credit_apply and payment.used_amount are updated accordingly.
