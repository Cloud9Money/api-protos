# API Protos - Cloud9 Shared Proto Definitions

**API Protos** is the shared gRPC/Protocol Buffers definitions package for the Cloud9 multi-jurisdiction banking platform. This repository provides common communication interfaces for all Cloud9 microservices.

**Version:** 1.0.0
**Last Updated:** 2025-10-30

---

## Table of Contents

1. [Overview](#overview)
2. [Services](#services)
3. [Installation](#installation)
4. [Building Maia](#building-maia)
5. [Using Maia in Services](#using-maia-in-services)
6. [Development](#development)
7. [Versioning](#versioning)
8. [Contributing](#contributing)

---

## Overview

### What is API Protos?

API Protos provides:
- **gRPC Service Definitions** - Email and SMS service interfaces
- **Protocol Buffers** - Strongly-typed message definitions
- **Generated Go Code** - Auto-generated client and server stubs
- **Service Decoupling** - Enables microservices to communicate without direct dependencies

### Architecture

```
┌─────────┐      ┌──────┐      ┌────────┐
│  Hama   │─────▶│ Maia │◀─────│ Valar  │
│ (Auth)  │      │      │      │(Email) │
└─────────┘      └──────┘      └────────┘
    │                              │
    │  gRPC Call                   │
    └──────────────────────────────┘
         (no direct dependency)
```

**Benefits:**
- ✅ No circular dependencies
- ✅ Independent deployment
- ✅ Type-safe communication
- ✅ Version-controlled contracts
- ✅ Easy testing with mocks

---

## Services

### Email Service (`email.v1`)

Provides email sending capabilities via Valar service.

**Available RPCs:**
- `SendEmail` - Send standard HTML/text email
- `SendTemplateEmail` - Send using predefined template
- `SendVerificationEmail` - Send email verification link
- `SendPasswordResetEmail` - Send password reset link
- `SendWelcomeEmail` - Send welcome email to new users
- `SendTransactionNotification` - Send transaction alerts

**Example:**
```go
import emailv1 "github.com/Cloud9Money/api-protos/proto/email/v1"

client := emailv1.NewEmailServiceClient(conn)
resp, err := client.SendVerificationEmail(ctx, &emailv1.SendVerificationEmailRequest{
    To:               "user@example.com",
    VerificationToken: "abc123",
    UserName:         "John Doe",
})
```

### SMS Service (`sms.v1`)

Provides SMS sending capabilities via Valar service.

**Available RPCs:**
- `SendSMS` - Send standard SMS
- `SendOTP` - Send one-time password
- `SendTransactionAlert` - Send transaction notification
- `SendBulkSMS` - Send to multiple recipients
- `VerifyOTP` - Verify OTP code

**Example:**
```go
import smsv1 "github.com/Cloud9Money/api-protos/proto/sms/v1"

client := smsv1.NewSMSServiceClient(conn)
resp, err := client.SendOTP(ctx, &smsv1.SendOTPRequest{
    To:            "+254712345678",
    OtpCode:       "123456",
    ExpiryMinutes: 5,
    Purpose:       "login",
    Jurisdiction:  "KE",
})
```

---

## Installation

### Prerequisites

**Required Tools:**
- Go 1.23 or later
- Protocol Buffers compiler (`protoc`)
- protoc-gen-go plugin
- protoc-gen-go-grpc plugin

### Install protoc

**macOS:**
```bash
brew install protobuf
```

**Linux (Debian/Ubuntu):**
```bash
sudo apt update
sudo apt install -y protobuf-compiler
```

**Linux (RHEL/CentOS):**
```bash
sudo yum install -y protobuf-compiler
```

**Or use the install script:**
```bash
cd scripts
./install-protoc.sh
```

### Install Go Plugins

```bash
make install-tools
```

This installs:
- `protoc-gen-go` - Generates Go structs from proto
- `protoc-gen-go-grpc` - Generates gRPC service stubs

### Verify Installation

```bash
protoc --version  # Should show libprotoc 3.x or later
protoc-gen-go --version
protoc-gen-go-grpc --version
```

---

## Building Proto Package

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/Cloud9Money/api-protos.git
cd maia

# Install Go dependencies
go mod download

# Install protoc plugins
make install-tools

# Generate gRPC code from proto files
make proto
```

### Generate Proto Code

```bash
# Generate all services
make proto

# Generate specific service
make proto-email
make proto-sms

# Verify proto files are valid
make verify

# Clean generated files
make clean
```

### What Gets Generated

After running `make proto`, you'll see:

```
protos/
├── proto/
│   ├── email/
│   │   └── v1/
│   │       ├── email.proto          # Proto definition
│   │       ├── email.pb.go          # Generated messages
│   │       └── email_grpc.pb.go     # Generated gRPC stubs
│   └── sms/
│       └── v1/
│           ├── sms.proto            # Proto definition
│           ├── sms.pb.go            # Generated messages
│           └── sms_grpc.pb.go       # Generated gRPC stubs
```

**Important:** Commit the generated `.pb.go` files to git so consumers don't need to regenerate them.

---

## Using Protos in Services

### Add API Protos as Dependency

**In your service's `go.mod`:**

```go
module github.com/Cloud9Money/hama

go 1.23

require (
    github.com/Cloud9Money/api-protos v1.0.0
    google.golang.org/grpc v1.59.0
    google.golang.org/protobuf v1.31.0
)
```

**Update dependencies:**
```bash
go mod download
go mod tidy
```

### Server Implementation (Valar)

**Implement the gRPC server:**

```go
package grpc

import (
    "context"
    emailv1 "github.com/Cloud9Money/api-protos/proto/email/v1"
)

type EmailServer struct {
    emailv1.UnimplementedEmailServiceServer
    emailService *service.EmailService
}

func NewEmailServer(emailService *service.EmailService) *EmailServer {
    return &EmailServer{emailService: emailService}
}

func (s *EmailServer) SendVerificationEmail(
    ctx context.Context,
    req *emailv1.SendVerificationEmailRequest,
) (*emailv1.SendEmailResponse, error) {
    // Call your internal email service
    messageID, err := s.emailService.SendVerificationEmail(ctx, req)
    if err != nil {
        return &emailv1.SendEmailResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }

    return &emailv1.SendEmailResponse{
        MessageId: messageID,
        Success:   true,
        Status:    "sent",
    }, nil
}
```

**Start gRPC server:**

```go
package main

import (
    "net"
    emailv1 "github.com/Cloud9Money/api-protos/proto/email/v1"
    "google.golang.org/grpc"
)

func main() {
    lis, _ := net.Listen("tcp", ":50051")
    grpcServer := grpc.NewServer()

    emailService := service.NewEmailService(/* deps */)
    emailServer := grpc.NewEmailServer(emailService)

    emailv1.RegisterEmailServiceServer(grpcServer, emailServer)

    grpcServer.Serve(lis)
}
```

### Client Implementation (Hama)

**Create gRPC client:**

```go
package clients

import (
    "context"
    "fmt"
    "time"

    emailv1 "github.com/Cloud9Money/api-protos/proto/email/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type EmailClient struct {
    client emailv1.EmailServiceClient
    conn   *grpc.ClientConn
}

func NewEmailClient(valarEndpoint string) (*EmailClient, error) {
    conn, err := grpc.Dial(
        valarEndpoint,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithTimeout(5*time.Second),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Valar: %w", err)
    }

    return &EmailClient{
        client: emailv1.NewEmailServiceClient(conn),
        conn:   conn,
    }, nil
}

func (c *EmailClient) SendVerificationEmail(
    ctx context.Context,
    email, token, userName string,
) error {
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

func (c *EmailClient) Close() error {
    return c.conn.Close()
}
```

**Use the client:**

```go
package handlers

func (h *AuthHandler) Register(ctx context.Context, req *RegisterRequest) error {
    // Create user, generate token...

    // Send verification email
    err := h.emailClient.SendVerificationEmail(ctx, req.Email, token, req.Name)
    if err != nil {
        log.Error("Failed to send verification email", "error", err)
        // Don't fail registration if email fails
    }

    return nil
}
```

---

## Development

### Project Structure

```
protos/
├── proto/                  # Protocol Buffer definitions
│   ├── email/
│   │   └── v1/
│   │       ├── email.proto
│   │       ├── email.pb.go (generated)
│   │       └── email_grpc.pb.go (generated)
│   └── sms/
│       └── v1/
│           ├── sms.proto
│           ├── sms.pb.go (generated)
│           └── sms_grpc.pb.go (generated)
├── scripts/                # Helper scripts
│   └── install-protoc.sh
├── Makefile                # Build automation
├── go.mod                  # Go module definition
├── go.sum                  # Go dependency checksums
├── README.md               # This file
└── .gitignore             # Git ignore rules
```

### Making Changes to Proto Files

1. **Edit proto file** (`proto/email/v1/email.proto`)
2. **Regenerate code:** `make proto`
3. **Test changes:** Update and test consuming services
4. **Commit everything:** Proto file + generated code
5. **Tag release:** `git tag v1.1.0 && git push --tags`
6. **Update consumers:** Services update to new version

### Adding a New Service

1. **Create proto directory:**
   ```bash
   mkdir -p proto/newservice/v1
   ```

2. **Create proto file:**
   ```protobuf
   syntax = "proto3";
   package newservice.v1;
   option go_package = "github.com/Cloud9Money/api-protos/proto/newservice/v1;newservicev1";

   service NewService {
       rpc DoSomething(Request) returns (Response);
   }
   ```

3. **Add to Makefile:**
   ```makefile
   proto-newservice:
       @protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/newservice/v1/newservice.proto
   ```

4. **Generate:** `make proto-newservice`

### Testing

**Validate proto files:**
```bash
make verify
```

**Test in consuming service:**
```bash
cd ../hama
go test ./internal/clients/...
```

---

## Versioning

### Semantic Versioning

API Protos follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (v2.0.0): Breaking changes to proto definitions
- **MINOR** (v1.1.0): New features (new RPCs, optional fields)
- **PATCH** (v1.0.1): Bug fixes, documentation

### Proto Versioning

Proto files use version directories (`v1`, `v2`) for major API versions:

```
proto/email/
├── v1/          # Current stable version
│   └── email.proto
└── v2/          # Future version (breaking changes)
    └── email.proto
```

**When to create v2:**
- Removing fields
- Changing field types
- Removing RPCs
- Major restructuring

**Backward-compatible changes (stay in v1):**
- Adding new fields (use optional)
- Adding new RPCs
- Adding new services

### Updating Consumers

**Check current version:**
```bash
cd hama
go list -m github.com/Cloud9Money/api-protos
```

**Update to latest:**
```bash
go get github.com/Cloud9Money/api-protos@latest
go mod tidy
```

**Update to specific version:**
```bash
go get github.com/Cloud9Money/api-protos@v1.2.0
go mod tidy
```

---

## Contributing

### Development Workflow

1. **Create branch:** `git checkout -b feature/add-notification-service`
2. **Make changes:** Edit proto files
3. **Generate code:** `make proto`
4. **Test:** Verify changes in consuming services
5. **Commit:** Include both proto and generated files
6. **Create PR:** Submit for review
7. **Tag release:** After merge, tag new version

### Commit Guidelines

```
feat(email): add attachment support to SendEmail RPC
fix(sms): correct OTP expiry field type
docs: update README with new service examples
chore: update dependencies
```

### Code Review Checklist

- [ ] Proto files follow style guide
- [ ] Generated code committed
- [ ] Backward compatibility maintained (or version bumped)
- [ ] Documentation updated
- [ ] Changes tested in at least one consumer service
- [ ] go.mod version updated

---

## Troubleshooting

### `protoc: command not found`

**Solution:**
```bash
# macOS
brew install protobuf

# Linux
sudo apt install protobuf-compiler

# Or use install script
./scripts/install-protoc.sh
```

### `protoc-gen-go: program not found`

**Solution:**
```bash
make install-tools
# Ensure $GOPATH/bin is in your $PATH
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Import errors in consuming service

**Solution:**
```bash
cd <service>
go get github.com/Cloud9Money/api-protos@latest
go mod tidy
```

### Generated files not updated

**Solution:**
```bash
make clean
make proto
```

---

## Additional Resources

- **Protocol Buffers Guide:** https://protobuf.dev/
- **gRPC Go Tutorial:** https://grpc.io/docs/languages/go/quickstart/
- **Cloud9 Architecture:** `../CLAUDE.md`
- **Valar Service:** `../valar/README.md`
- **Hama Service:** `../hama/README.md`

---

## Support

### Getting Help

1. Check this README
2. Review proto file comments
3. Check example code in consuming services
4. Open GitHub issue with:
   - Proto file excerpt
   - Error message
   - Steps to reproduce

### Maintainers

- Cloud9 Platform Team

---

**Last Updated:** 2025-10-30
**Version:** 1.0.0
**License:** Proprietary - Cloud9 Money

