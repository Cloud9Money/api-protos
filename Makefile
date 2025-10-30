# Maia - Cloud9 Shared Proto Definitions
# Makefile for generating gRPC/protobuf code

.PHONY: all proto clean install-tools help

# Default target
all: proto

# Generate all protobuf and gRPC code
proto: proto-email proto-sms

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
		proto/sms/v1/sms.proto
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
	@echo "  make all           - Generate all proto code (default)"
	@echo "  make proto         - Generate all proto code"
	@echo "  make proto-email   - Generate email service proto only"
	@echo "  make proto-sms     - Generate SMS service proto only"
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
