package message_broker

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func Connect(user, password, host, port string) (*amqp091.Channel, func() error) {
	addr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		user, password, host, port,
	)

	conn, err := amqp091.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(
		OrderCreatedEvent,
		"direct",
		true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(
		OrderPaidEvent,
		"fanout",
		true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return ch, ch.Close
}
