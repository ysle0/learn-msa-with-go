package gateway

import (
	"context"

	pb "github.com/ysle0/omsv2/common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
}
