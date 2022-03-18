ALTER TABLE accounts
    ADD COLUMN available_credit_limit DECIMAL(10, 2) NOT NULL AFTER document_number;