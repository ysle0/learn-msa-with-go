package main

import (
	"context"
	stripeProcessor "github.com/ysle0/omsv2/payments/processor/stripe"
	"log"
	"net"
	"time"

	"github.com/stripe/stripe-go/v82"
	"github.com/ysle0/omsv2/common"
	"github.com/ysle0/omsv2/common/discovery"
	"github.com/ysle0/omsv2/common/discovery/consul"
	message_broker "github.com/ysle0/omsv2/common/message-broker"
	"google.golang.org/grpc"
)

var (
	serviceName  = "payments"
	consulAddr   = common.EnvString("CONSUL_ADDR", "localhost:8500")
	grpcAddr     = common.EnvString("GRPC_ADDR", "localhost:2001")
	amqpUser     = common.EnvString("RABBITMQ_USER", "guest")
	amqpPassword = common.EnvString("RABBITMQ_PASSWORD", "guest")
	amqpHost     = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort     = common.EnvString("RABBITMQ_PORT", "5672")
	stripeKey    = common.EnvString("STRIPE_KEY", "")
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

	// stripe setup
	stripe.Key = stripeKey
	if stripe.Key == "" {
		log.Fatal("Stripe key is not set")
	}

	// broker connection
	ch, closeChannel := message_broker.Connect(
		amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		closeChannel()
		ch.Close()
	}()

	p := stripeProcessor.NewStripeProcessor()
	svc := NewService(p)
	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(context.Background(), ch)

	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer lis.Close()

	log.Println("GRPC server listening on ", grpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
