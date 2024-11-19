package service

import (
	"fmt"
	"hot-coffee/internal/customErrors"
	"hot-coffee/internal/models"
	"log/slog"
	"sort"
)

type OrderRepo interface {
	GetOrdersRepo() (map[string]models.Order, error)
	UpdateOrdersRepo(ordersMap map[string]models.Order) error
}

type MenuRepoForOrder interface {
	GetMenusRepo() (map[string]models.MenuItem, error)
}

type InventRepoForOrder interface {
	GetInventsRepo() (map[string]models.InventoryItem, error)
	UpdateInventsRepo(inventMap map[string]models.InventoryItem) error
}

type OrderServiceImpl struct {
	orderRepo  OrderRepo
	menuRepo   MenuRepoForOrder
	inventRepo InventRepoForOrder
}

func NewOrderServiceImpl(oR OrderRepo, mR MenuRepoForOrder, iR InventRepoForOrder) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepo:  oR,
		menuRepo:   mR,
		inventRepo: iR,
	}
}

func (s *OrderServiceImpl) CreateOrderService(newOrder models.Order) (models.TotalPrice, error) {
	menuMap, err := s.validateOrder(newOrder.Items)
	if err != nil {
		slog.Error("Order Service in CreateOrderService")
		return models.TotalPrice{}, err
	}

	orderMap, err := s.orderRepo.GetOrdersRepo()
	if err != nil {
		slog.Error("Order Service in CreateOrderService")
		return models.TotalPrice{}, err
	}
	orderId, err := getNewOrderID(orderMap)
	if err != nil {
		slog.Error("Order Service in CreateOrderService")
		return models.TotalPrice{}, err
	}

	newOrder.ID = orderId
	newOrder.Status = "open"

	orderMap[newOrder.ID] = newOrder
	if err := s.orderRepo.UpdateOrdersRepo(orderMap); err != nil {
		slog.Error("Order Service in CreateOrderService")
		return models.TotalPrice{}, err
	}

	totalPrice := models.NewTotalPrice()
	for _, orderItem := range newOrder.Items {
		totalPrice.TotalSale += float64(orderItem.Quantity) * menuMap[orderItem.ProductID].Price
	}

	return *totalPrice, nil
}

func (s *OrderServiceImpl) GetOrdersService() ([]models.Order, error) {
	orderMap, err := s.orderRepo.GetOrdersRepo()
	if err != nil {
		slog.Error("Order Service in GetOrdersService")
		return nil, err
	}
	var orders []models.Order
	for _, item := range orderMap {
		orders = append(orders, item)
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID < orders[j].ID
	})

	return orders, nil
}

func (s *OrderServiceImpl) GetOrderByIdService(id string) (models.Order, error) {
	orderMap, err := s.orderRepo.GetOrdersRepo()
	if err != nil {
		slog.Error("Order Service in GetOrderByIdService")
		return models.Order{}, err
	}

	order, exists := orderMap[id]
	if !exists {
		slog.Error("Order Service in GetOrderByIdService")
		return models.Order{}, fmt.Errorf("%w", customErrors.ErrNotExistConflict)
	}

	return order, nil
}

