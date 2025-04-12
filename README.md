[![CI](https://github.com/imhasandl/message-service/actions/workflows/ci.yml/badge.svg)](https://github.com/imhasandl/message-service/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/imhasandl/message-service)](https://goreportcard.com/report/github.com/imhasandl/message-service)
[![GoDoc](https://godoc.org/github.com/imhasandl/message-service?status.svg)](https://godoc.org/github.com/imhasandl/message-service)
[![Coverage](https://codecov.io/gh/imhasandl/message-service/branch/main/graph/badge.svg)](https://codecov.io/gh/imhasandl/message-service)
[![Go Version](https://img.shields.io/github/go-mod/go-version/imhasandl/message-service)](https://golang.org/doc/devel/release.html)

# Message Service

A microservice for handles sending messages between users in a social media application, built with Go and gRPC.

## Overview

The Message Service is a microservice designed to facilitate communication between users in a social media application. It provides robust and scalable functionality for sending and receiving messages, ensuring reliable interaction between users. Built with Go and leveraging gRPC for efficient communication, this service is optimized for performance and ease of integration with other services within the application ecosystem.

## Prerequisites

- Go 1.20 or later
- PostgreSQL database
- RabbitMQ

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
PORT=":YOUR_GRPC_PORT"
DB_URL="postgres://username:password@host:port/database?sslmode=disable"
# DB_URL="postgres://username:password@db:port/database?sslmode=disable" // FOR DOCKER COMPOSE
TOKEN_SECRET="YOUR_JWT_SECRET_KEY"
RABBITMQ_URL="amqp://username:password@host:port"
```

> **Note:** Make sure that you use same token secret in every services

## Database Migrations

This service uses Goose for database migrations:

```bash
# Install Goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir migrations postgres "YOUR_DB_CONNECTION_STRING" up
```

## gRPC Methods

The service implements the following gRPC methods:

---

### SendMessage

Insert a message to user into a messages table and stores that in database

#### Request format

```json
{
  "receiver_id": "id of a person who's gonna recieve the message",
  "content": "message content"
}
```

> **Note:** After the message is sent successfully, it sends a to message broker a json message and sends that to a user that received the message via push notification

#### Response format

```json
{
  "success": "a bool value that determines the result of of the querie TRUE if successfully completed, False otherwise"
}
```

---

### GetMessages

Gets the messages between users. In this case the current user, who wants to see the messages between himself and different person. We are gonna get current user's id using context, retrieve a token from it, and use it as current user id, and then we are gonna recieve a receiver_id's from client side then proceed the method.

#### Request format

```json
{
  "receiver_id": "gets messages of two users from db"
}
```

#### Response format

```json
{
  "message": [
    {
      "id": "string",
      "sent_at": "2025-04-11T19:44:23Z",
      "sender_id": "string",
      "receiver_id": "string",
      "content": "string"
    }
  ]

```

---

### ChangeMessage

Changes message content in database, using id of a message.

#### Request format

```json
{
  "id": "UUID of a message",
  "content": "New content"
}
```

#### Response format

```json
{
  "message": {
    "id": "string",
    "sent_at": "2025-04-11T19:44:23Z",
    "sender_id": "string",
    "receiver_id": "string",
    "content": "new content"
  }
}
```

---

### DeleteMessage

Deletes message from db using incoming message id

#### Request format

```json
{
  "id": "UUID string of a message"
}
```

#### Response format

```json
{
  "status": "boolean value for the result if TRUE the message is delted successfully FALSE otherwise"
}
```

---

## RabbitMQ Integration

The Notification Service consumes messages from RabbitMQ to process asynchronous notification requests from other services.

### Message Consumption

The service automatically sets up and listens to:
- **Exchange**: `notifications.topic` (topic exchange)
- **Queue**: `notification_service_queue`
- **Routing Key**: `#` (wildcard - receives all messages published to the exchange)

### Publishing Messages to the Notification Service

Other microservices can send notification requests by publishing messages to the `notifications.topic` exchange. Messages should be JSON formatted with the following structure:

```json
{
   "title": "Notification title",
   "sender_id": "username of sender",
   "receiver_id": "UUID of recipient user",
   "content": "Notification message content", 
   "sent_at": "2023-01-01T12:00:00Z"
}
```

## Running the Service

```bash
go run cmd/main.go
```

## Docker Support

The service can be run as part of a Docker Compose setup along with other microservices. When using Docker, make sure to use the Docker Compose specific DB_URL configuration.