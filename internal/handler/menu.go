package handler

import (
	"encoding/json"
	"errors"
	"hotC/internal/customErrors"
	"hotC/internal/models"
	"log/slog"
	"net/http"
)

type MenuServ interface {
	CreateMenuServ(menuNew models.MenuItem) error
	GetMenusServ() ([]models.MenuItem, error)
	GetMenuIdServ(id string) (models.MenuItem, error)
	UpdateMenuIdServ(menuNew models.MenuItem) error
	DeleteMenuIdServ(id string) error
}

type MenuHandler struct {
	menuServ MenuServ
}

func NewMenuHandler(mS MenuServ) *MenuHandler {
	return &MenuHandler{menuServ: mS}
}

func (h *MenuHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		return
	}

	var inputMenu models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&inputMenu); err != nil {
		slog.Error("Handler Error in CreateMenu: decoding JSON data ", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	menu, err := models.NewMenuItem(inputMenu.ID, inputMenu.Name, inputMenu.Description, inputMenu.Price, inputMenu.Ingredients)
	if err != nil {
		slog.Error("Handler Error in CreateMenu: Invalid input data", "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.menuServ.CreateMenuServ(*menu); err != nil {
		var status int
		if errors.Is(err, customErrors.ErrExistConflict) {
			status = http.StatusConflict
		} else if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in CreateMenu: creating menu", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	slog.Info("Menu created successfully", "menuID", menu.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *MenuHandler) GetMenus(w http.ResponseWriter, r *http.Request) {
	menus, err := h.menuServ.GetMenusServ()
	if err != nil {
		slog.Error("Handler Error in GetMenus: retrieving all menu", "error", err)
		writeError(w, "Failed to retrieve all menu", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(menus); err != nil {
		slog.Error("Handler Error in GetMenus: encoding JSON data", "error", err)
		writeError(w, "Failed to encode all menu to JSON", http.StatusInternalServerError)
		return
	}

	slog.Info("All menu retrieved successfully")
}

func (h *MenuHandler) GetMenuId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	menu, err := h.menuServ.GetMenuIdServ(id)
	if err != nil {
		var status int
		if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in GetMenuId: retrieving menu by ID ", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(menu); err != nil {
		slog.Error("Handler Error in GetMenuId: encoding JSON data", "error", err)
		writeError(w, "Failed to encode menu to JSON", http.StatusInternalServerError)
		return
	}

	slog.Info("Menu retrieved successfully", "menuID", id)
}

func (h *MenuHandler) UpdateMenuId(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		return
	}
	id := r.PathValue("id")

	var inputMenu models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&inputMenu); err != nil {
		slog.Error("Handler Error in UpdateMenuId: decoding JSON data", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	menu, err := models.NewMenuItem(inputMenu.ID, inputMenu.Name, inputMenu.Description, inputMenu.Price, inputMenu.Ingredients)
	if err != nil {
		slog.Error("Handler Error in UpdateMenuId: Invalid input data", "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	menu.ID = id

	if err := h.menuServ.UpdateMenuIdServ(*menu); err != nil {
		var status int
		if errors.Is(err, customErrors.ErrExistConflict) {
			status = http.StatusConflict
		} else if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in UpdateMenuId: updating menu", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	slog.Info("Menu updated successfully", "menuID", menu.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *MenuHandler) DeleteMenuId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.menuServ.DeleteMenuIdServ(id); err != nil {
		var status int
		if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in DeleteMenuId: deleting inventory by ID", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Menu deleted successfully")
}
