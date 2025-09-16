# Go Sample Server (Gin + Neo4j)

A RESTful API built with **Go**, **Gin**, and **Neo4j** for managing Persons, Transactions, Operators, Machines, Devices, and related entities. This project demonstrates CRUD operations, relationship management, and filtering using a graph database.

---

## Features

- Create, retrieve, and filter **Person** nodes.
- Create and retrieve **Transaction** nodes (multiple types).
- Manage **Operators**, **Machines**, **Devices**, **EA**, and **Stations**.
- Reset the database for development/testing.
- All endpoints return JSON responses.
- UUID generation for unique identifiers.

---

## Tech Stack

- **Go** (Golang)
- **Gin** web framework
- **Neo4j** Graph Database
- **UUID** for unique identifiers

---

## Prerequisites

- **Go** (≥ 1.18)
- **Neo4j** running locally or remotely
- **Postman** (optional, for API testing)

---

## Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/go-sample-server.git
   cd go-sample-server
   ```

2. **Initialize Go modules**
   ```bash
   go mod tidy
   ```

3. **Install dependencies**
   ```bash
   go get github.com/gin-gonic/gin
   go get github.com/neo4j/neo4j-go-driver/v5/neo4j
   go get github.com/google/uuid
   ```

4. **Configure Neo4j connection**
   Edit `db/db.go`:
   ```go
   uri := "neo4j://localhost:7687"
   username := "neo4j"
   password := "your_password"
   ```

---

## Running the Server

```bash
go run main.go
```
Server runs at: [http://localhost:8080](http://localhost:8080)

---

Running With Docker

````markdown
## Running with Docker

### 1. Build the Docker image
```bash
docker build -t sample-server .
````

This uses the provided `Dockerfile` (multi-stage build, based on Alpine) to produce a lightweight container image.

---

### 2. Run with environment variables

You can pass your `.env` file at runtime (recommended), so it stays outside the image:

```bash
docker run -p 8080:8080 --env-file .env sample-server
```

* `-p 8080:8080` maps container port **8080** to your machine’s port **8080**.
* `--env-file .env` injects all variables from your local `.env` file into the container.
* Your Go app will then read them with `os.Getenv`.

👉 If you don’t have a `.env` file, you can also pass variables directly:

```bash
docker run -p 8080:8080 \
  -e NEO4J_URI=bolt://localhost:7687 \
  -e NEO4J_USER=neo4j \
  -e NEO4J_PASSWORD=yourpassword \
  sample-server
```

---

### 3. Verify the server

Once running, the API will be available at:

[http://localhost:8080](http://localhost:8080)

---

### 4. Stop the container

To stop the running container, press **Ctrl+C** in the terminal or run:

```bash
docker ps   # find container ID
docker stop <container_id>
```

## API Endpoints

### Person Endpoints

- **GET** `/person/`  
  Get all persons.

- **GET** `/person/:uid`  
  Get person by UID.

### Transaction Endpoints

- **GET** `/transaction/`  
  Get all transactions.

- **GET** `/transaction/:id`  
  Get transaction by ID.

- **POST** `/transaction/:type`  
  Create a transaction (type: `type1` or `type2`).

### Database Utility

- **DELETE** `/reset`  
  Reset the database (delete all nodes and relationships).

---

## Example Postman JSON Bodies

### Create Person

```json
{
  "uid": "123e4567-e89b-12d3-a456-426614174000",
  "name": "John Doe",
  "age": 28,
  "gender": "male",
  "contact": "1234567890",
  "location": "London"
}
```

### Create Transaction (type1)

```json
{
  "id": "txn001",
  "person": {
    "uid": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "age": 28,
    "gender": "male",
    "contact": "1234567890",
    "location": "London"
  },
  "ip_address": "192.168.1.1",
  "timestamp": "2025-09-01T10:00:00Z",
  "operator_id": "op001",
  "machine_id": "mach001",
  "devices": [
    { "id": "dev001", "type": "fingerprint" },
    { "id": "dev002", "type": "camera" }
  ],
  "introducer_id": "123e4567-e89b-12d3-a456-426614174999",
  "relation_with_introducer": "friend"
}
```

### Create Transaction (type2)

```json
{
  "id": "txn002",
  "person": {
    "uid": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "age": 28,
    "gender": "male",
    "contact": "1234567890",
    "location": "London"
  },
  "ip_address": "192.168.1.2",
  "timestamp": "2025-09-01T11:00:00Z",
  "operator_id": "op002",
  "machine_id": "mach002",
  "devices": [
    { "id": "dev003", "type": "iris" }
  ]
}
```

---

## Project Structure

```
.
├── db/
│   └── db.go            # Neo4j connection setup
├── models/
│   └── person.go        # Model definitions
├── routes/
│   └── person_routes.go # API route handlers
├── main.go              # Entry point
└── go.mod               # Go module file
```

---

## License

MIT License

---
