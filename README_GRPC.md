# Users Service - API Documentation

## 🏗️ Architecture Overview

This service provides both **HTTP REST APIs** and **gRPC APIs** with different purposes:

### **HTTP REST APIs (Port 8001)**
- **Purpose**: Frontend integration and OAuth flows
- **Use Cases**: Web applications, mobile apps, OAuth authentication
- **Key Features**: 
  - Google OAuth integration
  - Email/password authentication
  - User registration
  - Session management

### **gRPC APIs (Port 9090)**
- **Purpose**: Internal service-to-service communication
- **Use Cases**: Microservices, backend services, internal operations
- **Key Features**:
  - High-performance inter-service communication
  - Strong typing with Protocol Buffers
  - Streaming capabilities

## 📋 API Endpoints

### **HTTP REST Endpoints (OAuth & Public APIs)**

#### **Google OAuth**
```
GET  /v1/google/auth-url?state=random123    # Get Google authorization URL
POST /v1/google/login                       # Complete Google OAuth login
```

#### **Authentication**
```
POST /v1/login                              # Email/password login
POST /v1/signup                             # User registration
POST /v1/logout                             # User logout
```

#### **OTP Management**
```
POST /v1/sendOTP                            # Send OTP
GET  /v1/verifyOTP                          # Verify OTP
GET  /v1/resendOTP                          # Resend OTP
```

#### **User Management (Authenticated)**
```
GET    /v1/getMyDetails                     # Get user profile
PATCH  /v1/users/me                         # Update user profile
DELETE /v1/users/me                         # Delete user account
```

### **gRPC Endpoints (Internal Services)**

#### **User Management**
```protobuf
rpc Signup (SignupRequest) returns (SignupResponse);
rpc Login (LoginRequest) returns (LoginResponse);
rpc Logout (LogoutRequest) returns (LogoutResponse);
rpc GetUser (GetUserRequest) returns (GetUserResponse);
rpc UpdateUser (UpdateUserRequest) returns (UpdateUserResponse);
rpc DeleteUser (DeleteUserRequest) returns (DeleteUserResponse);
```

## 🔄 When to Use Which API?

### **Use HTTP REST APIs When:**
- ✅ Building frontend applications (React, Vue, Angular)
- ✅ Implementing OAuth flows (Google, Facebook, etc.)
- ✅ Mobile app authentication
- ✅ Public API access
- ✅ Webhook integrations
- ✅ Browser-based authentication flows

### **Use gRPC APIs When:**
- ✅ Service-to-service communication
- ✅ High-performance requirements
- ✅ Internal microservices
- ✅ Batch operations
- ✅ Real-time data streaming
- ✅ Backend-to-backend operations

## 🚀 Getting Started

### **HTTP REST API**
```bash
# Start the HTTP server
go run main.go
# Server runs on http://localhost:8001
```

### **gRPC API**
```bash
# The gRPC server starts automatically with the main application
# Server runs on localhost:9090
```

## 📝 Example Usage

### **Frontend Integration (HTTP)**
```javascript
// 1. Get Google auth URL
const response = await fetch('/v1/google/auth-url?state=random123');
const { authUrl } = await response.json();

// 2. Redirect user to Google
window.location.href = authUrl;

// 3. Handle callback and login
const code = new URLSearchParams(window.location.search).get('code');
const loginResponse = await fetch('/v1/google/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    code: code,
    redirectUri: 'http://localhost:3000/callback',
    clientName: 'myapp'
  })
});
```

### **Service Integration (gRPC)**
```go
// Connect to gRPC server
conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
client := proto.NewUserServiceClient(conn)

// Call gRPC method
response, err := client.Signup(ctx, &proto.SignupRequest{
    Email: "user@example.com",
    Password: "password123",
    FirstName: "John",
    LastName: "Doe",
})
```

## 🔧 Configuration

### **HTTP Server Configuration**
```yaml
server:
  host: "http://localhost:8001"
  port: 8001
```

### **gRPC Server Configuration**
```yaml
grpc:
  port: 9090
  enableReflection: true
  maxConcurrentStreams: 100
```

### **Google OAuth Configuration**
```yaml
google:
  client_id: "your-google-client-id"
  client_secret: "your-google-client-secret"
  redirect_url: "http://localhost:3000/auth/google/callback"
```

## 🛡️ Security Considerations

### **HTTP REST APIs**
- CORS configuration for frontend access
- Rate limiting for public endpoints
- Input validation and sanitization
- OAuth state parameter for CSRF protection

### **gRPC APIs**
- Internal network access only
- Service-to-service authentication
- Request/response validation
- Error handling and logging

## 📊 Monitoring and Logging

Both HTTP and gRPC endpoints use the same logging infrastructure:
- Request/response logging
- Error tracking
- Performance metrics
- Structured logging with Zap 