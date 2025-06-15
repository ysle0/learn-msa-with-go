package main

import (
	"context"
	pb "github.com/ysle0/omsv2/common/api"
	"github.com/ysle0/omsv2/payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{processor}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	// connect to payment processor
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "", err
	}

	// update order with the link.

	return link, nil
}
