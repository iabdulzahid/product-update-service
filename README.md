# 🧩 Product Update Service

A lightweight **Golang-based event-driven microservice** that asynchronously processes product updates using an in-memory **queue** and a **worker pool**, while maintaining a **thread-safe in-memory store**.

---

## 🚀 Features

- **REST API Endpoints**
  - `POST /events`: Accepts product updates asynchronously.
  - `GET /products/{id}`: Retrieves the latest product state.
- **In-Memory Queue** for event buffering.
- **Worker Pool** (configurable 3–5 workers) for concurrent event processing.
- **Thread-Safe Store** with locking for data integrity.
- **Graceful Shutdown** handling.
- **Comprehensive Tests** for concurrency safety and endpoint behavior.

---

## ⚙️ Setup Instructions

### 1. Clone the Repository

git clone https://github.com/iabdulzahid/product-update-service.git
cd product-update-service

### 2. Install Dependencies
```
go mod tidy
```

### 3. Run the Application
```
go run cmd/server/main.go
```

### 4. Test the Service

Run the unit and concurrency tests:
```
go test ./test -v
```

## 🧱 API Endpoints
### 1. POST /events

Accepts JSON payloads representing product updates.

Example Request
```
curl -X POST http://localhost:8080/events \
-H "Content-Type: application/json" \
-d '{"product_id": "p1", "price": 99.99, "stock": 50}'
```

Example Response
```
{ "message": "Event received and queued successfully" }
```

### 2. GET /products/{id}

Retrieves the latest state of a product.

Example Request
```
curl http://localhost:8080/products/p1
```
Example Response
```
{ "product_id": "p1", "price": 99.99, "stock": 50 }
```
