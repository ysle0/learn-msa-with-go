package main

import (
	"context"
	"log"
	"net/http"
	"time"

	common "github.com/ysle0/omsv2/common"
	"github.com/ysle0/omsv2/common/discovery"
	"github.com/ysle0/omsv2/common/discovery/consul"
	"github.com/ysle0/omsv2/gateway/gateway"
)

var (
	httpAddr    = common.EnvString("HTTP_ADDR", ":8000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	serviceName = "gateway"
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.MakeInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
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

	gateway := gateway.NewGRPCGateway(registry)
	httpHandler := NewHttpHandler(gateway)

	mux := http.NewServeMux()
	httpHandler.registerRoutes(mux)

	log.Printf("HTTP server listening on %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal(err)
	}
}
