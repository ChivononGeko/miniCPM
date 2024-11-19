package router

import (
	"hot-coffee/internal/flags"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/repository"
	"hot-coffee/internal/service"
	"log/slog"
	"net/http"
	"path/filepath"
)

func SetupRoutes() (*http.ServeMux, error) {
	absDir, err := filepath.Abs(*flags.DIR)
	if err != nil {
		slog.Error("Error getting absolute path:", "error", err)
		return nil, err
	}
	inventoryJSON := filepath.Join(absDir, "inventory.json")
	menuJSON := filepath.Join(absDir, "menu_items.json")
	orderJSON := filepath.Join(absDir, "orders.json")

	inventRepo := repository.NewInventRepoImpl(inventoryJSON)
	inventServ := service.NewInventServImpl(inventRepo)
	inventHandler := handler.NewInventHandler(inventServ)

	menuRepo := repository.NewMenuRepoImpl(menuJSON)
	menuServ := service.NewMenuServImpl(menuRepo, inventRepo)
	menuHandler := handler.NewMenuHandler(menuServ)

	orderRepo := repository.NewOrderRepoImpl(orderJSON)
	orderServ := service.NewOrderServiceImpl(orderRepo, menuRepo, inventRepo)
	orderHandler := handler.NewOrderHandler(orderServ)

	serviceReports := service.NewReportsService(orderRepo, menuRepo)
	handlerReports := handler.NewReportsHandler(serviceReports)

	mux := http.NewServeMux()

	addRoutes(mux, "/inventory", InventoryRouter(inventHandler))
	addRoutes(mux, "/menu", MenuRouter(menuHandler))
	addRoutes(mux, "/orders", OrderRouter(orderHandler))
	addRoutes(mux, "/reports", ReportRouter(handlerReports))

	return mux, nil
}

func addRoutes(mux *http.ServeMux, path string, router http.Handler) {
	mux.Handle(path, router)
	mux.Handle(path+"/", router)
}
