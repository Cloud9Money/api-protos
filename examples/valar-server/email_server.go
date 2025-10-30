package grpcserver

// Example gRPC server implementation for Valar (Email Service)
// This shows how to implement the EmailService gRPC server using Maia protos

import (
	"context"
	"fmt"
	"time"

	emailv1 "github.com/Cloud9Money/maia/proto/email/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EmailServer implements the EmailService gRPC server
type EmailServer struct {
	emailv1.UnimplementedEmailServiceServer
	resendClient ResendClient // Your Resend email provider client
	logger       Logger       // Your logger interface
}

// ResendClient interface (implement this with actual Resend SDK)
type ResendClient interface {
	SendEmail(to, subject, html, text string) (messageID string, err error)
	SendWithTemplate(to, templateID string, variables map[string]string) (messageID string, err error)
}

// Logger interface
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// NewEmailServer creates a new EmailServer instance
func NewEmailServer(resendClient ResendClient, logger Logger) *EmailServer {
	return &EmailServer{
		resendClient: resendClient,
		logger:       logger,
	}
}

// SendEmail implements the SendEmail RPC
func (s *EmailServer) SendEmail(ctx context.Context, req *emailv1.SendEmailRequest) (*emailv1.SendEmailResponse, error) {
	s.logger.Info("Received SendEmail request", "to", req.To, "subject", req.Subject)

	// Validate request
	if req.To == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient email is required")
	}
	if req.Subject == "" {
		return nil, status.Error(codes.InvalidArgument, "subject is required")
	}
	if req.HtmlBody == "" && req.TextBody == "" {
		return nil, status.Error(codes.InvalidArgument, "email body is required")
	}

	// Send email via Resend
	messageID, err := s.resendClient.SendEmail(req.To, req.Subject, req.HtmlBody, req.TextBody)
	if err != nil {
		s.logger.Error("Failed to send email", "error", err, "to", req.To)
		return &emailv1.SendEmailResponse{
			Success:   false,
			Error:     err.Error(),
			Status:    "failed",
			Timestamp: time.Now().Unix(),
		}, nil
	}

	s.logger.Info("Email sent successfully", "messageID", messageID, "to", req.To)

	return &emailv1.SendEmailResponse{
		MessageId: messageID,
		Success:   true,
		Status:    "sent",
		Timestamp: time.Now().Unix(),
	}, nil
}

// SendVerificationEmail implements the SendVerificationEmail RPC
func (s *EmailServer) SendVerificationEmail(ctx context.Context, req *emailv1.SendVerificationEmailRequest) (*emailv1.SendEmailResponse, error) {
	s.logger.Info("Received SendVerificationEmail request", "to", req.To)

	// Validate request
	if req.To == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient email is required")
	}
	if req.VerificationToken == "" {
		return nil, status.Error(codes.InvalidArgument, "verification token is required")
	}

	// Build verification URL
	verificationURL := req.VerificationUrl
	if verificationURL == "" {
		// Use default verification URL if not provided
		verificationURL = fmt.Sprintf("https://app.cloud9.money/verify?token=%s", req.VerificationToken)
	}

	// Prepare template variables
	variables := map[string]string{
		"user_name":        req.UserName,
		"verification_url": verificationURL,
		"verification_token": req.VerificationToken,
		"app_name":         "Cloud9",
		"support_email":    "support@cloud9.money",
	}

	// Send email using verification template
	messageID, err := s.resendClient.SendWithTemplate(req.To, "verification-email", variables)
	if err != nil {
		s.logger.Error("Failed to send verification email", "error", err, "to", req.To)
		return &emailv1.SendEmailResponse{
			Success:   false,
			Error:     err.Error(),
			Status:    "failed",
			Timestamp: time.Now().Unix(),
		}, nil
	}

	s.logger.Info("Verification email sent successfully", "messageID", messageID, "to", req.To)

	return &emailv1.SendEmailResponse{
		MessageId: messageID,
		Success:   true,
		Status:    "sent",
		Timestamp: time.Now().Unix(),
	}, nil
}

