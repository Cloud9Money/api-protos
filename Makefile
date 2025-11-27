# Maia - Cloud9 Shared Proto Definitions
# Makefile for generating gRPC/protobuf code

.PHONY: all proto clean install-tools help

# Default target
all: proto

# Generate all protobuf and gRPC code
proto: proto-email proto-sms proto-common proto-accounts proto-transactions proto-events proto-auth proto-documents

# Generate email service proto
proto-email:
	@echo "Generating email service proto..."
	@protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		proto/email/v1/email.proto
	@echo "✅ Email proto generated"

# Generate SMS service proto
proto-sms:
	@echo "Generating SMS service proto..."
	@protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		proto/sms/v1/sms.proto
	@echo "✅ SMS proto generated"

# Generate common types proto (for adapters)
proto-common:
	@echo "Generating common types proto..."
	@protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		proto/common/*.proto
	@echo "✅ Common proto generated"

# Generate accounts proto (for adapters)
proto-accounts:
	@echo "Generating accounts proto..."
	@protoc \
		--proto_path=proto \
		--go_out=proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=proto \
		--go-grpc_opt=paths=source_relative \
		proto/accounts/*.proto
	@echo "✅ Accounts proto generated"

# Generate transactions proto (for adapters)
proto-transactions:
	@echo "Generating transactions proto..."
	@protoc \
		--proto_path=proto \
		--go_out=proto \
		--go_opt=paths=source_relative \
		proto/transactions/*.proto
	@echo "✅ Transactions proto generated"

# Generate events proto (for adapters)
proto-events:
	@echo "Generating events proto..."
	@protoc \
		--proto_path=proto \
		--go_out=proto \
		--go_opt=paths=source_relative \
		proto/events/*.proto
	@echo "✅ Events proto generated"

# Generate auth proto (for Hama service)
proto-auth:
	@echo "Generating auth proto..."
	@protoc \
		--proto_path=proto \
		--go_out=proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=proto \
		--go-grpc_opt=paths=source_relative \
		proto/auth/*.proto
	@echo "✅ Auth proto generated"

# Generate documents proto (for Mithiril service)
proto-documents:
	@echo "Generating documents proto..."
	@protoc \
		--proto_path=proto \
		--go_out=proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=proto \
		--go-grpc_opt=paths=source_relative \
		proto/documents/*.proto
	@echo "✅ Documents proto generated"

# Install required protoc plugins
install-tools:
	@echo "Installing protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Tools installed"
	@echo ""
	@echo "Make sure protoc is installed:"
	@echo "  macOS: brew install protobuf"
	@echo "  Linux: apt install -y protobuf-compiler"
	@echo ""
	@protoc --version

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@find proto -name "*.pb.go" -delete
	@echo "✅ Clean complete"

# Verify generated files
verify:
	@echo "Verifying proto files..."
	@protoc --proto_path=proto \
		--descriptor_set_out=/dev/null \
		proto/email/v1/email.proto \
		proto/sms/v1/sms.proto \
		proto/common/*.proto \
		proto/accounts/*.proto \
		proto/transactions/*.proto \
		proto/events/*.proto \
		proto/auth/*.proto \
		proto/documents/*.proto
	@echo "✅ Proto files valid"

# Update Go dependencies
update-deps:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Dependencies updated"

# Run tests (if any)
test:
	@echo "Running tests..."
	@go test -v ./...

# Format proto files (requires clang-format)
format:
	@echo "Formatting proto files..."
	@find proto -name "*.proto" -exec clang-format -i {} \;
	@echo "✅ Proto files formatted"

# Display help
help:
	@echo "Maia - Cloud9 Shared Proto Definitions"
	@echo ""
	@echo "Available targets:"
	@echo "  make all               - Generate all proto code (default)"
	@echo "  make proto             - Generate all proto code"
	@echo "  make proto-email       - Generate email service proto only"
	@echo "  make proto-sms         - Generate SMS service proto only"
	@echo "  make proto-common      - Generate common types proto only"
	@echo "  make proto-accounts    - Generate accounts proto only"
	@echo "  make proto-transactions - Generate transactions proto only"
	@echo "  make proto-events      - Generate events proto only"
	@echo "  make proto-auth        - Generate auth proto only"
	@echo "  make proto-documents   - Generate documents proto only"
	@echo "  make install-tools - Install required protoc plugins"
	@echo "  make clean         - Remove generated files"
	@echo "  make verify        - Verify proto file validity"
	@echo "  make update-deps   - Update Go dependencies"
	@echo "  make test          - Run tests"
	@echo "  make format        - Format proto files"
	@echo "  make help          - Show this help message"
	@echo ""
	@echo "Quick start:"
	@echo "  1. make install-tools"
	@echo "  2. make proto"
	@echo "  3. Commit generated files to git"
