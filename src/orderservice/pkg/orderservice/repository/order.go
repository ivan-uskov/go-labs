package repository

import (
	"context"
	"database/sql"
	"orderservice/pkg/orderservice/model"
)

type orderRepository struct {
	db *sql.DB
}

func (o *orderRepository) Add(order model.Order) error {
	return o.withTx(func(tx *sql.Tx, ctx context.Context, closeTx func(error) error) error {
		orderId, err := order.ID.MarshalBinary()
		if err != nil {
			return closeTx(err)
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO `order` (`order_id`, `cost`, `created_at`, `updated_at`, `deleted_at`) VALUES (?, ?, ?, ?, NULL)", orderId, order.Cost, order.OrderedAt, order.OrderedAt)
		if err != nil {
			return closeTx(err)
		}

		for _, item := range order.MenuItems {
			itemId, err := item.ID.MarshalBinary()
			if err != nil {
				return closeTx(err)
			}

			_, err = tx.ExecContext(ctx, "INSERT INTO order_item (order_id, menu_item_id, quantity) VALUES (?, ?, ?)", orderId, itemId, item.Quantity)
			if err != nil {
				return closeTx(err)
			}
		}

		return closeTx(nil)
	})
}

func NewOrderRepository(db *sql.DB) model.OrderRepository {
	return &orderRepository{db: db}
}
