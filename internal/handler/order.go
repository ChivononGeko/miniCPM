package handler

import (
	"encoding/json"
	"errors"
	"hot-coffee/internal/customErrors"
	"hot-coffee/internal/models"
	"log/slog"
	"net/http"
	"time"
)

type OrderService interface {
	CreateOrderService(newOrder models.Order) (models.TotalPrice, error)
	GetOrdersService() ([]models.Order, error)
	GetOrderByIdService(id string) (models.Order, error)
	UpdateOrderByIdService(updateOrder models.Order) (models.TotalPrice, error)
	DeleteOrderByIdService(id string) error
	CloseOrderByIdService(id string) error
}

type OrderHandler struct {
	orderService OrderService
}

func NewOrderHandler(os OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: os,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		return
	}

	var inputOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&inputOrder); err != nil {
		slog.Error("Handler Error in CreateOrderHandler: decoding JSON data", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	order, err := models.NewOrder(inputOrder.CustomerName, time.Now(), inputOrder.Items)
	if err != nil {
		slog.Error("Handler Error in CreateOrderHandler: Invalid input data", "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	totalPrice, err := h.orderService.CreateOrderService(*order)
	if err != nil {
		var status int
		if errors.Is(err, customErrors.ErrExistConflict) {
			status = http.StatusConflict
		} else if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in CreateOrderHandler: creating order", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	slog.Info("Order created successfully", "orderID", order.ID)
	writeJSON(w, http.StatusCreated, totalPrice)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	allOrders, err := h.orderService.GetOrdersService()
	if err != nil {
		slog.Error("Handler Error in GetOrders: retrieving all orders", "error", err)
		writeError(w, "Failed to retrieve all orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(allOrders); err != nil {
		slog.Error("Handler Error in GetOrders: encoding JSON data", "error", err)
		writeError(w, "Failed to encode all orders to JSON", http.StatusInternalServerError)

		return
	}

	slog.Info("All orders retrieved successfully")
}

func (h *OrderHandler) GetOrderId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !isValidID(w, id, "GetOrderByIdHandler", "get by id") {
		return
	}

	orderId, err := h.orderService.GetOrderByIdService(id)
	if err != nil {
		var status int
		if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in GetOrderId: retrieving order by ID ", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(orderId); err != nil {
		slog.Error("Handler Error in GetOrderId: encoding JSON data", "error", err)
		writeError(w, "Failed to encode order to JSON", http.StatusInternalServerError)
		return
	}

	slog.Info("Order retrieved successfully", "orderID", id)
}

func (h *OrderHandler) UpdateOrderId(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		return
	}

	id := r.PathValue("id")

	if !isValidID(w, id, "UpdateOrderByIdHandler", "update") {
		return
	}

	var inputOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&inputOrder); err != nil {
		slog.Error("Handler Error in UpdateOrderId: decoding JSON data", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	order, err := models.NewOrder(inputOrder.CustomerName, time.Now(), inputOrder.Items)
	if err != nil {
		slog.Error("Handler Error in UpdateOrderId: Invalid input data", "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.ID = id

	totalPrice, err := h.orderService.UpdateOrderByIdService(*order)
	if err != nil {
		var status int
		if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else if errors.Is(err, customErrors.ErrOrderClosed) {
			status = http.StatusBadRequest
		} else if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in UpdateOrderId: updating order", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	slog.Info("Order updated successfully", "orderID", order.ID)
	writeJSON(w, http.StatusOK, totalPrice)
}

func (h *OrderHandler) DeleteOrderId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !isValidID(w, id, "DeleteOrderId", "delete") {
		return
	}

	if err := h.orderService.DeleteOrderByIdService(id); err != nil {
		var status int
		if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else if errors.Is(err, customErrors.ErrOrderClosed) {
			status = http.StatusBadRequest
		} else {
			status = http.StatusInternalServerError
		}

		slog.Error("Handler Error in DeleteOrderId: deleting order by ID", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Order deleted successfully")
}

func (h *OrderHandler) CloseOrderId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !isValidID(w, id, "CloseOrderId", "close") {
		return
	}

	if err := h.orderService.CloseOrderByIdService(id); err != nil {
		var status int
		if errors.Is(err, customErrors.ErrNotExistConflict) {
			status = http.StatusNotFound
		} else if errors.Is(err, customErrors.ErrOrderClosed) {
			status = http.StatusBadRequest
		} else {
			status = http.StatusInternalServerError
		}
		slog.Error("Handler Error in CloseOrderId: closing order by ID ", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Order closed successfully")
}
