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

type UpdateOrderRequest struct {
	MenuItems []MenuItem `json:"menuItems"`
}

type orderService struct {
	repo model.OrderRepository
}

type OrderService interface {
	Add(r AddOrderRequest) error
	Update(id string, r UpdateOrderRequest) error
	Delete(id string) error
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

func calculateCost(items []model.MenuItem) int {
	return 42 * len(items)
}

func (os *orderService) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return os.repo.Delete(uid)
}

func (os *orderService) Add(r AddOrderRequest) error {
	items, err := validateOrderItems(r.MenuItems)
	if err != nil {
		return err
	}

	err = os.repo.Add(model.Order{
		ID:        uuid.New(),
		MenuItems: items,
		Cost:      calculateCost(items),
		OrderedAt: time.Now(),
	})

	if err != nil {
		log.Error(err)
		return errors.New("server error")
	}

	return nil
}

func (os *orderService) Update(id string, r UpdateOrderRequest) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
		return errors.New("server error")
	}

	o, err := os.repo.Get(uid)
	if err != nil {
		log.Error(err)
		return errors.New("server error")
	}

	if o == nil {
		return errors.New("order not found")
	}

	items, err := validateOrderItems(r.MenuItems)
	if err != nil {
		return err
	}

	o.MenuItems = items
	o.Cost = calculateCost(items)

	err = os.repo.Update(*o)

	if err != nil {
		log.Error(err)
		return errors.New("server error")
	}

	return nil
}
