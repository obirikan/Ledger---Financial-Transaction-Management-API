# Ledger - Financial Transaction Management API

A robust, secure RESTful API built with Go for managing financial transactions, user authentication, and account management. This system uses Gin framework for routing, PostgreSQL for data persistence, and JWT for authentication.

## Table of Contents
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Project Structure](#project-structure)
- [Security](#security)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

##  Features
- **User Management**
  - Secure registration and authentication
  - JWT-based session management
  - Role-based access control
- **Account Operations**
  - Create and manage multiple accounts
  - Real-time balance tracking
  - Account statement generation
- **Transaction Management**
  - Secure fund transfers between accounts
  - Transaction history
  - Atomic transactions with rollback support
- **Security**
  - Password hashing using bcrypt
  - JWT token-based authentication
  - Rate limiting and request validation
  - SQL injection prevention

##  Technology Stack
- **Go** (1.20+)
- **Gin Framework** - HTTP web framework
- **PostgreSQL** - Primary database
- **JWT-Go** - Authentication
- **Migrate** - Database migrations
- **Godotenv** - Environment configuration

## Getting Started

### Prerequisites
```bash
# Install Go (1.20 or later)
brew install go

# Install PostgreSQL
brew install postgresql

# Install migrate tool
brew install golang-migrate
```

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/ledger.git
cd ledger
```

2. **Set up environment variables**
```bash
cp .env.example .env
```

Update `.env` with your configuration:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=ledger_db
JWT_SECRET=your_secret_key
```

3. **Initialize the database**
```bash
# Start PostgreSQL service
brew services start postgresql

# Create database
createdb ledger_db

# Run migrations
migrate -path ./migrations -database "postgresql://your_user:your_password@localhost:5432/ledger_db?sslmode=disable" up
```

4. **Build and run**
```bash
go mod tidy
go run cmd/server/main.go
```

## API Documentation

### Authentication Endpoints

#### Register New User
```http
POST /api/v1/register
Content-Type: application/json

{
    "username": "user@example.com",
    "password": "securePassword123",
    "full_name": "John Doe"
}
```

#### Login
```http
POST /api/v1/login
Content-Type: application/json

{
    "username": "user@example.com",
    "password": "securePassword123"
}
```

### Account Operations

#### Create Account
```http
POST /api/v1/accounts
Authorization: Bearer <your_token>
Content-Type: application/json

{
    "initial_balance": "1000.00",
    "account_type": "savings"
}
```

#### Get Account Details
```http
GET /api/v1/accounts/{account_id}
Authorization: Bearer <your_token>
```

### Transaction Operations

#### Create Transfer
```http
POST /api/v1/transfers
Authorization: Bearer <your_token>
Content-Type: application/json

{
    "from_account_id": "123",
    "to_account_id": "456",
    "amount": "500.00",
    "description": "Rent payment"
}
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Accounts Table
```sql
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    account_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_account_id INTEGER REFERENCES accounts(id),
    to_account_id INTEGER REFERENCES accounts(id),
    amount DECIMAL(15,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

##  Security
- All passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- Rate limiting: 100 requests per minute
- All database queries use prepared statements
- Input validation on all endpoints


##  Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Verify PostgreSQL is running
   - Check credentials in .env
   - Ensure database exists

2. **JWT Token Invalid**
   - Check token expiration
   - Verify JWT_SECRET in .env
   - Ensure proper token format in Authorization header

##  License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
