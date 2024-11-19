package handler

import (
	"encoding/json"
	"errors"
	"hot-coffee/internal/customErrors"
	"hot-coffee/internal/models"
	"log/slog"
	"net/http"
)

type InventServ interface {
	CreateInventServ(invent models.InventoryItem) error
	GetInventsServ() ([]models.InventoryItem, error)
	GetInventIdServ(id string) (models.InventoryItem, error)
	UpdateInventIdServ(inventUpd models.InventoryItem) error
	DeleteInventIdServ(id string) error
}

type InventHandler struct {
	inventServ InventServ
}

func NewInventHandler(iS InventServ) *InventHandler {
	return &InventHandler{inventServ: iS}
}

func (h *InventHandler) CreateInvent(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		return
	}

	var inputInvent models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&inputInvent); err != nil {
		slog.Error("Handler Error in CreateInvent: decoding JSON data ", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	invent, err := models.NewInventoryItem(inputInvent.IngredientID, inputInvent.Name, inputInvent.Unit, inputInvent.Quantity)
	if err != nil {
		slog.Error("Handler Error in CreateInvent: Invalid input data", "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.inventServ.CreateInventServ(*invent); err != nil {
		var status int
		if errors.Is(err, customErrors.ErrExistConflict) {
			status = http.StatusConflict
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in CreateInvent: creating inventory", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	slog.Info("Inventory created successfully", "inventID", invent.IngredientID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *InventHandler) GetInvents(w http.ResponseWriter, r *http.Request) {
	allInvents, err := h.inventServ.GetInventsServ()
	if err != nil {
		slog.Error("Handler Error in GetInvents: retrieving all invents", "error", err)
		writeError(w, "Failed to retrieve invents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(allInvents); err != nil {
		slog.Error("Handler Error in GetInvents: encoding JSON data", "error", err)
		writeError(w, "Failed to encode invents to JSON", http.StatusInternalServerError)
		return
	}

	slog.Info("Invents retrieved successfully")
}

func (h *InventHandler) GetInventId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	inventId, err := h.inventServ.GetInventIdServ(id)
	if err != nil {
		var status int
		if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in GetInventId: retrieving invent by ID ", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(inventId); err != nil {
		slog.Error("Handler Error in GetInventId: encoding JSON data", "error", err)
		writeError(w, "Failed to encode invent to JSON", http.StatusInternalServerError)
		return
	}

	slog.Info("Inventory retrieved successfully", "inventID", id)
}

func (h *InventHandler) UpdateInventId(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		return
	}
	id := r.PathValue("id")

	var inputInvent models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&inputInvent); err != nil {
		slog.Error("Handler Error in UpdateInventId: decoding JSON data", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	invent, err := models.NewInventoryItem(inputInvent.IngredientID, inputInvent.Name, inputInvent.Unit, inputInvent.Quantity)
	if err != nil {
		slog.Error("Handler Error in UpdateInventId: Invalid input data", "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	invent.IngredientID = id

	if err := h.inventServ.UpdateInventIdServ(*invent); err != nil {
		var status int
		if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in UpdateInventId: updating inventory", "error", err)
		writeError(w, err.Error(), status)
		return
	}
	slog.Info("Invent updated successfully", "inventID", invent.IngredientID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *InventHandler) DeleteInventId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.inventServ.DeleteInventIdServ(id); err != nil {
		var status int
		if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in DeleteInventId: deleting inventory by ID", "error", err)
		writeError(w, err.Error(), status)
		return
	}
	w.WriteHeader(http.StatusOK)
	slog.Info("Invent deleted successfully")
}
