# Changelog

All notable changes to the API Protos repository will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.29] - 2026-06-04

### Added
- **BusinessDirector message** — new `role` field (field 17): role of this record on the business — `"director"` (LLC/LLP) or `"proprietor"` (SoleProprietorship). Defaults to `"director"` for pre-existing records. Lets consumers select a Sole Proprietorship's proprietor record explicitly rather than assuming `Directors[0]`.

## [0.0.28] - 2026-05-25

### Changed
- **LLCShareholder message** — expanded fields for shareholder records:
  - `id` (field 1) — unique identifier for the shareholder record
  - `entity_id` (field 2) — the LLC business entity UUID (was field 1)
  - `company_name` (field 3) — legal name of the parent/shareholder company (was field 2)
  - `reg_cert_doc_id` (field 4) — document ID for the registration certificate
- **EntityService.GetLLCShareholder RPC** — replaced with `ListLLCShareholders` returning a repeated list of shareholders
- **GetLLCShareholderRequest / GetLLCShareholderResponse** — replaced with `ListLLCShareholdersRequest` / `ListLLCShareholdersResponse`

## [0.0.26] - 2026-05-22

### Changed
- **BusinessDirector fields 13–16** — replaced base64 file strings with document ID strings:
  - `id_front_doc_id` (field 13) — document ID for front of ID document (was `id_front_side_file`)
  - `id_back_doc_id` (field 14) — document ID for back of ID document (was `id_back_side_file`)
  - `selfie_doc_id` (field 15) — document ID for selfie photo (was `selfie_file`)
  - `kra_pin_doc_id` (field 16) — document ID for KRA PIN certificate (was `kra_pin_file`)
  - Directors' documents are stored in Mithiril's entity document store; Citadel fetches base64 via `GetDocumentBase64` using these IDs

### Added
- **LLCShareholder message** — holds parent company details for an LLC business entity (`entity_id`, `company_name`)
- **GetLLCShareholderRequest / GetLLCShareholderResponse** — request/response wrapper messages
- **EntityService.GetLLCShareholder RPC** — returns shareholder company details for an LLC business entity

## [0.0.25] - 2026-05-22

### Added
- **BusinessDirector message** — 6 new fields required by Choice Bank's `submitCompanyMember` endpoint:
  - `gender` (field 11) — "Male" or "Female" string value
  - `kra_pin` (field 12) — KRA PIN number (e.g. A123456789B)
  - `id_front_side_file` (field 13) — base64-encoded front of ID document
  - `id_back_side_file` (field 14) — base64-encoded back of ID document
  - `selfie_file` (field 15) — base64-encoded selfie photo
  - `kra_pin_file` (field 16) — base64-encoded KRA PIN certificate
- All new fields are optional at the proto level; completeness is validated by Mithiril before KYB submission

### Technical Details
- Backward compatible — all existing `BusinessDirector` fields (1–10) unchanged
- Regenerated `proto/entities/entities.pb.go` and `proto/entities/entities_grpc.pb.go`

## [1.0.2] - 2025-10-30

### Added
- **Email Service - New RPCs:**
  - `SendBulkEmail` - Send multiple emails in batch
  - `GetEmailStatus` - Retrieve status of a sent email
  - `ListEmails` - List sent emails with filtering
- **SMS Service - New RPCs:**
  - `SendTemplateSMS` - Send SMS using predefined templates
  - `GetSMSStatus` - Retrieve status of a sent SMS
  - `ListSMS` - List sent SMS messages with filtering
- **New Message Types:**
  - `SendBulkEmailRequest`, `SendBulkEmailResponse`
  - `GetEmailStatusRequest`, `EmailStatus`
  - `ListEmailsRequest`, `ListEmailsResponse`
  - `SendTemplateSMSRequest`
  - `GetSMSStatusRequest`, `SMSStatus`
  - `ListSMSRequest`, `ListSMSResponse`

### Technical Details
- Email service now has 9 RPCs (was 6)
- SMS service now has 8 RPCs (was 5)
- Backward compatible - all existing RPCs unchanged
- Support for email/SMS tracking and status queries
- Support for bulk operations and filtering

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

[1.0.2]: https://github.com/Cloud9Money/api-protos/releases/tag/v1.0.2
[1.0.1]: https://github.com/Cloud9Money/api-protos/releases/tag/v1.0.1
[1.0.0]: https://github.com/Cloud9Money/api-protos/releases/tag/v1.0.0
