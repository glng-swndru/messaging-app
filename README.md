# Simple Messaging App

A lightweight messaging application built with Golang using Go Fiber, designed to provide basic chat functionality with user authentication and chat history storage. This project is deployed on AWS and leverages MySQL and MongoDB for persistent data storage.

## 📌 Project Overview

- **Version:** 1.0  
- **Date:** 26-04-2025  
- **Author:** Gilang Swandaru  

## 🧩 System Architecture

![System Architecture](diagram%20messaging%20app.png)

The system is structured as follows:

| Component | Description |
|----------|-------------|
| Website | Frontend interface for users |
| Golang Monolithic App | Backend handling user management and messaging |
| User Management | Subsystem for login, logout, and registration |
| Messaging System | Subsystem for sending and retrieving messages |
| MySQL | Stores user management data |
| MongoDB | Stores chat message history |
| AWS | Hosting platform |

**Connection Details:**
- Website ↔️ Golang App (User): **HTTP**
- Website ↔️ Golang App (Messaging): **WebSocket** (sending), **HTTP** (retrieving)

## 🛠️ Tech Stack

- **Language:** Go
- **Framework:** Go Fiber
- **Databases:** MySQL (users), MongoDB (messages)
- **Cloud:** AWS (EC2, RDS)
- **CI/CD:** GitHub Actions
- **Monitoring & Logging:** ELK Stack (Elasticsearch, Logstash, Kibana)

## ⚙️ CI/CD Pipeline (GitHub Actions)

This project uses GitHub Actions for continuous integration and deployment:

### CI Workflow
- Triggered on every `push` and `pull_request` to `main` branch
- Steps include:
  - Checkout repository
  - Set up Go environment
  - Run `go build` and `go test`
  - Lint and format checks
  - Docker build test

### CD Workflow
- Triggered on successful merge to `main` and new tag release
- Steps include:
  - Build Docker image
  - Push image to Docker Hub: `gilangswndru/simple-messaging-app:latest`
  - Deploy to AWS (e.g., via SSH or ECS)

CI/CD workflows can be found under `.github/workflows/`

## 🔐 API Contract

### Register
- **POST** `/user/v1/register`
- **Body:**
```json
{
  "Username": "string",
  "Password": "string",
  "full_name": "string",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```
- **Response:**
```json
{
  "message": "string",
  "data": {
    "Id": "int",
    "Username": "string",
    "full_name": "string"
  }
}
```

### Login
- **POST** `/user/v1/login`
- **Body:**
```json
{
  "Username": "string",
  "Password": "string"
}
```
- **Response:**
```json
{
  "username": "string",
  "full_name": "string",
  "token": "string",
  "refresh_token": "string"
}
```

### Logout
- **DELETE** `/user/v1/logout`
- **Headers:** `Authorization: JWT Token`
- **Response:**
```json
{
  "Id": "int",
  "Username": "string",
  "full_name": "string"
}
```

### Refresh Token
- **PUT** `/user/v1/refresh-token`
- **Headers:** `Authorization: JWT Refresh Token`
- **Response:**
```json
{
  "message": "string",
  "data": {
    "token": "string"
  }
}
```

### Get Chat History
- **GET** `/message/v1/history&page={1}&limit={10}`
- **Headers:** `Authorization: JWT Token`
- **Response:**
```json
[
  {
    "to": "int",
    "from": "string",
    "date": "datetime",
    "message": "string"
  }
]
```

### Send Chat
- **WebSocket** `/message/v1/send`
- **Message:**
```json
{
  "to": "int",
  "from": "string",
  "date": "datetime",
  "message": "string"
}
```

## 🗃️ Database Schema

### MySQL - Users Table
| Column | Type | Attributes |
|--------|------|------------|
| id | int | Primary key, Auto increment |
| username | varchar(20) | Not null, Default `""` |
| full_name | varchar(100) | Not null, Default `""` |
| password | varchar(255) | Not null, Default `""` |
| created_at | datetime | Default `CURRENT_TIMESTAMP` |
| updated_at | datetime | Default `CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |

### MongoDB - Message History
Example Document:
```json
{
  "_id": "681640f236b66ea30b87bce0",
  "from": "gilang swandaru",
  "message": "test",
  "date": { "timestamp": 1746284109, "i": 87 }
}
```

| Field | Type | Description |
|-------|------|-------------|
| _id | ObjectId | Unique identifier for the document |
| from | string | Sender's name |
| message | string | Message content |
| date | Timestamp | Message timestamp |

## 📊 Observability

- **Monitoring and Logging:** ELK Stack