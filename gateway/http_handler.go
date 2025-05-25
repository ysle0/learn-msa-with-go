package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ysle0/omsv2/common"
	pb "github.com/ysle0/omsv2/common/api"
)

type httpHandler struct {
	grpcClient pb.OrderServiceClient
}

func NewHttpHandler(c pb.OrderServiceClient) *httpHandler {
	return &httpHandler{grpcClient: c}
}

func (h *httpHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/customers/{customerID}/orders", h.createOrder)
}

func (h *httpHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("create order with customerID", r.PathValue("customerID"))
	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("items: %+v\n", items)

	h.grpcClient.CreateOrder(context.Background(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
}
