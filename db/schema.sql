CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists users (
    id uuid primary key default uuid_generate_v4 (),
    username varchar(255) not null unique,
    email varchar(255) not null unique,
    password varchar(255) not null,
    created_at numeric not null,
);

CREATE UNIQUE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_created_at ON users (created_at);

CREATE UNIQUE INDEX idx_users_email ON users (email);

CREATE TABLE IF NOT EXISTS personal_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    created_at numeric not null,
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    stock INT,
    price INT,
    weight INT,
    base_price INT,
    sku VARCHAR(100),
    variant_count INT,
    unique_code_type VARCHAR(100),
    sold BOOLEAN DEFAULT false,
    created_at BIGINT NOT NULL,
);

CREATE TABLE IF NOT EXISTS carts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    cart_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    FOREIGN KEY (cart_id) REFERENCES carts (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE INDEX idx_carts_user_id ON carts (user_id);

CREATE INDEX idx_cart_items_cart_id ON cart_items (cart_id);

CREATE INDEX idx_cart_items_product_id ON cart_items (product_id);

CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    account_id UUID NOT NULL REFERENCES accounts (id),
    amount BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    description TEXT,
    created_at BIGINT NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT check_amount_positive CHECK (amount > 0)
);

CREATE INDEX idx_accounts_user_id ON accounts (user_id);

CREATE INDEX idx_transactions_account_id ON transactions (account_id);

CREATE INDEX idx_transactions_created_at ON transactions (created_at);

CREATE TABLE IF NOT EXISTS accounts (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID NOT NULL,
balance BIGINT NOT NULL DEFAULT 0,
created_at BIGINT NOT NULL,
updated_at BIGINT NOT NULL,
deleted_at TIMESTAMP,
CONSTRAINT check_balance_non_negative CHECK (balance >= 0)
);

CREATE TABLE IF NOT EXISTS transactions (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
account_id UUID NOT NULL,
order_id UUID,
amount BIGINT NOT NULL,
type VARCHAR(20) NOT NULL,
status VARCHAR(20) NOT NULL,
description TEXT,
created_at BIGINT NOT NULL,
deleted_at TIMESTAMP,
);

CREATE TABLE IF NOT EXISTS orders (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID NOT NULL,
cart_id UUID NOT NULL,
status VARCHAR(20) NOT NULL,
total_amount BIGINT NOT NULL,
created_at BIGINT NOT NULL,
updated_at BIGINT NOT NULL,
paid_at BIGINT,
deleted_at TIMESTAMP,
);

CREATE TABLE IF NOT EXISTS order_items (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
order_id UUID NOT NULL,
product_id UUID NOT NULL,
quantity INTEGER NOT NULL,
price BIGINT NOT NULL,
created_at BIGINT NOT NULL,
);