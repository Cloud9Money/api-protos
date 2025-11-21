# Adapter Protocol Buffer Definitions

This directory contains Protocol Buffer definitions for Cloud9 banking adapters (Choice Bank SDK, Woodcore SDK). These proto files define common types, domain models, and event schemas used by the adapter modules.

## Overview

These proto definitions provide type-safe, language-agnostic interfaces for:
- **Common Types**: Money, address, contact, error types
- **Domain Models**: Accounts, transactions, ledger entries
- **Event Schemas**: Kafka event definitions for event-driven architecture

## Directory Structure

```
proto/
├── common/          # Common types used across all adapters
│   ├── money.proto       # Money, currency, exchange rate types
│   ├── address.proto     # Address, contact, identification types
│   ├── errors.proto      # Error and validation error types
│   └── common.proto      # Metadata, pagination, audit info
├── accounts/        # Account domain models
│   └── accounts.proto    # Account, account holder, tier definitions
├── transactions/    # Transaction domain models
│   └── transactions.proto # Transaction, transfer, bulk transfer types
├── events/          # Event schemas for Kafka
│   └── events.proto      # Event definitions for event-driven architecture
├── email/           # Email service (existing)
│   └── v1/email.proto
└── sms/             # SMS service (existing)
    └── v1/sms.proto
```

## Proto Files

### common/money.proto

Defines monetary types for financial operations:

```protobuf
message Money {
  string currency = 1;                    // ISO 4217 currency code (KES, USD, etc.)
  int64 amount_minor_units = 2;           // Amount in minor units (cents/pence)
  string display_amount = 3;              // Human-readable amount (e.g., "1,000.00")
}

enum Currency {
  CURRENCY_UNSPECIFIED = 0;
  CURRENCY_KES = 1;  // Kenyan Shilling
  CURRENCY_UGX = 2;  // Ugandan Shilling
  CURRENCY_TZS = 3;  // Tanzanian Shilling
  CURRENCY_USD = 4;  // US Dollar
  CURRENCY_EUR = 5;  // Euro
  CURRENCY_GBP = 6;  // British Pound
}

message ExchangeRate {
  string from_currency = 1;
  string to_currency = 2;
  double rate = 3;
  string effective_date = 4;  // ISO 8601
  string expiry_date = 5;     // ISO 8601
}
```

**Why minor units?**: Using integers for money amounts avoids floating-point precision issues in financial calculations.

### common/address.proto

Defines address and contact information:

```protobuf
message Address {
  string street = 1;
  string city = 2;
  string state_province = 3;
  string postal_code = 4;
  string country_code = 5;  // ISO 3166-1 alpha-2 (KE, UG, TZ)
}

message ContactInfo {
  string email = 1;
  string phone = 2;       // E.164 format
  string mobile = 3;      // E.164 format
  Address address = 4;
}

message Identification {
  string id_type = 1;     // national_id, passport, drivers_license
  string id_number = 2;
  string country_code = 3;
  string expiry_date = 4; // ISO 8601
}
```

### common/errors.proto

Standardized error types:

```protobuf
message Error {
  string code = 1;        // Error code (e.g., INVALID_AMOUNT)
  string message = 2;     // Human-readable error message
  map<string, string> details = 3;
  bool retryable = 4;     // Whether the operation can be retried
}

enum ErrorCode {
  ERROR_CODE_UNSPECIFIED = 0;
  ERROR_CODE_INVALID_REQUEST = 1;
  ERROR_CODE_UNAUTHORIZED = 2;
  ERROR_CODE_FORBIDDEN = 3;
  ERROR_CODE_NOT_FOUND = 4;
  ERROR_CODE_CONFLICT = 5;
  ERROR_CODE_RATE_LIMITED = 6;
  ERROR_CODE_INTERNAL_ERROR = 7;
}

message ValidationError {
  string field = 1;
  string message = 2;
  string code = 3;
}
```

### common/common.proto

Common utility types:

