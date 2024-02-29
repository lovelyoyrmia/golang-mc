package clients

import (
	"context"
	"fmt"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EmailServiceClient struct {
	Client pb.EmailServiceClient
}

func InitEmailServiceClient(url string) EmailServiceClient {
	cc, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := EmailServiceClient{
		Client: pb.NewEmailServiceClient(cc),
	}

	return c
}

func (mail *EmailServiceClient) SendEmail(
	email string,
	name string,
	url string,
	template string,
) (*pb.SendEmailResponse, error) {
	return mail.Client.SendEmail(context.Background(), &pb.PayloadSendEmail{
		Email: email,
		Name:  name,
		Url:   url,
		Type:  template,
	})
}
