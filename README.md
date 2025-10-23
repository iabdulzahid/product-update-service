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
go run cmd/product-update-service/main.go
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

## 🧠 Design Choices

### 1. **Event-Driven Architecture**
- The core concept is to **decouple producers (incoming requests)** from **consumers (update processors)**.
- Incoming product updates are pushed into an in-memory **`EventQueue`** (`chan *domain.Product`).
- A background worker continuously dequeues and applies updates to the **`ProductStore`**.

### 2. **In-Memory Data Structures**
- **`ProductStore`** uses a `sync.Map` for thread-safe concurrent reads and writes.
- The in-memory store simplifies implementation and testing while maintaining concurrency safety.

### 3. **Concurrency Control**
- The `EventQueue` and `ProductStore` are **thread-safe**.
- Tests simulate both **sequential** and **concurrent** product updates to validate safe behavior.

### 4. **Loose Coupling Between Layers**
- Separation between:
  - `handler` → HTTP layer
  - `queue` → Asynchronous event processing
  - `repository` → Data persistence (currently in-memory)
- This separation makes it easier to replace components (e.g., switch from in-memory to Redis).

---

## 🏭 Production Considerations

In production, you would likely evolve this design into a **distributed, fault-tolerant event-driven service**.

### 1. **Message Broker (RabbitMQ / Kafka)**
- Replace the in-memory `EventQueue` with a **persistent message queue**.
- This ensures events aren’t lost during crashes or restarts.
- Use **acknowledgment** and **retry mechanisms** to handle failed message deliveries.

### 2. **Persistent Storage (PostgreSQL / Redis)**
- Replace `ProductStore` with:
  - **PostgreSQL** for durable storage and complex queries.
  - **Redis** for high-speed access if data is primarily key-value-based.
- Include migrations and schema versioning tools like `golang-migrate`.

### 3. **High Throughput Handling**
- Scale horizontally by deploying multiple consumers reading from the same queue.
- Use **batch updates** or **worker pools** for efficient database writes.
- Apply **rate limiting** and **backpressure** on producers to prevent overload.

### 4. **Error Handling & Retry Mechanisms**
- Implement retry with exponential backoff for failed database writes.
- Use a **dead-letter queue (DLQ)** for failed events after multiple retries.
- Introduce structured logging and monitoring via Prometheus + Grafana.

---

## 🩺 Troubleshooting Strategies

### 1. **Data Consistency Problems**
**Symptoms:**
- Product data mismatched between requests and store.

**Possible Causes:**
- Race conditions in concurrent updates.
- Partial event processing or dropped messages.

**Approach:**
- Enable detailed logging for enqueue/dequeue operations.
- Verify message acknowledgment logic (if using RabbitMQ/Kafka).
- Add tests simulating concurrent writes with same `ProductID`.

---

### 2. **Products Aren’t Updating Despite Events Being Received**
**Possible Causes:**
- Consumer goroutine not running or blocked.
- Channel capacity full (deadlock or backpressure issue).
- Marshaling/unmarshaling errors in payload.

**Debugging Steps:**
1. Check if the consumer goroutine (`for p := range eq.Dequeue()`) is active.
2. Log each enqueue and dequeue operation.
3. Use `TryEnqueue()` and monitor when queue hits capacity.
4. Temporarily increase channel buffer size for testing.
5. Add `pprof` or runtime metrics to monitor goroutines.

---

## 🧪 Running Tests

Run all tests to validate functionality and concurrency safety:
```bash
go test ./test -v
```

The test suite covers:
- Basic POST + GET flow  
- Sequential updates for the same product  
- Concurrent updates for multiple products  
- Handling of non-existent products  

---

## 🧰 Future Enhancements
- ✅ Add structured logging (zap/logrus)
- ✅ Add metrics endpoint for queue depth and update rate
- ✅ Introduce middleware for tracing (OpenTelemetry)
- ✅ Implement configurable retry and DLQ
- ✅ Support for distributed queue (RabbitMQ/Kafka)
- ✅ Integrate PostgreSQL with transaction management

---


