# E-commerce Platform API

A multi-vendor e-commerce platform built with Go, following Clean Architecture and Domain-Driven Design (DDD) principles.

## 📋 Table of Contents

- [Features](#features)
- [System Requirements](#system-requirements)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Installation](#installation)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Development](#development)
- [Testing](#testing)

## ✨ Features

### For Customers

- Browse products from all shops
- Search and filter products
- Create orders and purchase products
- Track order status and shipping information
- View order history

### For Shop Owners

- Create and manage shop
- Add, edit, and delete products
- Manage inventory
- Process orders (change status, cancel)
- Add shipping/courier information
- View shop-specific orders

## 🔧 System Requirements

- Go 1.21 or higher
- PostgreSQL 14+
- Docker & Docker Compose (optional)

## 🛠 Tech Stack

- **Framework:** Echo

- [Echo Web Framework](https://github.com/labstack/echo)
- [Echo Web Framework](https://github.com/labstack/echo)
- [Echo Web Framework](https://github.com/labstack/echo)
  ecommerce/
  ├── domain/ # Domain entities and business logic
  │ ├── user.go
  │ ├── shop.go
  │ ├── product.go
  │ ├── order.go
  │ ├── order_item.go
  │ ├── shipping.go
  │ └── errors.go
  │
  ├── feature/ # Features organized by domain
  │ ├── auth/
  │ │ ├── delivery/ # HTTP handlers
  │ │ ├── usecase/ # Business logic
  │ │ └── repository/ # Data access
  │ ├── shop/
  │ │ ├── delivery/ # HTTP handlers
  │ │ ├── usecase/ # Business logic
  │ │ └── repository/ # Data access
  │ ├── product/
  │ │ ├── delivery/ # HTTP handlers
  │ │ ├── usecase/ # Business logic
  │ │ └── repository/ # Data access
  │ ├── order/
  │ │ ├── delivery/ # HTTP handlers
  │ │ ├── usecase/ # Business logic
  │ │ └── repository/ # Data access
  │ ├── shop/
  │ ├── product/
  │ ├── order/
  │ └── shipping/
  ├── middleware/ # HTTP middlewares
  │ ├── auth.go
  │ ├── role.go
  │ └── shop_owner.go
  │
  ├── internal/ # Shared packages
  │ ├── response/
  │ ├── validation/
  │ └── pagination/
  │
  ├── config/ # Configuration
  │ └── config.go
  │
  ├── migrations/ # Database migrations
  ├── docs/ # API documentation
  └── main.go

````

## 🚀 Installation

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/ecommerce-platform.git
cd ecommerce-platform

# Start services
docker-compose up -d

# Run migrations
docker-compose exec api go run migrations/migrate.go up
````

### Manual Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/ecommerce-platform.git
cd ecommerce-platform

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit .env with your configuration
nano .env

# Run migrations
go run migrations/migrate.go up

# Start the server
go run main.go
```

## ⚙️ Configuration

Create a `.env` file in the root directory:

```env
# Server
PORT=8080
ECHO_MODE=release

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=ecommerce
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24

# CORS
ALLOWED_ORIGINS=http://localhost:3000

# File Upload
MAX_UPLOAD_SIZE=10485760  # 10MB
UPLOAD_PATH=./uploads
```

## 📚 API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

All authenticated endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### API Endpoints

#### Authentication

```
POST   /auth/register          Register new user (shop/customer)
POST   /auth/login             Login
POST   /auth/logout            Logout
GET    /auth/me                Get current user info
```

#### Shops

```
GET    /shops                  Get all shops
GET    /shops/:shopId          Get shop details
GET    /shops/me               Get my shop (shop owner)
POST   /shops/me               Create shop (shop owner)
PUT    /shops/me               Update my shop (shop owner)
DELETE /shops/me               Delete my shop (shop owner)
```

#### Products

```
GET    /products               Get all products
GET    /products/search        Search products
GET    /products/:productId    Get product details
GET    /shops/:shopId/products Get products by shop

GET    /products/me            Get my products (shop owner)
POST   /products/me            Create product (shop owner)
PUT    /products/me/:productId Update product (shop owner)
DELETE /products/me/:productId Delete product (shop owner)
POST   /products/me/:productId/images  Upload images (shop owner)
```

#### Orders

```
POST   /orders                 Create order (customer)
GET    /orders/me              Get my orders (customer)
GET    /orders/:orderId        Get order details
PUT    /orders/:orderId/cancel Cancel order (customer)

GET    /orders/shop            Get shop orders (shop owner)
PUT    /orders/:orderId/status Update order status (shop owner)
PUT    /orders/:orderId/cancel Cancel order (shop owner)
POST   /orders/:orderId/shipping Add shipping info (shop owner)
```

#### Shipping

```
GET    /couriers               Get all couriers
GET    /orders/:orderId/tracking Track order
```

### Request Examples

#### Register as Customer

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123",
    "name": "John Doe",
    "role": "customer"
  }'
```

#### Register as Shop

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "shop@example.com",
    "password": "password123",
    "name": "Shop Owner",
    "role": "shop"
  }'
```

#### Create Product

```bash
curl -X POST http://localhost:8080/api/v1/products/me \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Product Name",
    "description": "Product Description",
    "price": 999.99,
    "stock": 100,
    "category": "Electronics"
  }'
```

#### Create Order

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": "uuid",
        "quantity": 2
      }
    ],
    "shipping_address": {
      "address": "123 Main St",
      "city": "Bangkok",
      "postal_code": "10110"
    }
  }'
```

## 🗄 Database Schema

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('customer', 'shop')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Shops Table

```sql
CREATE TABLE shops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    logo VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(owner_id)
);
```

### Products Table

```sql
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Orders Table

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES users(id),
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Order Items Table

```sql
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    shop_id UUID NOT NULL REFERENCES shops(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Shipping Table

```sql
CREATE TABLE shipping (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    courier_name VARCHAR(255) NOT NULL,
    tracking_number VARCHAR(255),
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 👨‍💻 Development

### Running in Development Mode

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Code Style

This project follows Go standard coding conventions:

- Run `gofmt` before committing
- Follow Go Code Review Comments
- Use meaningful variable names

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run
```

## 🧪 Testing

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Run Specific Feature Tests

```bash
go test ./feature/product/...
```

### Generate Coverage Report

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔐 Security Considerations

- Passwords are hashed using bcrypt
- JWT tokens for authentication
- Input validation on all endpoints
- SQL injection prevention using parameterized queries
- CORS configuration
- Rate limiting (recommended for production)

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📧 Contact

Your Name - your.email@example.com

Project Link: [https://github.com/yourusername/ecommerce-platform](https://github.com/yourusername/ecommerce-platform)

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
