-- CREATE DATABASE
CREATE DATABASE "Ecommerce-basic";

-- AUTH TABLE
CREATE TABLE auth (
    id SERIAL PRIMARY KEY,
    public_id VARCHAR(50) NOT NULL ,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- PRODUCTS TABLE
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    price INT NOT NULL DEFAULT 0,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- TABLE TRANSACTIONS
CREATE TABLE transactions
(
    id               SERIAL PRIMARY KEY,
    user_public_id   VARCHAR(100) NOT NULL,
    product_id       INT         NOT NULL,
    product_price    INT         NOT NULL,
    amount           INT         NOT NULL,
    sub_total        INT         NOT NULL,
    platform_fee     INT         NOT NULL   DEFAULT 0,
    grand_total      INT         NOT NULL,
    status           VARCHAR(10) NOT NULL,
    product_snapshot JSONB       ,
    created_at       TIMESTAMP   DEFAULT NOW(),
    updated_at       TIMESTAMP   DEFAULT NOW()
);

-- CHECKING TRANSACTIONS
SELECT * FROM transactions as t
JOIN auth AS a ON a.public_id = t.user_public_id