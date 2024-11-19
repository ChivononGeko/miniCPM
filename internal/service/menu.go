package service

import (
	"fmt"
	"hot-coffee/internal/customErrors"
	"hot-coffee/internal/models"
	"log/slog"
)

type MenuRepo interface {
	GetMenusRepo() (map[string]models.MenuItem, error)
	UpdateMenusRepo(menuMap map[string]models.MenuItem) error
}

type InventDal interface {
	GetInventsRepo() (map[string]models.InventoryItem, error)
}

type MenuServImpl struct {
	menuRepo  MenuRepo
	inventDal InventDal
}

func NewMenuServImpl(mR MenuRepo, iD InventDal) *MenuServImpl {
	return &MenuServImpl{
		menuRepo:  mR,
		inventDal: iD,
	}
}

func (s *MenuServImpl) CreateMenuServ(menuNew models.MenuItem) error {
	if err := s.validateMenuInventory(menuNew.Ingredients); err != nil {
		slog.Error("Menu Service in CreateMenuServ")
		return err
	}

	menuMap, err := s.menuRepo.GetMenusRepo()
	if err != nil {
		slog.Error("Menu Service in CreateMenuServ")
		return err
	}

	if _, exists := menuMap[menuNew.ID]; exists {
		slog.Error("Inventory Service in CreateInventServ: The inventory already exists.")
		return fmt.Errorf("%w", customErrors.ErrExistConflict)
	}

	menuMap[menuNew.ID] = menuNew
	return s.menuRepo.UpdateMenusRepo(menuMap)
}

func (s *MenuServImpl) GetMenusServ() ([]models.MenuItem, error) {
	menuMap, err := s.menuRepo.GetMenusRepo()
	if err != nil {
		slog.Error("Menu Service in GetMenusServ")
		return nil, err
	}

	var menus []models.MenuItem
	for _, menu := range menuMap {
		menus = append(menus, menu)
	}

	return menus, nil
}

func (s *MenuServImpl) GetMenuIdServ(id string) (models.MenuItem, error) {
	menuMap, err := s.menuRepo.GetMenusRepo()
	if err != nil {
		slog.Error("Menu Service in GetMenuIdServ")
		return models.MenuItem{}, err
	}

	menu, exists := menuMap[id]
	if !exists {
		slog.Error("Menu Service in GetMenuIdServ: doesn't exist")
		return models.MenuItem{}, fmt.Errorf("%w", customErrors.ErrNotExistConflict)
	}

	return menu, nil
}

func (s *MenuServImpl) UpdateMenuIdServ(menuNew models.MenuItem) error {
	if err := s.validateMenuInventory(menuNew.Ingredients); err != nil {
		slog.Error("Menu Service in UpdateMenuIdServ")
		return err
	}

	menuMap, err := s.menuRepo.GetMenusRepo()
	if err != nil {
		slog.Error("Menu Service in UpdateMenuIdServ")
		return err
	}

	_, exists := menuMap[menuNew.ID]
	if !exists {
		slog.Error("Menu Service in GetMenuIdServ: doesn't exist")
		return fmt.Errorf("%w", customErrors.ErrNotExistConflict)
	}

	menuMap[menuNew.ID] = menuNew

	return s.menuRepo.UpdateMenusRepo(menuMap)
}

func (s *MenuServImpl) DeleteMenuIdServ(id string) error {
	menuMap, err := s.menuRepo.GetMenusRepo()
	if err != nil {
		slog.Error("Menu Service in UpdateMenuIdServ")
		return err
	}
	_, exists := menuMap[id]
	if !exists {
		slog.Error("Menu Service in DeleteMenuIdServ: doesn't exist")
		return fmt.Errorf("%w", customErrors.ErrNotExistConflict)
	}

	delete(menuMap, id)

	return s.menuRepo.UpdateMenusRepo(menuMap)
}

func (s *MenuServImpl) validateMenuInventory(ingredients []models.MenuItemIngredient) error {
	inventMap, err := s.inventDal.GetInventsRepo()
	if err != nil {
		slog.Error("Menu Service in validateMenuInventory")
		return err
	}

	for _, ingredient := range ingredients {
		if _, exists := inventMap[ingredient.IngredientID]; !exists {
			slog.Error("Menu Service in validateMenuInventory: doesn't exist")
			return fmt.Errorf("%w", customErrors.ErrNotExistConflict)
		}
	}

	return nil
}
