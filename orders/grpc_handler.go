package main

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	pb "github.com/ysle0/omsv2/common/api"
	message_broker "github.com/ysle0/omsv2/common/message-broker"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	channel *amqp091.Channel
	service OrdersService
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrdersService, channel *amqp091.Channel) {
	h := &grpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrderServiceServer(grpcServer, h)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.Order, error) {
	//log.Printf("New order received: %v\n", r)

	order, err := h.service.CreateOrder(ctx, r)
	if err != nil {
		return nil, err
	}

	marshalledOrder, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	q, err := h.channel.QueueDeclare(
		message_broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	err = h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp091.Persistent,
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}
