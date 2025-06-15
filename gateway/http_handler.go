package main

import (
	"log"
	"net/http"

	"github.com/ysle0/omsv2/common"
	pb "github.com/ysle0/omsv2/common/api"
	"github.com/ysle0/omsv2/gateway/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type httpHandler struct {
	gateway gateway.OrdersGateway
}

func NewHttpHandler(gateway gateway.OrdersGateway) *httpHandler {
	return &httpHandler{gateway}
}

func (h *httpHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/customers/{customerID}/orders", h.createOrder)
}

func (h *httpHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("create order with customerID", r.PathValue("customerID"))
	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteHeaderErr(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("items: %v\n", items)

	err := validateItems(items)
	if err != nil {
		common.WriteHeaderErr(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	rStatus := status.Convert(err)
	if rStatus != nil {
		log.Println(rStatus.Message())
		if rStatus.Code() == codes.InvalidArgument {
			common.WriteHeaderErr(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteHeaderErr(w, http.StatusInternalServerError, rStatus.Message())
		return
	}

	if err != nil {
		common.WriteHeaderErr(w, http.StatusBadRequest, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, order)
}
