#!/bin/bash

# Install protoc compiler
# This script installs the Protocol Buffers compiler (protoc)

set -e

echo "Installing Protocol Buffers compiler (protoc)..."

# Detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Darwin)
        echo "Detected macOS"
        if command -v brew &> /dev/null; then
            echo "Installing via Homebrew..."
            brew install protobuf
        else
            echo "Error: Homebrew not found. Please install Homebrew first:"
            echo "  /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
            exit 1
        fi
        ;;
    Linux)
        echo "Detected Linux"

        # Check if apt is available (Debian/Ubuntu)
        if command -v apt &> /dev/null; then
            echo "Installing via apt..."
            sudo apt update
            sudo apt install -y protobuf-compiler
        # Check if yum is available (RHEL/CentOS)
        elif command -v yum &> /dev/null; then
            echo "Installing via yum..."
            sudo yum install -y protobuf-compiler
        else
            echo "Error: Could not detect package manager (apt or yum)"
            echo "Please install protoc manually from: https://github.com/protocolbuffers/protobuf/releases"
            exit 1
        fi
        ;;
    *)
        echo "Error: Unsupported OS: $OS"
        echo "Please install protoc manually from: https://github.com/protocolbuffers/protobuf/releases"
        exit 1
        ;;
esac

# Verify installation
echo ""
echo "Verifying installation..."
if command -v protoc &> /dev/null; then
    PROTOC_VERSION=$(protoc --version)
    echo "✅ protoc installed successfully: $PROTOC_VERSION"
else
    echo "❌ protoc installation failed"
    exit 1
fi

echo ""
echo "Next steps:"
echo "  1. Run 'make install-tools' to install Go plugins"
echo "  2. Run 'make proto' to generate gRPC code"
