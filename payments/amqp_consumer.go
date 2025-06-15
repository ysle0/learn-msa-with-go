package main

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	pb "github.com/ysle0/omsv2/common/api"
	message_broker "github.com/ysle0/omsv2/common/message-broker"
	"log"
)

type cosumer struct {
	service PaymentService
}

func NewConsumer(service PaymentService) *cosumer {
	return &cosumer{service}
}

func (c *cosumer) Listen(ctx context.Context, ch *amqp091.Channel) {

	q, err := ch.QueueDeclare(
		message_broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %v\n", d)
			err := d.Ack(false)
			if err != nil {
				log.Printf("Failed to ack: %v\n", err)
				continue
			}

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				log.Printf("Failed to unmarshal: %v\n", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(ctx, o)
			if err != nil {
				log.Printf("Failed to create payment: %v\n", err)
				continue
			}

			log.Printf("Payment link created: %v\n", paymentLink)
		}
	}()

	<-forever
}
