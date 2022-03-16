CREATE TABLE accounts
(
    id              INT         NOT NULL AUTO_INCREMENT PRIMARY KEY,
    document_number VARCHAR(14) NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE operation_types
(
    id          INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE transactions
(
    id                INT            NOT NULL AUTO_INCREMENT PRIMARY KEY,
    account_id        INT,
    operation_type_id INT,
    amount            DECIMAL(10, 2) NOT NULL,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id)
        REFERENCES accounts (id)
        ON DELETE CASCADE,
    FOREIGN KEY (operation_type_id)
        REFERENCES operation_types (id)
        ON DELETE CASCADE
);