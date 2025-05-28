package discovery

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	randAddr := addrs[rand.Intn(len(addrs))]
	return grpc.NewClient(
		randAddr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

}
