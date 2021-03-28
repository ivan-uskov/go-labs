package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"orderservice/pkg/orderservice/model"
	"strconv"
	"strings"
	"time"
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

func (o *orderRepository) Update(order model.Order) error {
	return o.withTx(func(tx *sql.Tx, ctx context.Context, closeTx func(error) error) error {
		_, err := tx.ExecContext(ctx, "UPDATE `order` SET cost = ?, updated_at = NOW() WHERE BIN_TO_UUID(order_id) = ?", order.Cost, order.ID)
		if err != nil {
			return closeTx(err)
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM order_item WHERE BIN_TO_UUID(order_id) = ?", order.ID)
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

func (o *orderRepository) Get(id uuid.UUID) (*model.Order, error) {
	rows, err := o.db.Query(""+
		"SELECT "+
		"BIN_TO_UUID(o.order_id) AS order_id, "+
		"o.cost, "+
		"o.created_at, "+
		"IFNULL(GROUP_CONCAT(CONCAT(BIN_TO_UUID(oi.menu_item_id), '=', oi.quantity)), '') AS items "+
		"FROM `order` o "+
		"LEFT JOIN order_item oi ON (o.order_id = oi.order_id) "+
		"WHERE o.deleted_at IS NULL AND BIN_TO_UUID(o.order_id) = ? "+
		"GROUP BY o.order_id", id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		order, err := parseOrder(rows)
		if err != nil {
			return nil, err
		}

		return order, nil
	}

	return nil, nil // not found
}

func parseMenuItems(itemsStr string) ([]model.MenuItem, error) {
	if len(itemsStr) == 0 {
		return make([]model.MenuItem, 0), nil
	}

	items := strings.Split(itemsStr, ",")
	result := make([]model.MenuItem, len(items))
	for i, pairStr := range items {
		pair := strings.Split(pairStr, "=")
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid pair: %s", pairStr)
		}

		quantity, err := strconv.Atoi(pair[1])
		if err != nil {
			return nil, fmt.Errorf("invalid quantity: %s", pair[1])
		}

		itemId, err := uuid.Parse(pair[0])
		if err != nil {
			return nil, err
		}

		result[i] = model.MenuItem{ID: itemId, Quantity: quantity}
	}

	return result, nil
}

func parseOrder(r *sql.Rows) (*model.Order, error) {
	var orderId string
	var cost int
	var createdAt time.Time
	var items string

	err := r.Scan(&orderId, &cost, &createdAt, &items)
	if err != nil {
		return nil, err
	}

	orderUid, err := uuid.Parse(orderId)
	if err != nil {
		return nil, err
	}

	menuItems, err := parseMenuItems(items)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		ID:        orderUid,
		MenuItems: menuItems,
		OrderedAt: createdAt,
		Cost:      cost,
	}, nil
}
