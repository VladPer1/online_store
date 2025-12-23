CREATE DATABASE Sports_supplement_store;

-- Users table (пользователи)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    status VARCHAR(20) DEFAULT 'user'
);

-- Products table (товары)
CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    stock_qty INTEGER NOT NULL DEFAULT 0 CHECK (stock_qty >= 0),
    image_url VARCHAR(255)
);

-- Categories table (категории)
CREATE TABLE categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

-- Producers table (производители)
CREATE TABLE producers (
    producer_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    country VARCHAR(100)
);

-- Products-Categories relationships (многие-ко-многим)
CREATE TABLE products_categories (
    product_id INTEGER REFERENCES products(product_id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(category_id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

-- Products-Producers relationships (многие-ко-многим)
CREATE TABLE products_producers (
    product_id INTEGER REFERENCES products(product_id) ON DELETE CASCADE,
    producer_id INTEGER REFERENCES producers(producer_id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, producer_id)
);

-- Shopping carts (корзины)
CREATE TABLE carts (
    cart_id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE
);

-- Cart items (элементы корзины)
CREATE TABLE cart_items (
    cart_item_id SERIAL PRIMARY KEY,
    cart_id INTEGER REFERENCES carts(cart_id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(product_id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    UNIQUE (cart_id, product_id)
);

-- Orders (заказы)
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    total_amount NUMERIC(12,2) NOT NULL CHECK (total_amount >= 0)
);

-- Order items (элементы заказа)
CREATE TABLE order_items (
    order_id INTEGER REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(product_id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(10,2) NOT NULL CHECK (unit_price >= 0),
    PRIMARY KEY (order_id, product_id)
);

-- Pending registrations (ожидающие регистрации)
CREATE TABLE pending_registrations (
    placeholder_user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    verification_code INTEGER NOT NULL
);

-- Forgot password data (восстановление пароля)
CREATE TABLE forgot_password_date (
    email VARCHAR(255) UNIQUE NOT NULL,
    verification_code INTEGER NOT NULL
);

CREATE INDEX idx_producers_name ON producers(name);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_price ON products(price);
