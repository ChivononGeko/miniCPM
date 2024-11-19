package router

import (
	"hot-coffee/internal/handler"
	"net/http"
)

func MenuRouter(h *handler.MenuHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /menu", h.CreateMenu)
	mux.HandleFunc("GET /menu", h.GetMenus)
	mux.HandleFunc("GET /menu/{id}", h.GetMenuId)
	mux.HandleFunc("PUT /menu/{id}", h.UpdateMenuId)
	mux.HandleFunc("DELETE /menu/{id}", h.DeleteMenuId)

	return mux
}
