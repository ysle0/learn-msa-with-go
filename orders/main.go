package main

import (
	"context"
	"log"
	"net"

	"github.com/ysle0/omsv2/common"
	"google.golang.org/grpc"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", ":3000")
)

func main() {
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer lis.Close()

	store := NewStore()
	svc := NewService(store)
	NewGrpcHandler(grpcServer, svc)

	svc.CreateOrder(context.Background())

	log.Println("GRPC server listening on ", grpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
