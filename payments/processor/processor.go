package processor

import pb "github.com/ysle0/omsv2/common/api"

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