```protobuf
message Metadata {
  map<string, string> fields = 1;
}

message Pagination {
  int32 page = 1;
  int32 page_size = 2;
  int32 total_count = 3;
  int32 total_pages = 4;
}

message AuditInfo {
  string created_by = 1;
  string created_at = 2;  // ISO 8601
  string updated_by = 3;
  string updated_at = 4;  // ISO 8601
}

enum Status {
  STATUS_UNSPECIFIED = 0;
  STATUS_ACTIVE = 1;
  STATUS_INACTIVE = 2;
  STATUS_SUSPENDED = 3;
  STATUS_CLOSED = 4;
}
```

### accounts/accounts.proto

Account domain models:

```protobuf
enum AccountType {
  ACCOUNT_TYPE_UNSPECIFIED = 0;
  ACCOUNT_TYPE_CHECKING = 1;
  ACCOUNT_TYPE_SAVINGS = 2;
  ACCOUNT_TYPE_BUSINESS = 3;
  ACCOUNT_TYPE_WALLET = 4;
}

enum AccountStatus {
  ACCOUNT_STATUS_UNSPECIFIED = 0;
  ACCOUNT_STATUS_PENDING = 1;
  ACCOUNT_STATUS_ACTIVE = 2;
  ACCOUNT_STATUS_SUSPENDED = 3;
  ACCOUNT_STATUS_CLOSED = 4;
}

enum AccountTier {
  ACCOUNT_TIER_UNSPECIFIED = 0;
  ACCOUNT_TIER_BASIC = 1;     // Limited features, lower limits
  ACCOUNT_TIER_STANDARD = 2;   // Standard features and limits
  ACCOUNT_TIER_PREMIUM = 3;    // Enhanced features, higher limits
  ACCOUNT_TIER_BUSINESS = 4;   // Business-specific features
}

message Account {
  string id = 1;
  string customer_id = 2;
  AccountType account_type = 3;
  AccountStatus status = 4;
  AccountTier tier = 5;
  string currency = 6;
  Money balance = 7;
  Metadata metadata = 8;
  AuditInfo audit_info = 9;
}
```

### transactions/transactions.proto

Transaction domain models:

```protobuf
enum TransactionType {
  TRANSACTION_TYPE_UNSPECIFIED = 0;
  TRANSACTION_TYPE_DEPOSIT = 1;
  TRANSACTION_TYPE_WITHDRAWAL = 2;
  TRANSACTION_TYPE_TRANSFER = 3;
  TRANSACTION_TYPE_PAYMENT = 4;
  TRANSACTION_TYPE_REFUND = 5;
  TRANSACTION_TYPE_FEE = 6;
  TRANSACTION_TYPE_REVERSAL = 7;
}

enum TransactionStatus {
  TRANSACTION_STATUS_UNSPECIFIED = 0;
  TRANSACTION_STATUS_PENDING = 1;
  TRANSACTION_STATUS_PROCESSING = 2;
  TRANSACTION_STATUS_COMPLETED = 3;
  TRANSACTION_STATUS_FAILED = 4;
  TRANSACTION_STATUS_REVERSED = 5;
}

message Transaction {
  string id = 1;
  TransactionType type = 2;
  TransactionStatus status = 3;
  string account_id = 4;
  Money amount = 5;
  string description = 6;
  string reference = 7;
  Metadata metadata = 8;
  AuditInfo audit_info = 9;
}
```

### events/events.proto

Event schemas for Kafka:

```protobuf
message Event {
  string id = 1;               // Event ID
  string type = 2;             // Event type
  string source = 3;           // Source service
  string timestamp = 4;        // ISO 8601
  string jurisdiction = 5;     // KE, UG, TZ, etc.

  oneof payload {
    AccountCreatedEvent account_created = 10;
    TransactionCompletedEvent transaction_completed = 11;
    ComplianceAlertEvent compliance_alert = 12;
    FXRateUpdatedEvent fx_rate_updated = 13;
  }
}

message AccountCreatedEvent {
  string account_id = 1;
  string customer_id = 2;
  AccountType account_type = 3;
  string currency = 4;
}

message TransactionCompletedEvent {
  string transaction_id = 1;
  string account_id = 2;
  TransactionType type = 3;
  Money amount = 4;
}
```

