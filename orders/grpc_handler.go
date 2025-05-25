package main

import (
	"context"
	"log"

	pb "github.com/ysle0/omsv2/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
}

func NewGrpcHandler(grpcServer *grpc.Server) {
	h := &grpcHandler{}
	pb.RegisterOrderServiceServer(grpcServer, h)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received: %v\n", r)
	order := &pb.Order{
		ID: r.CustomerID,
	}

	return order, nil
}
