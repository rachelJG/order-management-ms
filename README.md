# 🧩 Order Management Service

**Order Management Service** is a Go-based microservice that manages the lifecycle of delivery orders within a logistics platform.
It allows registering new orders, querying orders by client or status, and updating their delivery state.
The service leverages **MongoDB**, **Redis**, and **Kafka** to implement persistence, caching, and asynchronous event-driven communication.

---

## 🚀 Features

* **Register new orders**
* **Query orders** by client or status
* **Update order state** and emit asynchronous events
* **Cache frequent queries** using Redis
* **Persist data** in MongoDB
* **Publish domain events** to Kafka when an order changes state

---

## 🧠 Architecture Overview

This microservice follows clean architecture principles and includes:

* REST API for communication (`gin`)
* MongoDB repository for order persistence
* Redis cache layer for optimized reads
* Kafka producer for event publishing (In topic "order_events" when order state changes)
* Configurable environment variables

---

## ⚙️ Running Locally

### 1. Prerequisites

* Go 1.23+
* Docker & Docker Compose

### 2. Clone the repository

```bash
git clone https://github.com/your-org/order-management-ms.git
cd order-management-ms
```

### 3. Start dependencies (MongoDB, Redis, Kafka)

```bash
docker-compose up -d
```

### 4. Run the service

```bash
go run ./cmd/api
```

The service will start at `http://localhost:8080`

---

## 🧪 API Examples

### Create a new order

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
        "client_id": "client-123",
        "items": [
          {"sku": "P001", "name": "Package A", "quantity": 2, "price": 50.0}
        ]
      }'
```

### Query orders by client and status

```bash
curl "http://localhost:8080/orders?client_id=client-123&state=created"
```

### Update order state

```bash
curl -X PUT http://localhost:8080/orders/uuid-123/state \
  -H "Content-Type: application/json" \
  -d '{"new_state": "in_transit"}'
```

---

## 🧰 Technical Decisions

| Area             | Decision                  | Reason                                           |
| ---------------- | ------------------------- | ------------------------------------------------ |
| **Language**     | Go                        | Simplicity, concurrency, microservice-friendly   |
| **Framework**    | Gin                       | Lightweight and idiomatic HTTP routers           |
| **Database**     | MongoDB                   | Flexible NoSQL document model for dynamic orders |
| **Cache**        | Redis                     | Improves read performance and reduces DB load    |
| **Messaging**    | Kafka                     | Enables asynchronous and decoupled communication |
| **Architecture** | Modular                   | Separation of concerns and testability           |
| **Config**       | `.env` + `config` package | Centralized environment configuration            |
| **Logging**      | logrus                    | Structured, leveled logging                      |
| **Testing**      | testify                   | Unit testing and mocking support                 |

---

## 🧩 Future Improvements

* Add authentication and authorization
* Include metrics (Prometheus) and tracing (OpenTelemetry)
* Add pagination and sorting for query endpoints

---

## 🧑‍💻 Author

**Raquel Garcia**
Software Engineer | Microservices & Go Developer
