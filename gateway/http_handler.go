package main

import "net/http"

type httpHandler struct {
	// gateway

}

func NewHttpHandler() *httpHandler {
	return &httpHandler{}
}

func (h *httpHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/customers/{customerID}/orders", h.createOrder)
}

func (h *httpHandler) createOrder(w http.ResponseWriter, r *http.Request) {
}
