# Online Store


## ✨ Features

### 🛒 Full E-commerce System
- Product catalog with filtering and sorting
- Shopping cart functionality
- Order processing with payment simulation
- Stock management and validation
- Email order confirmations

### 🔐 Advanced Authentication
- User registration with email verification
- JWT token-based authentication
- Password reset with email confirmation
- Role-based access (User/Admin)
- Secure password hashing

### 👤 User Management
- Personal profile with order history
- Profile editing capabilities
- Password change functionality
- Account deletion with confirmation

### 🛠️ Admin Panel
- Complete product CRUD operations
- User management (ban/unban)
- CSV export/import of products
- Excel report generation
- Real-time stock monitoring

## 🚀 Quick Start / Installation

### Prerequisites
- Go 1.20+
- PostgreSQL 14+
- Gmail Account

### 1. Clone Repository
```bash
git clone https://github.com/VladPer1/online-store.git
cd online-store/backend
go mod download

```

### 2. Configure Environment

```bash
cp ../.env.example ../.env
# Edit .env file with your settings
```

### 3. Setup Database

```sql
CREATE DATABASE Sports_supplement_store;
```

### 4. Run the application:
```
go run main.go
```
Server starts at: http://localhost:8080

## Environment Variables

Create `.env` file with:
```env
# Database

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=Sports_supplement_store

# JWT

JWT_SECRET=your-secure-jwt-key

# Email

SMTP_PASSWORD=your-gmail-app-password


# File Paths (Windows example)

TEMPLATE_PATH=D:\online_store\frontend\templates
STATIC_PATH=D:\online_store\frontend\static
```

## Project Structure

```

online-store/
├── backend/                    # Go backend
│   ├── handlers/              # HTTP handlers
│   │   ├── auth_handler.go    # Authentication
│   │   ├── cart_handler.go    # Shopping cart
│   │   ├── catalog_handler.go # Product catalog
│   │   ├── payment_handler.go # Payments
│   │   ├── profile_handler.go # User profiles
│   │   ├── admin_handler.go   # Admin panel
│   │   └── routes.go          # URL routing
│   ├── database/              # Database layer
│   │   ├── connection.go      # DB connection
│   │   ├── user_repository.go # User operations
│   │   ├── cart_repository.go # Cart operations
│   │   ├── orders_repository.go # Orders
│   │   ├── filters_repository.go # Filtering
│   │   └── admin_repository.go # Admin functions
│   ├── models/                # Data models
│   ├── utils/                 # Utilities
│   ├── server/                # Server setup
│   └── main.go                # Entry point
├── frontend/                  # Frontend files
│   ├── templates/             # HTML templates
│   └── static/                # CSS/JS/images
│
├── .env.example               # Config template
├── go.mod                     # Go module
└── README.md                  # Documentation

```

## Database Setup

1. Create PostgreSQL database:
```sql

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

```


## Usage

The application provides the following functionality:

- User authentication (login, registration, password reset)
- Product catalog with filtering and sorting
- Shopping cart management
- Order processing
- User profile management
- Admin panel for managing products, users, and generating reports

## API

The application exposes the following API endpoints:

### Authentication

- `GET / `: Login page
- `POST /login`: User login
- `POST /register`: User registration
- `GET /verify-code-page`: Email verification
- `POST /verify-code`: Verify email code
- `GET /forgot-password`: Password reset
- `POST /forgot-password-send`: Send reset code
- `GET /forgot-password-verify-page`: Verify reset code page
- `POST /forgot-password-verify`: Verify reset code
- `GET /forgot-password-update-password-page`: New password form
- `POST /forgot-password-update-password`: Update password
- `POST /logout`: User logout

### Store

- `GET /catalog`: Product catalog with filters
- `POST /add-to-cart`: Add to shopping cart
- `GET /cart`: View cart
- `POST /delete-from-cart`: Remove from cart
- `POST /process-payment`: Create order
- `GET /success-payment`: Payment success
- `GET /error-payment`: Payment error

### User Profile

- `GET /profile`: View profile
- `GET /update-profile-page`: Edit profile form
- `POST /update-profile`: Update profile
- `GET /change-password-page`: Change password form
- `POST /change-password`: Change password
- `GET /delete-account-page`: Delete account page
- `POST /delete-account`: Delete account

### Admin Panel
- `GET /admin`: Admin dashboard
- `POST /admin/products/save`: Save/update product
- `POST /admin/products/delete`: Delete product
- `POST /admin/users/ban`: Ban user
- `POST /admin/users/unban`: Unban user
- `GET /admin/report`: Generate Excel report
- `GET /admin/export-csv`: Export to CSV
- `POST /admin/import-csv`: Import from CSV


## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push your changes to your forked repository.
5. Submit a pull request.
