package services

import (
	"context"

	"github.com/Foedie/foedie-server-v2/email/domain/pb"
	"github.com/Foedie/foedie-server-v2/email/pkg/email"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	recoverAccount    = "recover-account"
	verificationEmail = "verify-email"
)

func (server *Server) SendEmail(ctx context.Context, req *pb.PayloadSendEmail) (*pb.SendEmailResponse, error) {

	var template string

	switch req.GetType() {
	case recoverAccount:
		template = "pkg/templates/email_verification_template.html"
	case verificationEmail:
		template = "pkg/templates/email_verification_template.html"
	}

	mail := email.NewEmail(
		[]string{req.GetEmail()},
		req.GetName(),
		req.GetUrl(),
		template,
		server.config,
	)

	err := mail.SendEmail()

	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.SendEmailResponse{
		Code:    "success",
		Message: "Successfully send an email",
	}, nil
}
