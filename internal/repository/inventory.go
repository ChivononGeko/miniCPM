package repository

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/customErrors"
	"hot-coffee/internal/models"
	"log/slog"
)

type InventRepoImpl struct {
	filePath string
}

func NewInventRepoImpl(filepath string) *InventRepoImpl {
	return &InventRepoImpl{
		filePath: filepath,
	}
}

func (r *InventRepoImpl) GetInventsRepo() (map[string]models.InventoryItem, error) {
	data, err := readJSON(r.filePath)
	if err != nil {
		slog.Error("Inventory repository: GetInventsRepo method")
		return nil, err
	}

	var invents []models.InventoryItem

	if err := json.Unmarshal(data, &invents); err != nil {
		slog.Error("Inventory repository in GetInventsRepo method: decoding JSON")
		return nil, fmt.Errorf("%w: %s", customErrors.ErrJsonUnmarshal, err)
	}

	inventoryMap := make(map[string]models.InventoryItem)
	for _, item := range invents {
		inventoryMap[item.IngredientID] = item
	}

	return inventoryMap, nil
}

func (r *InventRepoImpl) UpdateInventsRepo(inventMap map[string]models.InventoryItem) error {
	var invents []models.InventoryItem
	for _, invent := range inventMap {
		invents = append(invents, invent)
	}

	return saveJSONToFile(r.filePath, invents)
}
