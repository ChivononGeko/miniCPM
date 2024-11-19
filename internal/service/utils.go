package service

import (
	"fmt"
	"hot-coffee/internal/models"
	"regexp"
	"strconv"
)

func getNewOrderID(orders map[string]models.Order) (string, error) {
	var maxNum int
	re := regexp.MustCompile(`\d+`)

	for id := range orders {
		match := re.FindString(id)
		if match == "" {
			return "", fmt.Errorf("non valid ID")
		}

		num, err := strconv.Atoi(match)
		if err != nil {
			return "", err
		}

		if num > maxNum {
			maxNum = num
		}
	}

	nextNum := strconv.Itoa(maxNum + 1)
	orderID := "order" + nextNum

	return orderID, nil
}
