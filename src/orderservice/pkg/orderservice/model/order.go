package model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID        uuid.UUID
	MenuItems []MenuItem
	Cost      int
	OrderedAt time.Time
}

type MenuItem struct {
	ID       uuid.UUID
	Quantity int
}

type OrderRepository interface {
	Add(order Order) error
	Update(order Order) error
	Delete(id uuid.UUID) error

	Get(id uuid.UUID) (*Order, error)
}
