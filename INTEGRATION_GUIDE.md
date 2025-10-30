# Maia Integration Guide

**Complete guide for integrating Maia gRPC services into Hama and Valar**

**Last Updated:** 2025-10-30

---

## Table of Contents

1. [Overview](#overview)
2. [Step 1: Build Maia](#step-1-build-maia)
3. [Step 2: Integrate into Valar (Server)](#step-2-integrate-into-valar-server)
4. [Step 3: Integrate into Hama (Client)](#step-3-integrate-into-hama-client)
5. [Step 4: Testing](#step-4-testing)
6. [Step 5: Deployment](#step-5-deployment)
7. [Troubleshooting](#troubleshooting)

---

## Overview

This guide shows you how to:
1. **Build Maia** - Generate gRPC code from proto definitions
2. **Implement Valar Server** - Create gRPC server for email/SMS
3. **Implement Hama Client** - Use gRPC client to call Valar
4. **Test Integration** - Verify end-to-end communication
5. **Deploy** - Update Kubernetes configuration

### Architecture

```
┌──────────────────────┐
│       Maia           │
│  (Proto Definitions) │
│                      │
│  - email.proto       │
│  - sms.proto         │
└──────────────────────┘
          │
          │ imports
          │
    ┌─────┴─────┐
    │           │
    ▼           ▼
┌────────┐  ┌─────────┐
│ Hama   │  │ Valar   │
│(Client)│  │(Server) │
│        │──│         │
│  gRPC  │  │  gRPC   │
│  Call  │  │ Service │
└────────┘  └─────────┘
```

---

## Step 1: Build Maia

### 1.1 Clone or Initialize Maia

```bash
cd /Users/mesongosibuti/Projects/Cloud9/api/maia

# Verify proto files exist
ls -la proto/email/v1/email.proto
ls -la proto/sms/v1/sms.proto
```

### 1.2 Install Dependencies

```bash
# Install protoc compiler (if not installed)
./scripts/install-protoc.sh

# Or manually:
# macOS: brew install protobuf
# Linux: sudo apt install protobuf-compiler

# Install Go plugins
make install-tools
```

### 1.3 Generate gRPC Code

```bash
# Generate all proto code
make proto

# Verify generated files
ls -la proto/email/v1/*.pb.go
ls -la proto/sms/v1/*.pb.go
```

**Expected output:**
```
proto/email/v1/email.pb.go         # Message definitions
proto/email/v1/email_grpc.pb.go    # gRPC client/server stubs
proto/sms/v1/sms.pb.go             # Message definitions
proto/sms/v1/sms_grpc.pb.go        # gRPC client/server stubs
```

### 1.4 Initialize Go Module (if needed)

```bash
go mod download
go mod tidy
```

### 1.5 Commit Generated Code

**Important:** Commit the generated `.pb.go` files so consumers don't need to regenerate them.

```bash
git add proto/
git commit -m "chore: generate gRPC code for email and SMS services"
git push origin main
```

---

## Step 2: Integrate into Valar (Server)

### 2.1 Add Maia Dependency

**In `valar/go.mod`:**

```go
module github.com/Cloud9Money/valar

go 1.23

require (
    github.com/Cloud9Money/maia v1.0.0  // Add this
    google.golang.org/grpc v1.59.0
    google.golang.org/protobuf v1.31.0
    // ... other dependencies
)
```

**Update dependencies:**
```bash
cd valar
go get github.com/Cloud9Money/maia@latest
go mod tidy
```

### 2.2 Create gRPC Server Implementation

**Create: `valar/internal/grpc/email_server.go`**

```go
package grpc

import (
    "context"
    "time"

    emailv1 "github.com/Cloud9Money/maia/proto/email/v1"
    "github.com/Cloud9Money/valar/internal/service"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type EmailServer struct {
    emailv1.UnimplementedEmailServiceServer
    emailService *service.EmailService
}

func NewEmailServer(emailService *service.EmailService) *EmailServer {
    return &EmailServer{
        emailService: emailService,
    }
}

func (s *EmailServer) SendVerificationEmail(ctx context.Context, req *emailv1.SendVerificationEmailRequest) (*emailv1.SendEmailResponse, error) {
    // Validate
    if req.To == "" {
        return nil, status.Error(codes.InvalidArgument, "recipient email is required")
    }

    // Call your internal email service
    messageID, err := s.emailService.SendVerificationEmail(ctx, req.To, req.VerificationToken, req.UserName)
    if err != nil {
        return &emailv1.SendEmailResponse{
            Success: false,
            Error:   err.Error(),
            Status:  "failed",
            Timestamp: time.Now().Unix(),
        }, nil
    }

    return &emailv1.SendEmailResponse{
        MessageId: messageID,
        Success:   true,
        Status:    "sent",
        Timestamp: time.Now().Unix(),
    }, nil
}

// Implement other methods...
```

**See complete example:** `maia/examples/valar-server/email_server.go`

### 2.3 Update Valar Main Server

**Update: `valar/cmd/server/main.go`**

```go
package main

import (
    "fmt"
    "log"
    "net"

    emailv1 "github.com/Cloud9Money/maia/proto/email/v1"
    grpcserver "github.com/Cloud9Money/valar/internal/grpc"
    "github.com/Cloud9Money/valar/internal/service"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

func main() {
    grpcPort := getEnv("GRPC_PORT", "50051")

    // Create gRPC listener
    lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    // Create gRPC server
    grpcServer := grpc.NewServer()

    // Initialize email service (with your dependencies)
    emailService := service.NewEmailService(/* dependencies */)

    // Register Email Service
    emailServer := grpcserver.NewEmailServer(emailService)
    emailv1.RegisterEmailServiceServer(grpcServer, emailServer)

    // Enable reflection for debugging
    reflection.Register(grpcServer)

    log.Printf("✅ gRPC server listening on :%s", grpcPort)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
```

### 2.4 Update Valar Dockerfile (if needed)

Ensure GRPC_PORT is exposed:

```dockerfile
# Expose gRPC port
EXPOSE 50051
```

---

## Step 3: Integrate into Hama (Client)

### 3.1 Add Maia Dependency

**In `hama/go.mod`:**

```go
module github.com/Cloud9Money/hama

go 1.23

require (
    github.com/Cloud9Money/maia v1.0.0  // Add this
    google.golang.org/grpc v1.59.0
    google.golang.org/protobuf v1.31.0
    // ... other dependencies
)
```

**Update dependencies:**
```bash
cd hama
go get github.com/Cloud9Money/maia@latest
go mod tidy
```

### 3.2 Create gRPC Client

**Create: `hama/internal/clients/email_client.go`**

```go
package clients

import (
    "context"
    "fmt"
    "time"

    emailv1 "github.com/Cloud9Money/maia/proto/email/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type EmailClient struct {
    client  emailv1.EmailServiceClient
    conn    *grpc.ClientConn
    timeout time.Duration
}

func NewEmailClient(valarEndpoint string) (*EmailClient, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    conn, err := grpc.DialContext(
        ctx,
        valarEndpoint,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Valar: %w", err)
    }

    return &EmailClient{
        client:  emailv1.NewEmailServiceClient(conn),
        conn:    conn,
        timeout: 10 * time.Second,
    }, nil
}

func (c *EmailClient) Close() error {
    return c.conn.Close()
}

func (c *EmailClient) SendVerificationEmail(ctx context.Context, email, token, userName string) error {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    resp, err := c.client.SendVerificationEmail(ctx, &emailv1.SendVerificationEmailRequest{
        To:                email,
        VerificationToken: token,
        UserName:          userName,
    })

    if err != nil {
        return fmt.Errorf("gRPC call failed: %w", err)
    }

    if !resp.Success {
        return fmt.Errorf("email send failed: %s", resp.Error)
    }

    return nil
}

// Implement other methods...
```

**See complete example:** `maia/examples/hama-client/email_client.go`

### 3.3 Update Hama Handlers

**Update: `hama/internal/handlers/auth_handler.go`**

```go
package handlers

import (
    "context"
    "github.com/Cloud9Money/hama/internal/clients"
)

type AuthHandler struct {
    emailClient *clients.EmailClient
    // ... other dependencies
}

func NewAuthHandler(valarEndpoint string, /* other deps */) (*AuthHandler, error) {
    emailClient, err := clients.NewEmailClient(valarEndpoint)
    if err != nil {
        return nil, fmt.Errorf("failed to create email client: %w", err)
    }

    return &AuthHandler{
        emailClient: emailClient,
    }, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *RegisterRequest) error {
    // 1. Create user
    user, err := h.userService.Create(ctx, req)
    if err != nil {
        return err
    }

    // 2. Generate token
    token, err := h.tokenService.GenerateVerificationToken(user.ID)
    if err != nil {
        return err
    }

    // 3. Send verification email via gRPC (async)
    go func() {
        err := h.emailClient.SendVerificationEmail(
            context.Background(),
            user.Email,
            token,
            user.Name,
        )
        if err != nil {
            log.Error("Failed to send verification email", "error", err)
        }
    }()

    return nil
}
```

### 3.4 Update Hama Main

**Update: `hama/cmd/server/main.go`**

```go
func main() {
    // Get Valar endpoint from environment
    valarEndpoint := os.Getenv("VALAR_GRPC_ENDPOINT")
    if valarEndpoint == "" {
        valarEndpoint = "valar-grpc.cloud9-api.svc.cluster.local:50051"
    }

    // Initialize handlers with email client
    authHandler, err := handlers.NewAuthHandler(valarEndpoint, /* other deps */)
    if err != nil {
        log.Fatalf("Failed to initialize auth handler: %v", err)
    }

    // ... rest of server setup
}
```

---

## Step 4: Testing

### 4.1 Unit Testing with Mocks

**Create mock for testing:**

```go
// hama/internal/clients/email_client_mock.go
package clients

import "context"

type MockEmailClient struct {
    SendVerificationEmailFunc func(ctx context.Context, email, token, userName string) error
}

func (m *MockEmailClient) SendVerificationEmail(ctx context.Context, email, token, userName string) error {
    if m.SendVerificationEmailFunc != nil {
        return m.SendVerificationEmailFunc(ctx, email, token, userName)
    }
    return nil
}
```

**Use in tests:**

```go
func TestRegister(t *testing.T) {
    mockEmail := &clients.MockEmailClient{
        SendVerificationEmailFunc: func(ctx context.Context, email, token, userName string) error {
            assert.Equal(t, "user@example.com", email)
            return nil
        },
    }

    handler := &AuthHandler{
        emailClient: mockEmail,
    }

    // ... test registration
}
```

### 4.2 Integration Testing

**Test with real gRPC connection (local):**

```bash
# Terminal 1: Start Valar
cd valar
go run cmd/server/main.go

# Terminal 2: Test with grpcurl
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 email.v1.EmailService/SendVerificationEmail

# Terminal 3: Run Hama with Valar endpoint
cd hama
VALAR_GRPC_ENDPOINT=localhost:50051 go run cmd/server/main.go
```

### 4.3 Testing in Kubernetes

```bash
# Deploy both services
kubectl apply -k infra/k8s/overlays/dev

# Check pods are running
kubectl get pods -n cloud9-api

# Test from Hama pod
kubectl exec -it dev-hama-xxxxx -n cloud9-api -c hama -- /bin/sh
# Inside pod:
grpcurl -plaintext valar-grpc.cloud9-api.svc.cluster.local:50051 list
```

---

## Step 5: Deployment

### 5.1 Verify Kubernetes Configuration

**Check Hama deployment has VALAR_GRPC_ENDPOINT:**

`infra/k8s/base/hama/deployment.yaml`:
```yaml
env:
  - name: VALAR_GRPC_ENDPOINT
    value: "valar-grpc.cloud9-api.svc.cluster.local:50051"
```

**Check Valar has gRPC service:**

`infra/k8s/base/valar/service.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: valar-grpc
spec:
  ports:
    - name: grpc
      port: 50051
      targetPort: 50051
  selector:
    app: valar
```

### 5.2 Build and Push Docker Images

```bash
# Build Valar with Maia dependency
cd valar
docker build -t africa-south1-docker.pkg.dev/cloud9-api-dev/cloud9-api/valar:latest .
docker push africa-south1-docker.pkg.dev/cloud9-api-dev/cloud9-api/valar:latest

# Build Hama with Maia dependency
cd ../hama
docker build -t africa-south1-docker.pkg.dev/cloud9-api-dev/cloud9-api/hama:latest .
docker push africa-south1-docker.pkg.dev/cloud9-api-dev/cloud9-api/hama:latest
```

### 5.3 Deploy to Kubernetes

```bash
# Deploy updated services
kubectl apply -k infra/k8s/overlays/dev

# Watch rollout
kubectl rollout status deployment/dev-valar -n cloud9-api
kubectl rollout status deployment/dev-hama -n cloud9-api

# Check logs
kubectl logs -f -l app=valar -n cloud9-api -c valar
kubectl logs -f -l app=hama -n cloud9-api -c hama
```

### 5.4 Verify gRPC Communication

```bash
# Shell into Hama pod
kubectl exec -it $(kubectl get pod -l app=hama -n cloud9-api -o jsonpath='{.items[0].metadata.name}') -n cloud9-api -c hama -- /bin/sh

# Test gRPC connection from inside Hama
grpcurl -plaintext valar-grpc.cloud9-api.svc.cluster.local:50051 list

# Expected output:
# email.v1.EmailService
# sms.v1.SMSService
# grpc.reflection.v1alpha.ServerReflection
```

---

## Troubleshooting

### Issue: Import errors in Valar/Hama

**Error:**
```
cannot find package "github.com/Cloud9Money/maia/proto/email/v1"
```

**Solution:**
```bash
go get github.com/Cloud9Money/maia@latest
go mod tidy
```

### Issue: gRPC connection refused

**Error:**
```
gRPC call failed: connection refused
```

**Solution:**
1. Check Valar gRPC server is running:
   ```bash
   kubectl logs -l app=valar -n cloud9-api -c valar | grep "gRPC server"
   ```

2. Check service exists:
   ```bash
   kubectl get svc valar-grpc -n cloud9-api
   ```

3. Check endpoint from Hama config:
   ```bash
   kubectl get deployment dev-hama -n cloud9-api -o yaml | grep VALAR_GRPC_ENDPOINT
   ```

### Issue: Proto generation fails

**Error:**
```
protoc: command not found
```

**Solution:**
```bash
# Install protoc
./maia/scripts/install-protoc.sh

# Or manually
brew install protobuf  # macOS
sudo apt install protobuf-compiler  # Linux
```

### Issue: Generated files out of sync

**Solution:**
```bash
cd maia
make clean
make proto
git add proto/
git commit -m "chore: regenerate proto files"
```

---

## Next Steps

After successful integration:

1. ✅ **Remove direct imports** - Remove any direct Valar imports from Hama
2. ✅ **Add SMS client** - Implement SMS client in Hama using `smsv1` protos
3. ✅ **Add monitoring** - Add gRPC metrics and tracing
4. ✅ **Add retry logic** - Implement retries for failed gRPC calls
5. ✅ **Add circuit breaker** - Prevent cascading failures
6. ✅ **Update CI/CD** - Ensure Maia is built before services

---

## Additional Resources

- **Maia README:** `maia/README.md`
- **Example Server:** `maia/examples/valar-server/`
- **Example Client:** `maia/examples/hama-client/`
- **gRPC Go Tutorial:** https://grpc.io/docs/languages/go/
- **Protocol Buffers Guide:** https://protobuf.dev/

---

**Last Updated:** 2025-10-30
**Status:** ✅ Ready for Integration

