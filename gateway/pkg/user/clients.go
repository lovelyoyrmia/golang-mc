package user

import (
	"fmt"

	"github.com/Foedie/foedie-server-v2/gateway/internal/config"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/user/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client pb.UserServiceClient
}

func InitServiceClient(c *config.Config) pb.UserServiceClient {
	cc, err := grpc.Dial(c.UserSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewUserServiceClient(cc)
}
