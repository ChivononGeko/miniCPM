package router

import (
	"hot-coffee/internal/handler"
	"net/http"
)

func OrderRouter(h *handler.OrderHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", h.CreateOrder)
	mux.HandleFunc("GET /orders", h.GetOrders)
	mux.HandleFunc("GET /orders/{id}", h.GetOrderId)
	mux.HandleFunc("PUT /orders/{id}", h.UpdateOrderId)
	mux.HandleFunc("DELETE /orders/{id}", h.DeleteOrderId)
	mux.HandleFunc("POST /orders/{id}/close", h.CloseOrderId)

	return mux
}
