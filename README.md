# 🧩 Order Management MicroService

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

### 3. Build and start the services

```bash
docker-compose up --build -d
```
This will:
- Build all Docker images based on your local Dockerfile
- Start all services defined in docker-compose.yml (for example: API, database, Kafka, Zookeeper, etc.)

The service will start at `http://localhost:8080`

### 4. Stop the services

```bash
docker-compose down
```

### 5. Check running containers
```bash
docker ps
```

### 6. Check logs
```bash
docker logs <container_id>
```

---

## 🧪 API Examples

### Create a new order

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "1233",
    "items": [
      {
        "sku": "JNS-CLS-32",
        "quantity": 1
      }
    ]
  }'

```

### Get order by id

```bash
curl -X GET http://localhost:8080/api/v1/orders/ORD-e7825df7
```

### Query orders by client and status

```bash
curl -X GET "http://localhost:8080/api/v1/orders?status=NEW&page=1&limit=10"
```

### Update order state

```bash
curl -X PATCH http://localhost:8080/api/v1/orders/ORD-12b72b69/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "IN_PROGRESS"
  }'

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

### Project Structure

The project follows a clean and modular architecture to maintain separation of concerns and scalability.

```bash
src/
└── main/
    ├── config/
    ├── controllers/
    ├── models/
    ├── pkg/
    ├── repositories/
    ├── services/
    └── main.go
``` 

# config/

Contains configuration files and initialization logic for external dependencies such as environment variables, database connections, Kafka, etc.

# controllers/

You define the HTTP handlers that receive incoming requests, validate input data, and call the corresponding service layer functions.
They represent the presentation layer of the application.

# models/

Includes the data models and structures used across the project, such as database entities and API request/response models.

# pkg/

Holds reusable helper packages and utilities (e.g., logging, Kafka producer/consumer, middlewares).
This folder can be imported by other layers when needed.

# repositories/

Implements the data access logic.
Each repository interacts directly with the database or external storage using libraries like GORM or MongoDB drivers.

# services/

Contains the business logic of the application.
Services act as intermediaries between controllers and repositories, encapsulating use case–specific operations.

# main.go

The entry point of the application.
It loads the configuration, initializes dependencies, sets up routes, and starts the HTTP server.



## 🧩 Future Improvements

* Add authentication and authorization
* Add metrics and tracing
* Configure CI/CD pipeline
* Add Validation layer
* Add swagger documentation
* Improve error messages
* Add unit tests

---

## 🧑‍💻 Author

**Raquel Garcia**
Software Engineer | Microservices & Go Developer
