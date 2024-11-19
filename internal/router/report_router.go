package router

import (
	"hot-coffee/internal/handler"
	"net/http"
)

func ReportRouter(h *handler.ReportsHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /reports/total-sales", h.TotalSalesReportsHandler)
	mux.HandleFunc("GET /reports/popular-items", h.PopularItemsReportsHandler)

	return mux
}
