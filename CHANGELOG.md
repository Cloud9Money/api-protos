# Changelog

All notable changes to the API Protos repository will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2025-10-30

### Changed
- **Repository renamed** from `api-maia` to `api-protos`
- Updated Go module path from `github.com/Cloud9Money/api-maia` to `github.com/Cloud9Money/api-protos`
- Updated all proto `go_package` options to use `api-protos`
- Regenerated all proto code with new import paths
- Updated README.md, INTEGRATION_GUIDE.md with new repository name
- **Repository made public** - No longer requires authentication for access
- Reserved "maia" name for future non-proto shared functionality

### Migration Guide
If you're using v1.0.0, update your service's `go.mod`:
```bash
# Update to new repository
go get github.com/Cloud9Money/api-protos@v1.0.1
go mod tidy

# Update imports in your code from:
# import emailv1 "github.com/Cloud9Money/api-maia/proto/email/v1"
# to:
# import emailv1 "github.com/Cloud9Money/api-protos/proto/email/v1"
```

## [1.0.0] - 2025-01-30

### Added
- Initial release of API Protos gRPC proto definitions
- Email service v1 with 6 RPCs:
  - SendEmail - Send standard email with HTML and text content
  - SendTemplateEmail - Send email using predefined template
  - SendVerificationEmail - Send verification email with token
  - SendPasswordResetEmail - Send password reset email
  - SendWelcomeEmail - Send welcome email to new users
  - SendTransactionNotification - Send transaction-related emails
- SMS service v1 with 5 RPCs:
  - SendSMS - Send standard SMS message
  - SendOTP - Send one-time password via SMS
  - SendTransactionAlert - Send transaction notification via SMS
  - SendBulkSMS - Send SMS to multiple recipients
  - VerifyOTP - Verify an OTP code
- Complete proto message definitions with validation
- Generated Go code with Protocol Buffers and gRPC
- Comprehensive documentation (README.md and INTEGRATION_GUIDE.md)
- Example server implementation for Valar
- Example client implementation for Hama
- Makefile for proto generation
- Installation script for protoc compiler

### Technical Details
- Go package: `github.com/Cloud9Money/api-maia` (renamed to `api-protos` in v1.0.1)
- Proto package: `email.v1` and `sms.v1`
- Generated code uses `google.golang.org/grpc v1.59.0`
- Generated code uses `google.golang.org/protobuf v1.31.0`

[1.0.1]: https://github.com/Cloud9Money/api-protos/releases/tag/v1.0.1
[1.0.0]: https://github.com/Cloud9Money/api-protos/releases/tag/v1.0.0
