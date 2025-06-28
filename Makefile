.PHONY: help build run test clean proto client

# Default target
help:
	@echo "Available commands:"
	@echo "  make proto     - Generate protocol buffer code"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the server (REST + gRPC)"
	@echo "  make client    - Run the gRPC client example"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean generated files"
	@echo "  make deps      - Install dependencies"

# Generate protocol buffer code
proto:
	@echo "Generating protocol buffer code..."
	@chmod +x scripts/generate-proto.sh
	@./scripts/generate-proto.sh

# Build the application
build: proto
	@echo "Building application..."
	@go build -o bin/users-service main.go

# Run the server
run: proto
	@echo "Starting server (REST + gRPC)..."
	@go run main.go

# Run the gRPC client example
client: proto
	@echo "Running gRPC client example..."
	@go run examples/grpc_client.go

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@rm -f proto/*.pb.go
	@rm -f bin/users-service

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Install protoc (macOS)
install-protoc-mac:
	@echo "Installing protoc on macOS..."
	@brew install protobuf

# Install protoc (Ubuntu/Debian)
install-protoc-ubuntu:
	@echo "Installing protoc on Ubuntu/Debian..."
	@sudo apt-get update
	@sudo apt-get install -y protobuf-compiler

# Install protoc (Windows)
install-protoc-windows:
	@echo "Installing protoc on Windows..."
	@choco install protoc
