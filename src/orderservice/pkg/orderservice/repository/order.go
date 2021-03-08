package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"orderservice/pkg/orderservice/model"
)

type orderRepository struct {
	db *sql.DB
}

func (o *orderRepository) Add(order model.Order) error {
	return o.withTx(func(tx *sql.Tx, ctx context.Context, closeTx func(error) error) error {
		_, err := tx.ExecContext(ctx, "INSERT INTO `order` (`order_id`, `cost`, `created_at`, `updated_at`, `deleted_at`) VALUES (UUID_TO_BIN(?), ?, ?, ?, NULL)", order.ID, order.Cost, order.OrderedAt, order.OrderedAt)
		if err != nil {
			return closeTx(err)
		}

		for _, item := range order.MenuItems {
			_, err = tx.ExecContext(ctx, "INSERT INTO order_item (order_id, menu_item_id, quantity) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?)", order.ID, item.ID, item.Quantity)
			if err != nil {
				return closeTx(err)
			}
		}

		return closeTx(nil)
	})
}

func (o *orderRepository) Delete(id uuid.UUID) error {
	return o.withTx(func(tx *sql.Tx, ctx context.Context, closeTx func(error) error) error {
		_, err := tx.ExecContext(ctx, "UPDATE `order` SET deleted_at = NOW() WHERE BIN_TO_UUID(order_id) = ?", id)

		return closeTx(err)
	})
}

func NewOrderRepository(db *sql.DB) model.OrderRepository {
	return &orderRepository{db: db}
}
