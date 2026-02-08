# TODO - CLI and API Server for Task Management

## Project Description

Full-featured task management system (TODO list) in Go with client-server architecture. The project consists of a CLI application for users, REST API server, and client library for interacting with the server.

Educational project for learning Go, REST API, Cobra CLI, Viper configuration, and client-server architecture.

## Features

- 🚀 **CLI Interface** - convenient task management through terminal
- 🌐 **REST API Server** - backend for storing and processing tasks
- 📦 **Client Library** - Go library for API interaction
- 💾 **JSON Storage** - data stored in JSON format
- ⚙️ **Viper Configuration** - flexible settings management
- 🎯 **CRUD Operations** - complete Create, Read, Update, Delete operations
- 🔧 **Cobra Commands** - professional CLI structure

## Technology Stack

**Programming Language:**
- Go 1.16+

**Libraries and Frameworks:**
- **Cobra** - CLI framework for creating commands
- **Viper** - configuration management
- **Encoding/JSON** - data serialization
- **net/http** - HTTP server and client
- **Gorilla/mux** (possibly) - HTTP routing

## Project Structure

```
apis/
├── todo/              # CLI application
│   ├── cmd/          # Cobra commands
│   ├── main.go       # CLI entry point
│   └── config/       # Configuration files
│
├── todo_server/      # REST API server
│   ├── main.go       # Server entry point
│   ├── handlers/     # HTTP handlers
│   ├── models/       # Data structures
│   └── storage/      # JSON storage
│
├── todo_client/      # Client library
│   ├── client.go     # HTTP client
│   └── models.go     # Data models
│
├── Makefile          # Build scripts
└── go.mod            # Go modules
```

## System Components

### 1. TODO CLI (todo/)

Console application for task management:

**Commands:**
- `todo add "task description"` - add new task
- `todo list` - show all tasks
- `todo complete <id>` - mark task as completed
- `todo delete <id>` - delete task
- `todo update <id> "new description"` - update task description

**Examples:**
```bash
# Add task
./todo add "Buy milk"

# View all tasks
./todo list

# Complete task
./todo complete 1

# Delete task
./todo delete 2
```

### 2. TODO Server (todo_server/)

REST API server for task management:

**Endpoints:**
- `GET /todos` - get all tasks
- `GET /todos/:id` - get task by ID
- `POST /todos` - create new task
- `PUT /todos/:id` - update task
- `DELETE /todos/:id` - delete task

**Data Format:**
```json
{
  "id": 1,
  "title": "Buy milk",
  "completed": false,
  "created_at": "2024-01-15T10:30:00Z"
}
```

### 3. TODO Client (todo_client/)

Go library for interacting with TODO server:

**Usage:**
```go
import "github.com/PopovMarko/apis/todo_client"

client := todoclient.New("http://localhost:8080")

// Create task
task, err := client.CreateTodo("New task")

// Get all tasks
todos, err := client.GetTodos()

// Update task
err = client.UpdateTodo(1, "Updated text")

// Delete task
err = client.DeleteTodo(1)
```

## Installation and Setup

### Requirements

- Go 1.16 or higher
- Git

### Step 1: Clone Repository

```bash
git clone https://github.com/PopovMarko/apis.git
cd apis
```

### Step 2: Install Dependencies

```bash
go mod download
```

### Step 3A: Run Server

```bash
# Compile server
cd todo_server
go build -o todo_server

# Run server
./todo_server

# Or via Make (if Makefile exists)
make server
```

Server will start by default at `http://localhost:8080`

### Step 3B: Run CLI Application

**Option 1: Via go run**
```bash
cd todo
go run main.go list
```

**Option 2: Build binary**
```bash
cd todo
go build -o todo
./todo list
```

**Option 3: Via Make**
```bash
make build-cli
./todo list
```

### Step 4: Configuration

Create configuration file `config.yaml`:

```yaml
server:
  host: localhost
  port: 8080
  
storage:
  type: json
  path: ./data/todos.json

cli:
  server_url: http://localhost:8080
```

Viper will automatically load the configuration.

## Usage

### CLI Mode

```bash
# Add task
./todo add "Write documentation"

# View all tasks
./todo list

# Complete task #1
./todo complete 1

# Update task #2
./todo update 2 "Updated documentation"

# Delete task #3
./todo delete 3

# Show help
./todo --help
```

### API Mode

**Create task:**
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "New task"}'
```

**Get all tasks:**
```bash
curl http://localhost:8080/todos
```

**Get task by ID:**
```bash
curl http://localhost:8080/todos/1
```

**Update task:**
```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated task", "completed": true}'
```

**Delete task:**
```bash
curl -X DELETE http://localhost:8080/todos/1
```

## Architecture

### Client-Server Model

```
┌─────────────┐         HTTP/REST         ┌──────────────┐
│   TODO CLI  │ ◄─────────────────────► │ TODO Server  │
│   (Client)  │                           │   (Backend)  │
└─────────────┘                           └──────────────┘
      │                                          │
      │                                          │
      │ Uses                                     │ Stores
      ▼                                          ▼
