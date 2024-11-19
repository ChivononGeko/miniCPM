package repository

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/models"
	"log/slog"
)

type MenuRepoImpl struct {
	filePath string
}

func NewMenuRepoImpl(filepath string) *MenuRepoImpl {
	return &MenuRepoImpl{
		filePath: filepath,
	}
}

func (r *MenuRepoImpl) GetMenusRepo() (map[string]models.MenuItem, error) {
	data, err := readJSON(r.filePath)
	if err != nil {
		slog.Error("Menu repository: GetMenusRepo method")
		return nil, err
	}

	var menus []models.MenuItem

	if err := json.Unmarshal(data, &menus); err != nil {
		slog.Error("Menu repository in GetMenusRepo method: decoding JSON")
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	menuMap := make(map[string]models.MenuItem)
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	return menuMap, nil
}

func (r *MenuRepoImpl) UpdateMenusRepo(menuMap map[string]models.MenuItem) error {
	var menus []models.MenuItem
	for _, menu := range menuMap {
		menus = append(menus, menu)
	}

	return saveJSONToFile(r.filePath, menus)
}
