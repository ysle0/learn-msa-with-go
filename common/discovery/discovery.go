package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(ctx context.Context, instanceID, serviceName, hostPort string) error
	Unregister(ctx context.Context, instanceID, serviceName string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceID, serviceName string) error
}

func MakeInstanceID(serviceName string) string {
	randNum := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	return fmt.Sprintf(
		"%s-%d",
		serviceName,
		randNum,
	)
}
