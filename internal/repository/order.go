package repository

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/models"
	"log/slog"
)

type OrderRepoImpl struct {
	filePath string
}

func NewOrderRepoImpl(filepath string) *OrderRepoImpl {
	return &OrderRepoImpl{
		filePath: filepath,
	}
}

func (r *OrderRepoImpl) GetOrdersRepo() (map[string]models.Order, error) {
	data, err := readJSON(r.filePath)
	if err != nil {
		slog.Error("Order repository: GetOrdersRepo method")
		return nil, err
	}

	var orders []models.Order

	if err := json.Unmarshal(data, &orders); err != nil {
		slog.Error("Order repository in GetOrdersRepo method: decoding JSON")
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	ordersMap := make(map[string]models.Order)
	for _, order := range orders {
		ordersMap[order.ID] = order
	}

	return ordersMap, nil
}

func (r *OrderRepoImpl) UpdateOrdersRepo(ordersMap map[string]models.Order) error {
	var orders []models.Order
	for _, order := range ordersMap {
		orders = append(orders, order)
	}

	return saveJSONToFile(r.filePath, orders)
}
