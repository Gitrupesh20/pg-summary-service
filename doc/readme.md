# Pg Summary Service

A lightweight Go microservice to **summarize PostgreSQL databases** and expose results via REST APIs.

---

## Table of Contents

* [Features](#features)
* [Architecture](#architecture)
* [Prerequisites](#prerequisites)
* [Setup & Run](#setup--run)
* [Configuration](#configuration)
* [API Endpoints](#api-endpoints)
* [Testing](#testing)
* [Potential Improvements](#potential-improvements)

---

## Features

* Fetch summaries from a **local PostgreSQL database**.
* Sync summaries from an **external API**.
* Store external summaries locally for querying.
* RESTful APIs to:
  * Sync summaries (`POST /summary/sync`)
  * Get summaries list (`GET /summaries`)
  * Get summary by ID (`GET /summaries/{id}`)
* Retry mechanism for external API calls.
* Structured logging using **Zap**.
* Dockerized setup for local development.
* Mockable interfaces for **unit testing**.

---

## Architecture

```

Client
|
v
\[HTTP Handler / Middleware]  <--- Logging, etc
|
v
\[SummaryService]  <-- Orchestrates local & external repos
\|          &#x20;
\|           &#x20;
\[LocalRepo]   \[ExternalRepo]  <-- Fetch & store summaries
|
v
PostgreSQL DB

````

* **Service Layer:** Coordinates fetching from external API, validation, and storing to local DB.  
* **Repository Layer:** Abstracted interfaces (`LocalRepository`, `ExternalRepository`) allow easy testing & DI.  
* **Configurable:** Supports switching between local DB or cloud DB via environment variables.  

---

## Prerequisites

* Docker & Docker Compose  
* Go >= 1.24  
* External API running locally or in cloud  

---

## Setup & Run

1. **Navigate to repository**

```bash
cd PgDataSummaryService-AppVersal
````

2. **Set up environment variables**

Environment variables can be loaded from a JSON/YAML config file or set directly:

```env
LOCAL_DB_URL=postgres://postgres:postgres@db:5432/localdb?sslmode=disable
EXTERNAL_API_URL=http://host.docker.internal:3000/api/summary
```

> **Notes:**
>
> * If you don’t set `LOCAL_DB_URL`, the service will default to the bundled Postgres container (`db` in docker-compose).
> * `EXTERNAL_API_URL` must be reachable from Docker. Use `host.docker.internal` for APIs running on your host machine.

3. **Build and run with Docker**

```bash
docker-compose up --build
```

* App will be available on `http://localhost:8080`
* Local Postgres will be available on `localhost:5432`

---

## Configuration

Configuration can be provided through environment variables or generated from **YAML/JSON** config files.

Example Go config struct:

```go
conf := config{
    port:                  8080,
    localDbUrl:            "postgres://postgres:postgres@db:5432/localdb?sslmode=disable",
    maxConnections:        5,
    maxIdleConnections:    2,
    maxConnectionLifeTime: 5 * time.Minute,
    externalDbUrl:         "http://host.docker.internal:3000/api/summary",
    retries:               3,
    isDebug:               false,
    logDir:                "./logs",
    logFile:               "server.log",
}
```
> Note: For now, I’m not loading from JSON or YAML; I’m just hardcoding some values in the code.
---

## API Endpoints

### 1. Sync Summaries

**POST** `/summary/sync`

**Request Body:**

```json
{
  "host": "aaaaa-db.example.com",
  "port": 5432,
  "user": "readonly",
  "password": "pass",
  "dbname": "sample"
}
```

**Response:**

```json
{
  "status": "sync triggered"
}
```

---

### 2. Get Summaries List

**GET** `/summaries`

**Response:**

```json
[
  {
    "id": "sum-1757835089142",
    "db_name": "aaaaa-db.example.com:sample",
    "synced_at": "2025-09-14T07:31:29.737242Z"
  }
]
```

---

### 3. Get Summary by ID

**GET** `/summaries/{id}`

**Response:**

```json
{
  "summary_id": "sum-1757835089142",
  "source": "aaaaa-db.example.com:sample",
  "synced_at": "2025-09-14T07:31:29.737242Z",
  "schemas": [
    {
      "id": "41ab450e-b615-4b08-9ceb-288ece05c063",
      "name": "sales",
      "table_count": 2,
      "total_rows": 450,
      "total_size_mb": 7.5
    },
    {
      "id": "4a695036-0950-459f-8f4a-fdc80e0d3926",
      "name": "public",
      "table_count": 2,
      "total_rows": 1820,
      "total_size_mb": 20.8
    }
  ]
}
```

---

## Testing

Unit tests are written using `stretchr/testify` with **mocked repos**.

Run tests locally:

```bash
go test ./... -v
```

---

## Potential Improvements

1. **Worker Pool**

    * Fetch summaries concurrently for multiple databases using a configurable worker pool.
    * Reduces latency when syncing multiple remote databases.

2. **TTL Cache**

    * Implement cache with TTL and LRU for frequently requested summaries to reduce DB load.

3. **Authentication / Authorization**

    * Secure APIs using token-based or basic auth with role-based access control.

```