func (s *OrderServiceImpl) UpdateOrderByIdService(updateOrder models.Order) (models.TotalPrice, error) {
	_, err := s.validateOrder(updateOrder.Items)
	if err != nil {
		slog.Error("Order Service in UpdateOrderByIdService")
		return models.TotalPrice{}, err
	}

	orderMap, err := s.orderRepo.GetOrdersRepo()
	if err != nil {
		slog.Error("Order Service in UpdateOrderByIdService")
		return models.TotalPrice{}, err
	}

	order, exists := orderMap[updateOrder.ID]
	if !exists {
		slog.Error("Order Service in UpdateOrderByIdService")
		return models.TotalPrice{}, fmt.Errorf("%w", customErrors.ErrNotExistConflict)
	}

	if order.Status == "closed" {
		slog.Error("Order Service in UpdateOrderByIdService")
		return models.TotalPrice{}, fmt.Errorf("%w", customErrors.ErrOrderClosed)
	}

	updateOrder.Status = order.Status
	updateOrder.CreatedAt = order.CreatedAt
	updateOrder.Items = append(updateOrder.Items, order.Items...)

	orderMap[updateOrder.ID] = updateOrder

	if err := s.orderRepo.UpdateOrdersRepo(orderMap); err != nil {
		slog.Error("Order Service in UpdateOrderByIdService")
		return models.TotalPrice{}, err
	}

	menuMap, err := s.menuRepo.GetMenusRepo()
	if err != nil {
		slog.Error("Order Service in UpdateOrderByIdService")
		return models.TotalPrice{}, err
	}

	totalPrice := models.NewTotalPrice()
	for _, orderItem := range updateOrder.Items {
		totalPrice.TotalSale += float64(orderItem.Quantity) * menuMap[orderItem.ProductID].Price
	}

	return *totalPrice, nil
}

func (s *OrderServiceImpl) DeleteOrderByIdService(id string) error {
	orderMap, err := s.orderRepo.GetOrdersRepo()
	if err != nil {
		slog.Error("Order Service in DeleteOrderByIdService")
		return err
	}

	_, exists := orderMap[id]
	if !exists {
		slog.Error("Order Service in DeleteOrderByIdService")
		return fmt.Errorf("%w", customErrors.ErrNotExistConflict)
	}

	delete(orderMap, id)
	return s.orderRepo.UpdateOrdersRepo(orderMap)
}

func (s *OrderServiceImpl) CloseOrderByIdService(id string) error {
	orderMap, err := s.orderRepo.GetOrdersRepo()
	if err != nil {
		slog.Error("Order Service in CloseOrderByIdService")
		return err
	}

	order, exists := orderMap[id]
	if !exists {
		slog.Error("Order Service in CloseOrderByIdService")
		return fmt.Errorf("%w", customErrors.ErrNotExistConflict)
	}

	order.Status = "closed"
	orderMap[id] = order

	return s.orderRepo.UpdateOrdersRepo(orderMap)
}

func (s *OrderServiceImpl) validateOrder(orderItems []models.OrderItem) (map[string]models.MenuItem, error) {
	menuMap, err := s.menuRepo.GetMenusRepo()
	if err != nil {
		slog.Error("Order Service in validateOrder")
		return nil, err
	}
	inventoryMap, err := s.inventRepo.GetInventsRepo()
	if err != nil {
		slog.Error("Order Service in validateOrder")
		return nil, err
	}

	requiredIngredients := make(map[string]float64)
	for _, orderItem := range orderItems {
		menuItem, exists := menuMap[orderItem.ProductID]
		if !exists {
			return nil, customErrors.ErrNotExistConflict
		}

		for _, ingredient := range menuItem.Ingredients {
			requiredIngredients[ingredient.IngredientID] += ingredient.Quantity * float64(orderItem.Quantity)
		}
	}

	for ingredientID, requiredQuantity := range requiredIngredients {
		inventoryItem, exists := inventoryMap[ingredientID]
		if !exists || inventoryItem.Quantity < requiredQuantity {
			slog.Error("Insufficient ingredient in inventory", "ingredientID", ingredientID)
			return nil, fmt.Errorf("insufficient ingredient: %s", ingredientID)
		}
	}

	for ingredientID, usedQuantity := range requiredIngredients {
		inventoryItem := inventoryMap[ingredientID]
		inventoryItem.Quantity -= usedQuantity
		inventoryMap[ingredientID] = inventoryItem
	}

	if err := s.inventRepo.UpdateInventsRepo(inventoryMap); err != nil {
		slog.Error("Order Service in validateOrder")
		return nil, err
	}

	return menuMap, nil
}