┌─────────────┐                           ┌──────────────┐
│ TODO Client │                           │  todos.json  │
│  (Library)  │                           │   (Storage)  │
└─────────────┘                           └──────────────┘
```

### Cobra CLI Structure

```
todo
├── add       # Add task
├── list      # Show all
├── complete  # Mark as completed
├── delete    # Delete
├── update    # Update
└── help      # Help
```

### Viper Configuration

Viper supports:
- Configuration files (YAML, JSON, TOML)
- Environment variables
- Command-line flags
- Default values

**Priority:** CLI flags > ENV vars > Config file > Defaults

## Development

### Adding New Command

1. **Create new command file:**
```go
// todo/cmd/search.go
package cmd

import (
    "github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
    Use:   "search [keyword]",
    Short: "Search todos by keyword",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(searchCmd)
}
```

2. **Build project:**
```bash
go build -o todo
```

### Adding New API Endpoint

1. **Add handler:**
```go
// todo_server/handlers/search.go
func SearchTodosHandler(w http.ResponseWriter, r *http.Request) {
    keyword := r.URL.Query().Get("q")
    // Search logic
}
```

2. **Register route:**
```go
// todo_server/main.go
http.HandleFunc("/todos/search", handlers.SearchTodosHandler)
```

### Testing

```bash
# Test CLI
cd todo
go test ./...

# Test server
cd todo_server
go test ./...

# Test client
cd todo_client
go test ./...
```

## Build and Deploy

### Cross-Platform Compilation

**Linux:**
```bash
GOOS=linux GOARCH=amd64 go build -o todo-linux
```

**Windows:**
```bash
GOOS=windows GOARCH=amd64 go build -o todo.exe
```

**macOS:**
```bash
GOOS=darwin GOARCH=amd64 go build -o todo-mac
```

### Docker (if needed)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o todo_server ./todo_server

FROM alpine:latest
COPY --from=builder /app/todo_server /todo_server
EXPOSE 8080
CMD ["/todo_server"]
```

**Build and run:**
```bash
docker build -t todo-server .
docker run -p 8080:8080 todo-server
```

## Libraries Used

### Cobra
**Purpose:** CLI framework  
**Usage:** Creating commands and subcommands  
**Official site:** https://github.com/spf13/cobra

### Viper
**Purpose:** Configuration management  
**Usage:** Reading configs from files, ENV, flags  
**Official site:** https://github.com/spf13/viper

### Encoding/JSON
**Purpose:** Data serialization  
**Usage:** Saving and reading tasks in JSON  
**Package:** Go standard library

## Technical Features

### CRUD Operations
- **Create** - POST /todos
- **Read** - GET /todos, GET /todos/:id
- **Update** - PUT /todos/:id
- **Delete** - DELETE /todos/:id

### JSON Storage
```json
[
  {
    "id": 1,
    "title": "Task 1",
    "completed": false,
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "title": "Task 2",
    "completed": true,
    "created_at": "2024-01-15T11:00:00Z"
  }
]
```

### HTTP Client/Server
- RESTful API design
- JSON request/response
- Error handling
- Status codes

## Future Improvements

- [ ] Add filters (completed, pending)
- [ ] Task priorities (high, medium, low)
- [ ] Due dates
- [ ] Categories/tags
- [ ] Task search
- [ ] Database instead of JSON (PostgreSQL/SQLite)
- [ ] User authentication
- [ ] Web UI (React/Vue)
- [ ] gRPC instead of REST
- [ ] GraphQL API
- [ ] Docker Compose for deployment

## Troubleshooting

**Server won't start:**
```bash
# Check if port is free
lsof -i :8080
# Or change port in configuration
```

**CLI can't connect to server:**
```bash
# Check if server is running
curl http://localhost:8080/todos
# Check server_url in configuration
```

**"command not found" error:**
```bash
# Make sure binary is compiled
go build -o todo
# Or use full path
./todo list
```

## Learning Objectives

This project was created to learn:
- **Go CLI development** with Cobra
- **REST API** creation and design
- **Client-Server** architecture
- **JSON** serialization/deserialization
- **HTTP** protocol
- **Configuration management** with Viper
- **Go modules** and package management
- **CRUD operations** in practice

## Author

Popov Marko Vyacheslavovych

## License

Educational project

---

**Note:** This is an educational project for learning Go, REST API, and CLI development. Use for learning and practice!