## Usage in Adapters

### Choice Bank SDK

```go
import (
    "github.com/Cloud9Money/api-protos/proto/common"
    "github.com/Cloud9Money/api-protos/proto/accounts"
    "github.com/Cloud9Money/api-protos/proto/transactions"
)

// Use proto types
money := &common.Money{
    Currency:         "KES",
    AmountMinorUnits: 100000,  // 1000.00 KES
    DisplayAmount:    "1,000.00",
}

account := &accounts.Account{
    CustomerId:  "cust_123",
    AccountType: accounts.AccountType_ACCOUNT_TYPE_CHECKING,
    Status:      accounts.AccountStatus_ACCOUNT_STATUS_ACTIVE,
    Currency:    "KES",
    Balance:     money,
}
```

### Woodcore SDK

```go
import (
    "github.com/Cloud9Money/api-protos/proto/transactions"
    "github.com/Cloud9Money/api-protos/proto/common"
)

txn := &transactions.Transaction{
    Id:     "txn_456",
    Type:   transactions.TransactionType_TRANSACTION_TYPE_DEPOSIT,
    Status: transactions.TransactionStatus_TRANSACTION_STATUS_COMPLETED,
    Amount: money,
}
```

## Generating Code

From the root `protos/` directory:

```bash
# Generate all proto files (including adapters)
make proto

# Or generate adapter protos only
make proto-common
make proto-accounts
make proto-transactions
make proto-events
```

This generates `.pb.go` files alongside each `.proto` file.

## Best Practices

### 1. Use Minor Units for Money

Always represent monetary amounts as integers in minor units (cents, pence, etc.):

```protobuf
// Good ✅
int64 amount_minor_units = 2;  // 100000 = 1000.00 KES

// Bad ❌
double amount = 2;  // Floating point precision issues
```

### 2. Use Enums for Fixed Values

Define enums for values that have a fixed set of options:

```protobuf
enum AccountStatus {
  ACCOUNT_STATUS_UNSPECIFIED = 0;  // Always have UNSPECIFIED = 0
  ACCOUNT_STATUS_ACTIVE = 1;
  ACCOUNT_STATUS_SUSPENDED = 2;
  ACCOUNT_STATUS_CLOSED = 3;
}
```

### 3. ISO 8601 for Timestamps

Always use ISO 8601 format for timestamps:

```protobuf
string created_at = 8;  // "2024-11-20T10:30:00Z"
```

### 4. Use oneof for Polymorphic Data

Use `oneof` for event payloads or polymorphic types:

```protobuf
message Event {
  oneof payload {
    AccountCreatedEvent account_created = 10;
    TransactionCompletedEvent transaction_completed = 11;
  }
}
```

### 5. Namespace Enums

Prefix enum values with the enum name to avoid conflicts:

```protobuf
enum AccountType {
  ACCOUNT_TYPE_UNSPECIFIED = 0;  // Not just UNSPECIFIED
  ACCOUNT_TYPE_CHECKING = 1;     // Not just CHECKING
}
```

## Versioning

These proto definitions follow semantic versioning:

- **Breaking changes**: Require version bump (v1 → v2)
- **Backward compatible additions**: Patch version bump
- **Field additions**: Always add new fields, never remove or renumber existing fields

## Integration with Services

Services consuming these protos should:

1. Import via the `api-protos` module
2. Use the generated Go types
3. Handle all enum values (including UNSPECIFIED)
4. Validate proto messages before use

## Contributing

When adding or modifying proto definitions:

1. Follow the naming conventions above
2. Add comprehensive comments
3. Run `make verify` to validate syntax
4. Regenerate code with `make proto`
5. Update this documentation
6. Test with consuming services

## Related Documentation

- [Main README](../README.md) - API Protos overview and gRPC services
- [Choice Bank SDK](../../adapters/choice-bank-sdk/README.md)
- [Woodcore SDK](../../adapters/woodcore-sdk/README.md)
- [Adapter Architecture](../../adapter_architecture.md)

---

**Version**: 1.0.0
**Last Updated**: 2024-11-20
