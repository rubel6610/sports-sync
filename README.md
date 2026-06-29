# SpotSync

A modern parking management system API built with Go, providing comprehensive parking reservation and zone management capabilities.

**Live Demo:** https://sports-sync.onrender.com/

---

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [Setup & Installation](#setup--installation)
- [Environment Configuration](#environment-configuration)
- [Database Schema](#database-schema)
- [Authentication](#authentication)
- [Usage Examples](#usage-examples)

---

## ✨ Features

- **User Authentication** - Secure registration and login with JWT tokens
- **Role-Based Access Control** - Driver and Admin roles with permission enforcement
- **Parking Zone Management** - Create, retrieve, and manage parking zones with capacity tracking
- **Real-Time Availability** - Live parking spot availability calculation
- **Reservation System** - Create, view, and cancel parking reservations
- **Atomic Transactions** - Race condition prevention with row-level locking
- **Connection Pooling** - Optimized database connection management
- **Error Handling** - Standardized error responses with proper HTTP status codes
- **API Documentation** - Home endpoint providing complete API reference

---

## 🛠 Tech Stack

- **Language:** Go 1.26.2
- **Framework:** Echo v4.15.4 (HTTP web framework)
- **Database:** PostgreSQL (NeonDB)
- **ORM:** GORM v1.31.2
- **Authentication:** JWT (golang-jwt/jwt/v5)
- **Password Hashing:** bcrypt
- **Environment Config:** godotenv

---

## 📁 Project Structure

```
spotsync-api/
├── main.go                    # Application entry point & routing
├── go.mod                     # Go module dependencies
├── config/
│   └── db.go                 # Database connection & pooling config
├── models/
│   └── models.go             # User, ParkingZone, Reservation entities
├── dto/
│   └── dto.go                # Request/response data structures
├── handler/
│   ├── auth_handler.go       # Authentication endpoints
│   ├── zone_handler.go       # Parking zone endpoints
│   ├── reservation_handler.go # Reservation endpoints
│   ├── home_handler.go       # API documentation endpoint
│   └── middleware.go         # JWT & role-based middleware
├── service/
│   ├── auth_service.go       # Authentication business logic
│   ├── zone_service.go       # Parking zone business logic
│   └── reservation_service.go # Reservation business logic
├── repository/
│   ├── user_repository.go    # User database operations
│   ├── zone_repository.go    # Parking zone database operations
│   └── reservation_repository.go # Reservation database operations
└── utils/
    ├── jwt.go                # JWT token operations
    └── errors.go             # Error handling utilities
```

---

## 🔌 API Endpoints

### Home (Documentation)

```
GET /
```

Returns API overview and complete endpoint documentation.

### Authentication

| Method | Endpoint                | Description           | Auth |
| ------ | ----------------------- | --------------------- | ---- |
| POST   | `/api/v1/auth/register` | Register new user     | None |
| POST   | `/api/v1/auth/login`    | Login & get JWT token | None |

### Parking Zones

| Method | Endpoint            | Description                     | Auth        |
| ------ | ------------------- | ------------------------------- | ----------- |
| POST   | `/api/v1/zones`     | Create new parking zone         | JWT + Admin |
| GET    | `/api/v1/zones`     | Get all zones with availability | None        |
| GET    | `/api/v1/zones/:id` | Get specific zone details       | None        |

### Reservations

| Method | Endpoint                               | Description                  | Auth        |
| ------ | -------------------------------------- | ---------------------------- | ----------- |
| POST   | `/api/v1/reservations`                 | Create new reservation       | JWT         |
| GET    | `/api/v1/reservations/my-reservations` | Get user's reservations      | JWT         |
| DELETE | `/api/v1/reservations/:id`             | Cancel reservation           | JWT         |
| GET    | `/api/v1/reservations`                 | Get all reservations (admin) | JWT + Admin |

---

## 🚀 Setup & Installation

### Prerequisites

- Go 1.26.2 or higher
- PostgreSQL database (NeonDB recommended)
- Git

### Steps

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd spotsync-api
   ```

2. **Install dependencies**

   ```bash
   go mod download
   go mod tidy
   ```

3. **Configure environment variables** (create `.env` file)

   ```bash
   PORT=8080
   DB_HOST=your-neondb-host
   DB_USER=your-db-user
   DB_PASSWORD=your-db-password
   DB_NAME=your-db-name
   DB_PORT=5432
   JWT_SECRET=your-secret-key
   ```

4. **Run the server**
   ```bash
   go run main.go
   ```
   Server starts on `http://localhost:8080`

---

## 🔐 Environment Configuration

Create a `.env` file in the project root with the following variables:

```env
# Server
PORT=8080

# Database (NeonDB PostgreSQL)
DB_HOST=pg-xxxxx.neon.tech
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=neondb
DB_PORT=5432

# JWT Authentication
JWT_SECRET=your-secret-key-min-32-chars
```

**Security Note:** Never commit `.env` to version control. Use environment variables in production.

---

## 🗄 Database Schema

### Users Table

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(50) DEFAULT 'driver',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Parking Zones Table

```sql
CREATE TABLE parking_zones (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(50),
  total_capacity INT,
  price_per_hour DECIMAL(10, 2),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Reservations Table

```sql
CREATE TABLE reservations (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
  zone_id INT REFERENCES parking_zones(id) ON DELETE RESTRICT ON UPDATE CASCADE,
  license_plate VARCHAR(50),
  status VARCHAR(50) DEFAULT 'active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Relationships:**

- User → Reservations (One-to-Many)
- ParkingZone → Reservations (One-to-Many)
- Cascade updates, restrict deletes

---

## 🔑 Authentication

### JWT Token Structure

The API uses JWT for stateless authentication.

**Token Claims:**

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "role": "driver",
  "exp": 1234567890
}
```

**Token Expiration:** 24 hours

### Usage

1. Register or login to get JWT token
2. Include token in request header: `Authorization: Bearer <token>`
3. Server validates token and authorizes requests

### Role-Based Access Control

- **Driver:** Can create & manage own reservations, view zones
- **Admin:** Can create zones, view all reservations

---

## 💡 Usage Examples

### 1. Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

Response includes JWT token to use for authenticated endpoints.

### 3. Create Parking Zone (Admin Only)

```bash
curl -X POST http://localhost:8080/api/v1/zones \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "name": "Downtown Parking Lot A",
    "type": "general",
    "total_capacity": 100,
    "price_per_hour": 5.00
  }'
```

### 4. Get All Parking Zones

```bash
curl http://localhost:8080/api/v1/zones
```

Returns zones with calculated available spots.

### 5. Create Reservation

```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "zone_id": 1,
    "license_plate": "ABC-1234"
  }'
```

### 6. Get My Reservations

```bash
curl http://localhost:8080/api/v1/reservations/my-reservations \
  -H "Authorization: Bearer <your-jwt-token>"
```

### 7. Cancel Reservation

```bash
curl -X DELETE http://localhost:8080/api/v1/reservations/1 \
  -H "Authorization: Bearer <your-jwt-token>"
```

### 8. View API Documentation

```bash
curl http://localhost:8080/
```

---

## 📊 Response Format

All API responses follow a standardized format:

**Success Response:**

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "role": "driver"
  }
}
```

**Error Response:**

```json
{
  "success": false,
  "message": "Invalid credentials",
  "data": null
}
```

---

## 🔒 Security Features

- **Password Hashing:** bcrypt with salt
- **JWT Authentication:** Stateless token-based auth
- **Role-Based Access Control:** Middleware enforces permissions
- **Race Condition Prevention:** Atomic transactions with row-level locking
- **Input Validation:** Struct tags for request validation
- **Error Masking:** Internal errors masked in responses
- **CORS Support:** Cross-origin requests configured

---

## 🐛 Common Issues

### Port Already in Use

```bash
# Use different port
PORT=3000 go run main.go
```

### Database Connection Timeout

- Verify NeonDB credentials in `.env`
- Check network connectivity to database
- Ensure SSL mode is enabled

### Slow Query Warnings

These appear during `AutoMigrate()` due to cloud database latency. They don't affect functionality.

---

## 📝 License

This project is part of Level 2 Assignment 6.

---

## 👨‍💻 Author

SpotSync API Development Team

---

## 📞 Support

For issues or questions, please refer to the project documentation or contact the development team.

---

**Last Updated:** June 29, 2026  
**Version:** 1.0.0
