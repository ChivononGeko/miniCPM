package models

import (
	"hot-coffee/internal/customErrors"
)

type InventoryItem struct {
	IngredientID string  `json:"ingredient_id"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
}

func NewInventoryItem(id, name, unit string, quantity float64) (*InventoryItem, error) {
	if name == "" || unit == "" || quantity <= 0 {
		return nil, customErrors.ErrInvalidInput
	}
	if id == "" {
		id = fromNameToID(name)
	}

	return &InventoryItem{
		IngredientID: id,
		Name:         name,
		Quantity:     quantity,
		Unit:         unit,
	}, nil
}
