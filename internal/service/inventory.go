package service

import (
	"fmt"
	"hot-coffee/internal/customErrors"
	"hot-coffee/internal/models"
	"log/slog"
)

type InventRepo interface {
	GetInventsRepo() (map[string]models.InventoryItem, error)
	UpdateInventsRepo(inventMap map[string]models.InventoryItem) error
}

type InventServImpl struct {
	inventRepo InventRepo
}

func NewInventServImpl(iR InventRepo) *InventServImpl {
	return &InventServImpl{inventRepo: iR}
}

func (s *InventServImpl) CreateInventServ(invent models.InventoryItem) error {
	inventoryMap, err := s.inventRepo.GetInventsRepo()
	if err != nil {
		slog.Error("Inventory Service in CreateInventServ")
		return err
	}

	if _, exists := inventoryMap[invent.IngredientID]; exists {
		slog.Error("Inventory Service in CreateInventServ: The inventory already exists.")
		return fmt.Errorf("%w", customErrors.ErrExistConflict)
	}

	inventoryMap[invent.IngredientID] = invent

	return s.inventRepo.UpdateInventsRepo(inventoryMap)
}

func (s *InventServImpl) GetInventsServ() ([]models.InventoryItem, error) {
	invents, err := s.inventRepo.GetInventsRepo()
	if err != nil {
		slog.Error("Inventory Service in GetInventsServ")
		return nil, err
	}
	var inventoryItems []models.InventoryItem
	for _, item := range invents {
		inventoryItems = append(inventoryItems, item)
	}

	return inventoryItems, nil
}

func (s *InventServImpl) GetInventIdServ(id string) (models.InventoryItem, error) {
	invents, err := s.inventRepo.GetInventsRepo()
	if err != nil {
		slog.Error("Inventory Service in GetInventIdServ")
		return models.InventoryItem{}, err
	}
	invent, exists := invents[id]
	if !exists {
		slog.Error("Inventory Service in GetInventIdServ: doesn't exist")
		return models.InventoryItem{}, fmt.Errorf("%w", customErrors.ErrNotExistConflict)
	}

	return invent, nil
}

func (s *InventServImpl) UpdateInventIdServ(inventUpd models.InventoryItem) error {
	invents, err := s.inventRepo.GetInventsRepo()
	if err != nil {
		slog.Error("Inventory Service in UpdateInventIdServ")
		return err
	}
	_, exists := invents[inventUpd.IngredientID]
	if !exists {
		slog.Error("Inventory Service in UpdateInventIdServ: doesn't exist")
		return fmt.Errorf("%w", customErrors.ErrNotExistConflict)
	}

	invents[inventUpd.IngredientID] = inventUpd

	return s.inventRepo.UpdateInventsRepo(invents)
}

func (s *InventServImpl) DeleteInventIdServ(id string) error {
	invents, err := s.inventRepo.GetInventsRepo()
	if err != nil {
		slog.Error("Inventory Service in DeleteInventIdServ")
		return err
	}
	_, exists := invents[id]
	if !exists {
		slog.Error("Inventory Service in DeleteInventIdServ: doesn't exist")
		return fmt.Errorf("%w", customErrors.ErrNotExistConflict)
	}

	delete(invents, id)

	return s.inventRepo.UpdateInventsRepo(invents)
}
