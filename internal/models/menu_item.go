package models

import (
	"hot-coffee/internal/customErrors"
)

type MenuItem struct {
	ID          string               `json:"product_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
}

type MenuItemIngredient struct {
	IngredientID string  `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}

func NewMenuItemIngredient(ingId string, quantity float64) *MenuItemIngredient {
	return &MenuItemIngredient{
		IngredientID: ingId,
		Quantity:     quantity,
	}
}

func NewMenuItem(id, name, description string, price float64, ingredients []MenuItemIngredient) (*MenuItem, error) {
	if name == "" || price <= 0 {
		return nil, customErrors.ErrInvalidInput
	}
	if description == "" {
		description = "Very tasty " + name
	}

	for _, ingredient := range ingredients {
		if ingredient.IngredientID == "" || ingredient.Quantity <= 0 {
			return nil, customErrors.ErrInvalidInput
		}
	}
	if id == "" {
		id = fromNameToID(name)
	}
	return &MenuItem{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Ingredients: ingredients,
	}, nil
}
