package main

// Example main.go for Hama showing how to initialize and use email client

import (
	"context"
	"log"
	"os"

	"github.com/Cloud9Money/hama/internal/clients"
)

func main() {
	// Get Valar gRPC endpoint from environment
	valarEndpoint := getEnv("VALAR_GRPC_ENDPOINT", "valar-grpc.cloud9-api.svc.cluster.local:50051")

	log.Printf("Connecting to Valar at %s", valarEndpoint)

	// Initialize email client
	emailClient, err := clients.NewEmailClient(valarEndpoint)
	if err != nil {
		log.Fatalf("Failed to create email client: %v", err)
	}
	defer emailClient.Close()

	// Example: Send verification email
	ctx := context.Background()
	err = emailClient.SendVerificationEmail(
		ctx,
		"user@example.com",
		"verification-token-123",
		"John Doe",
	)
	if err != nil {
		log.Printf("Failed to send verification email: %v", err)
	} else {
		log.Println("✅ Verification email sent successfully")
	}

	// Example: Send welcome email
	err = emailClient.SendWelcomeEmail(
		ctx,
		"newuser@example.com",
		"Jane Doe",
		"personal",
	)
	if err != nil {
		log.Printf("Failed to send welcome email: %v", err)
	} else {
		log.Println("✅ Welcome email sent successfully")
	}

	// Example: Send transaction notification
	err = emailClient.SendTransactionNotification(
		ctx,
		"user@example.com",
		"txn-123456",
		"credit",
		1000.00,
		"KES",
		5000.00,
	)
	if err != nil {
		log.Printf("Failed to send transaction notification: %v", err)
	} else {
		log.Println("✅ Transaction notification sent successfully")
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
