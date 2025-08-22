# ScyllaDB CRUD Application

A complete Go application demonstrating CRUD (Create, Read, Update, Delete) operations with ScyllaDB using the gocqlx library.

## Features

- ✅ **Database Initialization**: Automatically creates keyspace and table
- ✅ **Complete CRUD Operations**: Create, Read, Update, Delete users
- ✅ **Error Handling**: Proper error handling with descriptive messages
- ✅ **UUID Generation**: Automatic unique ID generation for users
- ✅ **Structured Code**: Clean separation of concerns with dedicated functions
- ✅ **Comprehensive Demo**: Full demonstration of all operations

## Prerequisites

- Go 1.21 or higher
- ScyllaDB running on `localhost:9042`

## ScyllaDB Setup

### Using Docker (Recommended)

```bash
# Pull and run ScyllaDB container
docker run --name scylladb -p 9042:9042 -d scylladb/scylla:latest

# Wait for ScyllaDB to start (usually takes 30-60 seconds)
docker logs -f scylladb
```

### Using Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: '3.8'
services:
  scylladb:
    image: scylladb/scylla:latest
    container_name: scylladb
    ports:
      - "9042:9042"
    environment:
      - SCYLLA_CLUSTER_NAME=test-cluster
    volumes:
      - scylla_data:/var/lib/scylla
    command: --smp 1 --memory 750M --overprovisioned 1 --api-address 0.0.0.0

volumes:
  scylla_data:
```

Then run:

```bash
docker-compose up -d
```

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd crud-scylladb
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Prerequisites
1. **Start Docker Desktop or OrbStack** (make sure Docker daemon is running)
2. **Start ScyllaDB using Docker Compose:**
   ```bash
   docker-compose up -d
   ```
   
   Wait for ScyllaDB to be ready (usually takes 30-60 seconds):
   ```bash
   docker-compose logs -f scylladb
   ```
   Look for "Scylla version" in the logs to confirm it's ready.

### Running the Application

#### Option 1: REST API Server (Default)
```bash
go run main.go
```

This starts the REST API server on `http://localhost:8080` with the following endpoints:

- `GET /api/v1/health` - Health check
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

#### Option 2: CRUD Demo
```bash
go run main.go demo
```

### API Usage Examples

#### 1. Health Check
```bash
curl http://localhost:8080/api/v1/health
```

#### 2. Create a User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

#### 3. Get All Users
```bash
curl http://localhost:8080/api/v1/users
```

#### 4. Get User by ID
```bash
curl http://localhost:8080/api/v1/users/{user-id}
```

#### 5. Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{"name": "John Smith", "email": "johnsmith@example.com"}'
```

#### 6. Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user-id}
```

### Automated Testing

Use the provided test script to automatically test all API endpoints:

```bash
./test_api.sh
```

This script will:
1. Check API health
2. Create a test user
3. Retrieve all users
4. Get user by ID
5. Update the user
6. Verify the update
7. Delete the user
8. Verify deletion

### Expected Output

#### REST API Server Startup:
```
Connected to ScyllaDB successfully!
Database initialized successfully!
🚀 Starting REST API server on http://localhost:8080
📚 API Documentation:
   GET    /api/v1/health          - Health check
   GET    /api/v1/users           - Get all users
   POST   /api/v1/users           - Create user
   GET    /api/v1/users/{id}      - Get user by ID
   PUT    /api/v1/users/{id}      - Update user
   DELETE /api/v1/users/{id}      - Delete user

💡 Run with 'go run main.go demo' to see CRUD demo
```

#### CRUD Demo Output:
```
Connected to ScyllaDB successfully!
Database initialized successfully!

=== CRUD Operations Demo ===

1. Creating user...
✓ User created successfully with ID: 123e4567-e89b-12d3-a456-426614174000

2. Reading user...
✓ Found user: {ID:123e4567-e89b-12d3-a456-426614174000 Name:John Doe Email:john@example.com CreatedAt:2024-01-15 10:30:45}

3. Updating user...
✓ User updated successfully
✓ Updated user: {ID:123e4567-e89b-12d3-a456-426614174000 Name:John Smith Email:johnsmith@example.com CreatedAt:2024-01-15 10:30:45}

4. Listing all users...
✓ Found 1 users:
   1. John Smith (johnsmith@example.com) - 2024-01-15 10:30:45

5. Deleting user...
✓ User deleted successfully
✓ Confirmed: User no longer exists

=== CRUD Operations Demo Completed ===
```

## Code Structure

### Data Model

```go
type User struct {
    ID        string    `db:"id"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
}
```

### Available Functions

- `initializeDatabase(session)` - Creates keyspace and table
- `createUser(session, user)` - Inserts a new user
- `getUserByID(session, id)` - Retrieves user by ID
- `updateUser(session, user)` - Updates existing user
- `deleteUser(session, id)` - Deletes user by ID
- `getAllUsers(session)` - Retrieves all users

## Database Schema

### Keyspace: `example`
- Replication Strategy: SimpleStrategy
- Replication Factor: 1

### Table: `users`
```sql
CREATE TABLE users (
    id text PRIMARY KEY,
    name text,
    email text,
    created_at timestamp
);
```

## Dependencies

- `github.com/gocql/gocql` - Cassandra/ScyllaDB driver
- `github.com/google/uuid` - UUID generation
- `github.com/scylladb/gocqlx/v2` - Enhanced ScyllaDB client

## Error Handling

The application includes comprehensive error handling:

- Connection failures
- Database initialization errors
- CRUD operation failures
- Not found errors
- Proper error wrapping with context

## Troubleshooting

### Common Issues

1. **Connection refused**: Ensure ScyllaDB is running on localhost:9042
2. **Timeout errors**: Increase connection timeout in cluster configuration
3. **Import errors**: Run `go mod tidy` to download dependencies

### Checking ScyllaDB Status

```bash
# Check if ScyllaDB container is running
docker ps | grep scylladb

# Check ScyllaDB logs
docker logs scylladb

# Connect to ScyllaDB shell
docker exec -it scylladb cqlsh
```

## License

MIT License