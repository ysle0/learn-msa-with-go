package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/ysle0/omsv2/common"
	pb "github.com/ysle0/omsv2/common/api"
	"log"
)

var gatewayHTTPAddr = common.EnvString("GATEWAY_HTTP_ADDR", "http://localhost:8080")

type stripeProcessor struct {
}

func NewStripeProcessor() *stripeProcessor {
	return &stripeProcessor{}
}

func (p *stripeProcessor) CreatePaymentLink(o *pb.Order) (string, error) {
	log.Printf("Creating payment link for order %v\n", o)

	gatewaySuccessURL := fmt.Sprintf("%s/success.html", gatewayHTTPAddr)
	log.Printf("Gateway success URL: %s\n", gatewaySuccessURL)

	var items []*stripe.CheckoutSessionLineItemParams
	for _, item := range o.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			//Price: stripe.String("price_1RaIJzFM42raIhkYqptQbQNR"),
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessURL),
	}

	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}
