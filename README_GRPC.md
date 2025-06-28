# gRPC Implementation for User Service

This document describes the gRPC implementation added to your user service.

## Overview

The user service now supports both REST API (existing) and gRPC (new) protocols. The gRPC server runs alongside the REST server, providing the same functionality through a more efficient binary protocol.

## Features

- **Dual Protocol Support**: REST API and gRPC running simultaneously
- **Protocol Buffers**: Type-safe message definitions
- **Reflection**: Built-in gRPC reflection for debugging
- **Interceptors**: Request/response logging and error handling
- **Graceful Shutdown**: Proper server shutdown handling

## Prerequisites

1. **Install Protocol Buffers Compiler**:
   ```bash
   # On Windows (using Chocolatey)
   choco install protoc
   
   # On macOS
   brew install protobuf
   
   # On Linux
   sudo apt-get install protobuf-compiler
   ```

2. **Install Go Protocol Buffers Plugins**:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

## Project Structure

```
users-service/
├── proto/
│   └── user-service.proto          # Protocol buffer definitions
├── grpc/
│   ├── server.go                   # gRPC service implementation
│   └── manager.go                  # gRPC server management
├── examples/
│   └── grpc_client.go              # Example gRPC client
├── scripts/
│   └── generate-proto.sh           # Script to generate Go code
└── main.go                         # Updated to run both servers
```

## Configuration

The gRPC server configuration is in `resources/application.yml`:

```yaml
grpc:
  port: 9090
  enableReflection: true
  maxConcurrentStreams: 100
  maxConnectionIdle: 300s
  maxConnectionAge: 600s
  time: 120s
  timeout: 20s
```

## Usage

### 1. Generate Protocol Buffer Code

First, generate the Go code from the protobuf definitions:

```bash
# Make the script executable
chmod +x scripts/generate-proto.sh

# Run the generation script
./scripts/generate-proto.sh
```

This will create:
- `proto/user-service.pb.go` - Protocol buffer message types
- `proto/user-service_grpc.pb.go` - gRPC service definitions

### 2. Run the Server

The server now runs both REST and gRPC servers:

```bash
go run main.go
```

You should see output like:
```
Running Server on port : 8001
Starting gRPC server on port : 9090
```

### 3. Test with gRPC Client

Use the example client to test the gRPC functionality:

```bash
go run examples/grpc_client.go
```

### 4. Use gRPC Reflection (Optional)

You can use tools like `grpcurl` to interact with the server:

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:9090 list

# List methods for UserService
grpcurl -plaintext localhost:9090 list userservice.UserService

# Call CreateUser method
grpcurl -plaintext -d '{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "password123"
}' localhost:9090 userservice.UserService/CreateUser
```

## API Methods

The gRPC service provides the following methods:

### User Management
- `CreateUser` - Create a new user
- `GetUser` - Retrieve user by ID, email, or phone
- `UpdateUser` - Update user information
- `DeleteUser` - Delete a user
- `ListUsers` - List users with pagination

### Authentication
- `Login` - User login
- `Logout` - User logout

### Password Management
- `ChangePassword` - Change user password
- `ResetPassword` - Reset user password

### User Status
- `SuspendUser` - Suspend a user
- `ActivateUser` - Activate a suspended user

## Implementation Status

Currently implemented:
- ✅ Server infrastructure
- ✅ Protocol buffer definitions
- ✅ CreateUser (partially - calls existing business logic)
- ✅ Basic error handling and logging
- ✅ Graceful shutdown

To be implemented:
- ⏳ Complete business logic integration for all methods
- ⏳ Authentication and authorization
- ⏳ Input validation
- ⏳ Error mapping from business logic to gRPC status codes

## Development

### Adding New Methods

1. **Update Protocol Buffer Definition** (`proto/user-service.proto`):
   ```protobuf
   rpc NewMethod(NewMethodRequest) returns (NewMethodResponse);
   ```

2. **Regenerate Go Code**:
   ```bash
   ./scripts/generate-proto.sh
   ```

3. **Implement in Server** (`grpc/server.go`):
   ```go
   func (s *UserServiceServer) NewMethod(ctx context.Context, req *proto.NewMethodRequest) (*proto.NewMethodResponse, error) {
       // Implementation here
   }
   ```

### Error Handling

Use gRPC status codes for errors:

```go
import "google.golang.org/grpc/status"
import "google.golang.org/grpc/codes"

// Return appropriate error
return nil, status.Error(codes.InvalidArgument, "Invalid input")
return nil, status.Error(codes.NotFound, "User not found")
return nil, status.Error(codes.Internal, "Internal server error")
```

## Performance Benefits

gRPC provides several advantages over REST:

- **Binary Protocol**: More efficient than JSON
- **HTTP/2**: Multiplexing, compression, and streaming
- **Type Safety**: Compile-time type checking
- **Code Generation**: Automatic client/server code generation
- **Streaming**: Support for real-time communication

## Troubleshooting

### Common Issues

1. **"protoc: command not found"**
   - Install Protocol Buffers compiler (see Prerequisites)

2. **"protoc-gen-go: command not found"**
   - Install Go protobuf plugins (see Prerequisites)

3. **Port already in use**
   - Change the gRPC port in `application.yml`

4. **Import errors**
   - Run `go mod tidy` to update dependencies
   - Ensure protobuf files are generated correctly

### Debugging

Enable gRPC reflection for debugging:
```yaml
grpc:
  enableReflection: true
```

Use tools like:
- `grpcurl` for command-line testing
- `grpcui` for web-based testing
- `grpc-gateway` for REST-to-gRPC proxy

## Next Steps

1. Complete the implementation of all gRPC methods
2. Add comprehensive testing
3. Implement authentication middleware
4. Add metrics and monitoring
5. Consider implementing streaming methods for real-time features 