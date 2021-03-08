package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"orderservice/pkg/orderservice/model"
	"time"
)

type AddOrderRequest struct {
	MenuItems []MenuItem `json:"menuItems"`
}

type orderService struct {
	repo model.OrderRepository
}

type OrderService interface {
	AddOrder(r AddOrderRequest) error
}

func NewOrderService(repo model.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func validateOrderItems(reqItems []MenuItem) ([]model.MenuItem, error) {
	itemIds := map[string]bool{}
	items := make([]model.MenuItem, len(reqItems))
	for i, item := range reqItems {
		itemId, err := uuid.Parse(item.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %s", item.ID)
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("item: %s has invalid quantity, accepted only positive numbers", item.ID)
		}
		if _, found := itemIds[item.ID]; found {
			return nil, fmt.Errorf("has duplicate menu items: %s", item.ID)
		} else {
			itemIds[item.ID] = true
		}

		items[i].ID = itemId
		items[i].Quantity = item.Quantity
	}

	return items, nil
}

func (os *orderService) AddOrder(r AddOrderRequest) error {
	items, err := validateOrderItems(r.MenuItems)
	if err != nil {
		return err
	}

	err = os.repo.Add(model.Order{
		ID:        uuid.New(),
		MenuItems: items,
		Cost:      42,
		OrderedAt: time.Now(),
	})

	if err != nil {
		log.Error(err)
		return errors.New("server error")
	}

	return nil
}