// SendPasswordResetEmail implements the SendPasswordResetEmail RPC
func (s *EmailServer) SendPasswordResetEmail(ctx context.Context, req *emailv1.SendPasswordResetEmailRequest) (*emailv1.SendEmailResponse, error) {
	s.logger.Info("Received SendPasswordResetEmail request", "to", req.To)

	// Validate request
	if req.To == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient email is required")
	}
	if req.ResetToken == "" {
		return nil, status.Error(codes.InvalidArgument, "reset token is required")
	}

	// Build reset URL
	resetURL := req.ResetUrl
	if resetURL == "" {
		resetURL = fmt.Sprintf("https://app.cloud9.money/reset-password?token=%s", req.ResetToken)
	}

	// Prepare template variables
	variables := map[string]string{
		"user_name":       req.UserName,
		"reset_url":       resetURL,
		"reset_token":     req.ResetToken,
		"expiry_minutes":  fmt.Sprintf("%d", req.ExpiryMinutes),
		"app_name":        "Cloud9",
		"support_email":   "support@cloud9.money",
	}

	// Send email using password reset template
	messageID, err := s.resendClient.SendWithTemplate(req.To, "password-reset", variables)
	if err != nil {
		s.logger.Error("Failed to send password reset email", "error", err, "to", req.To)
		return &emailv1.SendEmailResponse{
			Success:   false,
			Error:     err.Error(),
			Status:    "failed",
			Timestamp: time.Now().Unix(),
		}, nil
	}

	s.logger.Info("Password reset email sent successfully", "messageID", messageID, "to", req.To)

	return &emailv1.SendEmailResponse{
		MessageId: messageID,
		Success:   true,
		Status:    "sent",
		Timestamp: time.Now().Unix(),
	}, nil
}

// SendWelcomeEmail implements the SendWelcomeEmail RPC
func (s *EmailServer) SendWelcomeEmail(ctx context.Context, req *emailv1.SendWelcomeEmailRequest) (*emailv1.SendEmailResponse, error) {
	s.logger.Info("Received SendWelcomeEmail request", "to", req.To)

	if req.To == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient email is required")
	}

	variables := map[string]string{
		"user_name":     req.UserName,
		"account_type":  req.AccountType,
		"app_name":      "Cloud9",
		"dashboard_url": "https://app.cloud9.money/dashboard",
	}

	messageID, err := s.resendClient.SendWithTemplate(req.To, "welcome-email", variables)
	if err != nil {
		s.logger.Error("Failed to send welcome email", "error", err, "to", req.To)
		return &emailv1.SendEmailResponse{
			Success:   false,
			Error:     err.Error(),
			Status:    "failed",
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

// SendTransactionNotification implements the SendTransactionNotification RPC
func (s *EmailServer) SendTransactionNotification(ctx context.Context, req *emailv1.SendTransactionNotificationRequest) (*emailv1.SendEmailResponse, error) {
	s.logger.Info("Received SendTransactionNotification request", "to", req.To, "txnID", req.TransactionId)

	if req.To == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient email is required")
	}

	variables := map[string]string{
		"transaction_id":   req.TransactionId,
		"transaction_type": req.TransactionType,
		"amount":           fmt.Sprintf("%.2f", req.Amount),
		"currency":         req.Currency,
		"recipient_name":   req.RecipientName,
		"timestamp":        req.Timestamp,
		"balance":          fmt.Sprintf("%.2f", req.Balance),
		"app_name":         "Cloud9",
	}

	messageID, err := s.resendClient.SendWithTemplate(req.To, "transaction-notification", variables)
	if err != nil {
		s.logger.Error("Failed to send transaction notification", "error", err, "to", req.To)
		return &emailv1.SendEmailResponse{
			Success:   false,
			Error:     err.Error(),
			Status:    "failed",
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

// SendTemplateEmail implements the SendTemplateEmail RPC
func (s *EmailServer) SendTemplateEmail(ctx context.Context, req *emailv1.SendTemplateEmailRequest) (*emailv1.SendEmailResponse, error) {
	s.logger.Info("Received SendTemplateEmail request", "to", req.To, "templateID", req.TemplateId)

	if req.To == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient email is required")
	}
	if req.TemplateId == "" {
		return nil, status.Error(codes.InvalidArgument, "template ID is required")
	}

	messageID, err := s.resendClient.SendWithTemplate(req.To, req.TemplateId, req.Variables)
	if err != nil {
		s.logger.Error("Failed to send template email", "error", err, "to", req.To, "templateID", req.TemplateId)
		return &emailv1.SendEmailResponse{
			Success:   false,
			Error:     err.Error(),
			Status:    "failed",
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
