package clients

// Example gRPC client implementation for Hama (consuming Valar's Email Service)
// This shows how to use Maia protos to call Valar gRPC server

import (
	"context"
	"fmt"
	"time"

	emailv1 "github.com/Cloud9Money/maia/proto/email/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// EmailClient wraps the gRPC client for email operations
type EmailClient struct {
	client  emailv1.EmailServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

// NewEmailClient creates a new email client connected to Valar
func NewEmailClient(valarEndpoint string) (*EmailClient, error) {
	// Connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create gRPC connection
	conn, err := grpc.DialContext(
		ctx,
		valarEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Valar at %s: %w", valarEndpoint, err)
	}

	return &EmailClient{
		client:  emailv1.NewEmailServiceClient(conn),
		conn:    conn,
		timeout: 10 * time.Second,
	}, nil
}

// Close closes the gRPC connection
func (c *EmailClient) Close() error {
	return c.conn.Close()
}

// SendVerificationEmail sends an email verification link to the user
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

// SendPasswordResetEmail sends a password reset email to the user
func (c *EmailClient) SendPasswordResetEmail(ctx context.Context, email, resetToken, userName string, expiryMinutes int32) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.SendPasswordResetEmail(ctx, &emailv1.SendPasswordResetEmailRequest{
		To:            email,
		ResetToken:    resetToken,
		UserName:      userName,
		ExpiryMinutes: expiryMinutes,
	})

	if err != nil {
		return fmt.Errorf("gRPC call failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("email send failed: %s", resp.Error)
	}

	return nil
}

// SendWelcomeEmail sends a welcome email to newly registered users
func (c *EmailClient) SendWelcomeEmail(ctx context.Context, email, userName, accountType string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.SendWelcomeEmail(ctx, &emailv1.SendWelcomeEmailRequest{
		To:          email,
		UserName:    userName,
		AccountType: accountType,
	})

	if err != nil {
		return fmt.Errorf("gRPC call failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("email send failed: %s", resp.Error)
	}

	return nil
}

// SendTransactionNotification sends transaction notification email
func (c *EmailClient) SendTransactionNotification(ctx context.Context, email, txnID, txnType string, amount float64, currency string, balance float64) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.SendTransactionNotification(ctx, &emailv1.SendTransactionNotificationRequest{
		To:              email,
		TransactionId:   txnID,
		TransactionType: txnType,
		Amount:          amount,
		Currency:        currency,
		Timestamp:       time.Now().Format(time.RFC3339),
		Balance:         balance,
	})

	if err != nil {
		return fmt.Errorf("gRPC call failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("email send failed: %s", resp.Error)
	}

	return nil
}

// SendCustomEmail sends a custom email with subject and body
func (c *EmailClient) SendCustomEmail(ctx context.Context, to, subject, htmlBody, textBody string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.SendEmail(ctx, &emailv1.SendEmailRequest{
		To:       to,
		Subject:  subject,
		HtmlBody: htmlBody,
		TextBody: textBody,
	})

	if err != nil {
		return fmt.Errorf("gRPC call failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("email send failed: %s", resp.Error)
	}

	return nil
}

// Example usage in Hama's auth handler
/*
package handlers

import (
	"context"
	"github.com/Cloud9Money/hama/internal/clients"
)

type AuthHandler struct {
	emailClient *clients.EmailClient
	// ... other dependencies
}

func NewAuthHandler(valarEndpoint string) (*AuthHandler, error) {
	emailClient, err := clients.NewEmailClient(valarEndpoint)
	if err != nil {
		return nil, err
	}

	return &AuthHandler{
		emailClient: emailClient,
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *RegisterRequest) error {
	// 1. Create user in database
	user, err := h.userRepo.Create(ctx, req)
	if err != nil {
		return err
	}

	// 2. Generate verification token
	token, err := h.tokenService.GenerateVerificationToken(user.ID)
	if err != nil {
		return err
	}

	// 3. Send verification email via gRPC (non-blocking)
	go func() {
		err := h.emailClient.SendVerificationEmail(
			context.Background(),
			user.Email,
			token,
			user.Name,
		)
		if err != nil {
			// Log error but don't fail registration
			log.Error("Failed to send verification email", "error", err, "userID", user.ID)
		}
	}()

	return nil
}

func (h *AuthHandler) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	// 1. Find user
	user, err := h.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	// 2. Generate reset token
	resetToken, err := h.tokenService.GeneratePasswordResetToken(user.ID)
	if err != nil {
		return err
	}

	// 3. Send password reset email via gRPC
	err = h.emailClient.SendPasswordResetEmail(
		ctx,
		user.Email,
		resetToken,
		user.Name,
		30, // 30 minutes expiry
	)
	if err != nil {
		// In this case, we want to fail the request if email fails
		return err
	}

	return nil
}
*/
