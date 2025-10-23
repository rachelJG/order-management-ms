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
* File .env configured with the environment variables. (Example: .env.example)

### 2. Clone the repository

```bash
git clone https://github.com/rachelJG/order-management-ms.git
cd order-management-ms
```
Copy the environment variables file:

```bash
cp .env.example .env
```

Then open .env and adjust the variables as needed (ports, credentials, topic names, etc).

### 3. Start the services

```bash
docker-compose up -d
```

This will start all the containers defined in the docker-compose.yml file (for example: API, database, Kafka, Zookeeper, etc).

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
| **Language**     | Go 1.23.0                       | Simplicity, concurrency, microservice-friendly   |
| **Framework**    | Gin                       | Lightweight and idiomatic HTTP routers           |
| **Database**     | MongoDB                   | Flexible NoSQL document model for dynamic orders |
| **Cache**        | Redis                     | Improves read performance and reduces DB load    |
| **Messaging**    | Kafka (segmentio/kafka-go)| Enables asynchronous and decoupled communication |
| **Architecture** | Modular                   | Separation of concerns and testability           |
| **Config**       | `.env` + `config` package | Centralized environment configuration            |
| **Logging**      | zap                       | Structured, leveled logging                      |
| **Testing**      | testify                   | Unit testing and mocking support                 |


---

## 🧩 Future Improvements

* Add authentication and authorization
* Add metrics and tracing
* Configure CI/CD pipeline
* Add Validation layer
* Add swagger documentation
* Improve error messages

---

## 🧑‍💻 Author

**Raquel Garcia**
Software Engineer | Microservices & Go Developer
