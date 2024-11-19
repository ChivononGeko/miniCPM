package models

import (
	"hot-coffee/internal/customErrors"
	"time"
)

type Order struct {
	ID           string      `json:"order_id"`
	CustomerName string      `json:"customer_name"`
	Items        []OrderItem `json:"items"`
	Status       string      `json:"status"`
	CreatedAt    string      `json:"created_at"`
}

type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func NewOrderItem(productId string, quantity int) *OrderItem {
	return &OrderItem{
		ProductID: productId,
		Quantity:  quantity,
	}
}

func NewOrder(name string, createdTime time.Time, items []OrderItem) (*Order, error) {
	if name == "" {
		return nil, customErrors.ErrInvalidInput
	}

	for _, orderItem := range items {
		if orderItem.ProductID == "" || orderItem.Quantity <= 0 {
			return nil, customErrors.ErrInvalidInput
		}
	}
	return &Order{
		CustomerName: name,
		Items:        items,
		CreatedAt:    createdTime.Format("2006-01-02 15:04:05"),
	}, nil
}
