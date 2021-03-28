package service

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"orderservice/pkg/orderservice/application/data"
	"orderservice/pkg/orderservice/model"
	"time"
)

type AddOrderRequest struct {
	MenuItems []data.MenuItem `json:"menuItems"`
}

type UpdateOrderRequest struct {
	MenuItems []data.MenuItem `json:"menuItems"`
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

func validateOrderItems(reqItems []data.MenuItem) ([]model.MenuItem, error) {
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
		log.Debug(err)
		return fmt.Errorf("invalid uuid: %s", id)
	}

	err = os.repo.Delete(uid)
	if err != nil {
		log.Error(err)
		return data.InternalError
	}

	return nil
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
		return data.InternalError
	}

	return nil
}

func (os *orderService) Update(id string, r UpdateOrderRequest) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Debug(err)
		return fmt.Errorf("invalid uuid: %s", id)
	}

	o, err := os.repo.Get(uid)
	if err != nil {
		log.Error(err)
		return data.InternalError
	}

	if o == nil {
		return fmt.Errorf("order %s not found", id)
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
		return data.InternalError
	}

	return nil
}
