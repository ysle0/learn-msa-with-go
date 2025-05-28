package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/ysle0/omsv2/common"
	"github.com/ysle0/omsv2/common/discovery"
	"github.com/ysle0/omsv2/common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2000")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.MakeInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("Failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()
	defer registry.Unregister(ctx, instanceID, serviceName)

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
