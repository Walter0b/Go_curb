-- INSERT INTO "customer" ("customer_name",
--                         "street",
--                         "city",
--                         "state",
--                         "zip_code",
--                         "notes",
--                         "terms",
--                         "account_number",
--                         "tax_id",
--                         "balance",
--                         "is_active",
--                         "is_sub_agency",
--                         "language",
--                         "slug",
--                         "id_currency",
--                         "id_country",
--                         "irs_share_key",
--                         "currency_rate",
--                         "agency",
--                         "avoid_deletion",
--                         "is_editable",
--                         "alias",
--                         "already_used",
--                         "ab_key",
--                         "tmc_client_number",
--                         "id")
-- VALUES ('John Doe',
--         '123 Main St',
--         'New York',
--         'NY',
--         '10001',
--         'A new customer',
--         30,
--         '12345678',
--         '789012345',
--         '1000.0',
--         true,
--         false,
--         'English',
--         1,
--         1,
--         1,
--         'IRSShare',
--         1,
--         'MyAgency',
--         false,
--         true,
--         'Alias123',
--         0,
--         'AB123',
--         'TMC123',
--         52) RETURNING "id"
--         ;
INSERT INTO
        "invoice" (
                "creation_date",
                "invoice_number",
                "status",
                "due_date",
                "amount",
                "balance",
                "purchase_order",
                "customer_notes",
                "terms",
                "terms_conditions",
                "credit_apply",
                "rate",
                "net_amount",
                "tax_amount",
                "base_amount",
                "detail_taxes",
                "slug",
                "id_customer",
                "void_summary",
                "credit_used",
                "email",
                "printed_name",
                "hidden_field",
                "hidden_identifier",
                "already_used",
                "is_opening_balance",
                "tag"
        )
VALUES
        (
                '2023-11-16 12:25:37.128',
                '',
                'unpaid',
                '2023-12-31 00:00:00',
                '1',
                '',
                '',
                '',
                0,
                '',
                '',
                '',
                '',
                '',
                '',
                'map[]',
                0,
                123,
                'map[]',
                '',
                '',
                '',
                '',
                '',
                0,
                false,
                ''
        ) RETURNING "id"
SELECT
        "id",
        "creation_date",
        "invoice_number",
        "status",
        "due_date",
        "amount",
        "balance",
        "net_amount",
        "tax_amount",
        "base_amount",
        "purchase_order",
        "customer_notes",
        "terms",
        "terms_conditions",
        "credit_apply",
        "credit_used",
        "email",
        "printed_name",
        "hidden_field",
        "hidden_identifier",
        "already_used",
        "is_opening_balance",
        "tag",
        "id_customer"
FROM
        "invoice_payment_received"
WHERE
        tag = '2'
LIMIT
        7