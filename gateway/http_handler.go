package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ysle0/omsv2/common"
	pb "github.com/ysle0/omsv2/common/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	err := validateItems(items)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.grpcClient.CreateOrder(context.Background(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() == codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, rStatus.Message())
		return
	}
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, o)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}

	for _, i := range items {
		if i.ID == "" {
			return common.ErrItemIdRequired
		}

		if i.Quantity <= 0 {
			return common.ErrItemQuantityRequired
		}
	}

	return nil
}
