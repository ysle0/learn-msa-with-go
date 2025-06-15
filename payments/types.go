package main

import (
	"context"
	pb "github.com/ysle0/omsv2/common/api"
)

type PaymentService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}

type PaymentStore interface {
	//
}
