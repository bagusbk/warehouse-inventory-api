# Warehouse Inventory Management System API

A RESTful API built with Golang and PostgreSQL for managing warehouse inventory, including product management, stock tracking, purchases, and sales transactions.

## Features

- **Master Data Management**: CRUD operations for products (barang)
- **Stock Management**: Real-time stock tracking with history
- **Purchase Transactions**: Create purchases with automatic stock updates
- **Sales Transactions**: Create sales with stock validation
- **JWT Authentication**: Secure API endpoints
- **Role-based Authorization**: Admin and staff roles
- **Docker Support**: Easy deployment with Docker

## Tech Stack

- **Language**: Go 1.24
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd warehouse-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/warehouse?sslmode=disable"
export SESSION_SECRET="your-secret-key"
export PORT=5000
```

4. Start the server:
```bash
go run main.go
```

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/register | Register new user |
| POST | /api/login | Login and get JWT token |
| GET | /api/profile | Get current user profile |

### Master Data (Barang)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/barang | Get all products with pagination |
| GET | /api/barang/stok | Get products with stock info |
| GET | /api/barang/:id | Get product by ID |
| POST | /api/barang | Create new product |
| PUT | /api/barang/:id | Update product |
| DELETE | /api/barang/:id | Delete product |

### Stock Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/stok | Get all stock |
| GET | /api/stok/:barang_id | Get stock by product ID |
| GET | /api/history-stok | Get stock history |
| GET | /api/history-stok/:barang_id | Get stock history by product |

### Purchase Transactions (Pembelian)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/pembelian | Create purchase |
| GET | /api/pembelian | Get all purchases |
| GET | /api/pembelian/:id | Get purchase by ID |

### Sales Transactions (Penjualan)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/penjualan | Create sale |
| GET | /api/penjualan | Get all sales |
| GET | /api/penjualan/:id | Get sale by ID |

### Reports (Laporan)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/laporan/stok | Stock report |
| GET | /api/laporan/pembelian | Purchase report |
| GET | /api/laporan/penjualan | Sales report |

## API Response Format

### Success Response
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "code": "ERROR_CODE"
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| INSUFFICIENT_STOCK | 400 | Stock is not sufficient |
| ITEM_NOT_FOUND | 404 | Item not found |
| VALIDATION_ERROR | 422 | Input validation error |
| INTERNAL_ERROR | 500 | Server error |
| UNAUTHORIZED | 401 | Authentication required |
| FORBIDDEN | 403 | Access denied |

## Test Users

After seeding, you can use these credentials:

| Username | Password | Role |
|----------|----------|------|
| admin | password123 | admin |
| staff1 | password123 | staff |
| staff2 | password123 | staff |

## Example Requests

### Login
```bash
curl -X POST http://localhost:5000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

### Create Product
```bash
curl -X POST http://localhost:5000/api/barang \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "kode_barang": "BRG006",
    "nama_barang": "USB Flash Drive 32GB",
    "deskripsi": "USB 3.0 Flash Drive",
    "satuan": "pcs",
    "harga_beli": 50000,
    "harga_jual": 75000
  }'
```

### Create Purchase
```bash
curl -X POST http://localhost:5000/api/pembelian \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "supplier": "PT Supplier Baru",
    "details": [
      {"barang_id": 1, "qty": 5, "harga": 15000000},
      {"barang_id": 2, "qty": 20, "harga": 250000}
    ]
  }'
```

### Create Sale
```bash
curl -X POST http://localhost:5000/api/penjualan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "customer": "PT Customer Baru",
    "details": [
      {"barang_id": 1, "qty": 2, "harga": 17500000},
      {"barang_id": 2, "qty": 5, "harga": 350000}
    ]
  }'
```

## License

MIT License
