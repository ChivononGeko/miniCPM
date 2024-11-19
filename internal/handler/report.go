package handler

import (
	"hot-coffee/internal/models"
	"log/slog"
	"net/http"
)

type ReportsService interface {
	TotalSalesReportService() (models.TotalPrice, error)
	PopularItemsReportService() ([]models.PopularItem, error)
}

type ReportsHandler struct {
	reportsService ReportsService
}

func NewReportsHandler(rs ReportsService) *ReportsHandler {
	return &ReportsHandler{reportsService: rs}
}

func (rp *ReportsHandler) TotalSalesReportsHandler(w http.ResponseWriter, r *http.Request) {
	totalSales, err := rp.reportsService.TotalSalesReportService()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Get total sales successful")
	writeJSON(w, http.StatusOK, totalSales)
}

func (rp *ReportsHandler) PopularItemsReportsHandler(w http.ResponseWriter, r *http.Request) {
	popularItems, err := rp.reportsService.PopularItemsReportService()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Get popular items successful")
	writeJSON(w, http.StatusOK, map[string][]models.PopularItem{"the most popular item:": popularItems})
}
