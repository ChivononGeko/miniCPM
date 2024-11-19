package service

import (
	"hot-coffee/internal/customErrors"
	"hot-coffee/internal/models"
	"sort"
)

type OrderRepoForReport interface {
	GetOrdersRepo() (map[string]models.Order, error)
	UpdateOrdersRepo(ordersMap map[string]models.Order) error
}

type MenuRepoForReports interface {
	GetMenusRepo() (map[string]models.MenuItem, error)
}

type ReportsServiceImplementation struct {
	ordersRepository OrderRepoForReport
	menuRepository   MenuRepoForReports
}

func NewReportsService(or OrderRepoForReport, mr MenuRepoForReports) *ReportsServiceImplementation {
	return &ReportsServiceImplementation{
		ordersRepository: or,
		menuRepository:   mr,
	}
}

func (rs *ReportsServiceImplementation) TotalSalesReportService() (models.TotalPrice, error) {
	ordersMap, err := rs.ordersRepository.GetOrdersRepo()
	if err != nil {
		return models.TotalPrice{}, err
	}

	menuMap, err := rs.menuRepository.GetMenusRepo()
	if err != nil {
		return models.TotalPrice{}, err
	}

	var totalSale models.TotalPrice
	for _, items := range ordersMap {
		for _, item := range items.Items {
			totalSale.TotalSale += menuMap[item.ProductID].Price * float64(item.Quantity)
		}
	}

	return totalSale, nil
}

func (rs *ReportsServiceImplementation) PopularItemsReportService() ([]models.PopularItem, error) {
	ordersMap, err := rs.ordersRepository.GetOrdersRepo()
	if err != nil {
		return nil, err
	}

	popularItemsMap := make(map[string]models.PopularItem)

	for _, items := range ordersMap {
		for _, item := range items.Items {
			_, exists := popularItemsMap[item.ProductID]
			if exists {
				tempItem := popularItemsMap[item.ProductID]
				tempItem.QuantityOfSales += item.Quantity
				popularItemsMap[item.ProductID] = tempItem
			} else {
				popularItemsMap[item.ProductID] = models.PopularItem{ItemName: item.ProductID, QuantityOfSales: item.Quantity}
			}
		}
	}

	var popularItems []models.PopularItem
	for _, item := range popularItemsMap {
		popularItems = append(popularItems, item)
	}
	sort.Slice(popularItems, func(i, j int) bool {
		return popularItems[i].QuantityOfSales > popularItems[j].QuantityOfSales
	})

	if len(popularItems) > 0 {
		return popularItems[:1], nil
	} else {
		return nil, customErrors.ErrNotExistConflict
	}
}
