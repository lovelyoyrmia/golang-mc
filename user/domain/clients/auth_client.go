package clients

import (
	"context"
	"fmt"

	"github.com/Foedie/foedie-server-v2/user/domain/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient interface {
	VerifyEmail(email string) (*pb.VerifyEmailResponse, error)
}

type AuthClient struct {
	Client pb.AuthServiceClient
}

func InitAuthServiceClient(url string) AuthServiceClient {
	cc, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := &AuthClient{
		Client: pb.NewAuthServiceClient(cc),
	}

	return c
}

func (auth *AuthClient) VerifyEmail(email string) (*pb.VerifyEmailResponse, error) {
	req := &pb.VerifyEmailParams{
		Email: email,
	}
	return auth.Client.VerifyEmail(context.Background(), req)
}
