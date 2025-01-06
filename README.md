# Simple Thread-Safe Key-Value Store in Go

This project implements a basic but thread-safe key-value store using Go. It provides persistent storage with transaction logging capabilities.

## Features

- Thread-safe operations for concurrent access
- Persistent storage through transaction logging
- Support for basic Put and Delete operations
- Transaction log replay for data recovery
- Error handling and reporting through channels

## Implementation Details

### Transaction Logger

The store uses a transaction logging system to maintain data persistence:

- Each operation (Put/Delete) is logged as an event with a sequence number
- Events are written to a transaction log file
- On startup, the log can be replayed to recover the store's state
- The logger runs in a separate goroutine for non-blocking operations

### Data Structure

Events in the transaction log follow this format:
SequenceNumber	EventType	Key	"Value"

## Running

### Prerequisites

- Go 1.16 or higher
- Git

### Steps to Build and Run

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd kvstore/api/main
   ```

2. Run the server:
   ```bash
   go run server.go
   ```

### Using the Key-Value Store

Once running, you can interact with the store using HTTP requests:

- Put a value:
  ```bash
  curl -X PUT "localhost:8080/v1/key/hello" -d "hello world"
  ```

- Get a value:
  ```bash
  curl "localhost:8080/v1/key/hello"
  ```

- Delete a value:
  ```bash
  curl -X DELETE "localhost:8080/v1/key/hello"
  ```
