package router

import (
	"hot-coffee/internal/handler"
	"net/http"
)

func InventoryRouter(h *handler.InventHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /inventory", h.CreateInvent)
	mux.HandleFunc("GET /inventory", h.GetInvents)
	mux.HandleFunc("GET /inventory/{id}", h.GetInventId)
	mux.HandleFunc("PUT /inventory/{id}", h.UpdateInventId)
	mux.HandleFunc("DELETE /inventory/{id}", h.DeleteInventId)

	return mux
}
