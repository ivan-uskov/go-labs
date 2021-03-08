package query

import (
	"database/sql"
	"fmt"
	"orderservice/pkg/orderservice/service"
	"strconv"
	"strings"
	"time"
)

type orderQueryService struct {
	db *sql.DB
}

func NewOrderQueryService(db *sql.DB) service.OrderQueryService {
	return &orderQueryService{db: db}
}

func parseMenuItems(itemsStr string) ([]service.MenuItem, error) {
	if len(itemsStr) == 0 {
		return make([]service.MenuItem, 0), nil
	}

	items := strings.Split(itemsStr, ",")
	result := make([]service.MenuItem, len(items))
	for i, pairStr := range items {
		pair := strings.Split(pairStr, "=")
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid pair: %s", pairStr)
		}

		quantity, err := strconv.Atoi(pair[1])
		if err != nil {
			return nil, fmt.Errorf("invalid quantity: %s", pair[1])
		}

		result[i] = service.MenuItem{ID: pair[0], Quantity: quantity}
	}

	return result, nil
}

func parseOrder(r *sql.Rows) (*service.OrderInfo, error) {
	var orderId string
	var cost int
	var createdAt time.Time
	var items string

	err := r.Scan(&orderId, &cost, &createdAt, &items)
	if err != nil {
		return nil, err
	}

	menuItems, err := parseMenuItems(items)
	if err != nil {
		return nil, err
	}

	return &service.OrderInfo{
		ID:        orderId,
		MenuItems: menuItems,
		OrderedAt: createdAt,
		Cost:      cost,
	}, nil
}

func (qs *orderQueryService) GetOrders() (*service.OrdersList, error) {
	rows, err := qs.db.Query("" +
		"SELECT " +
		"BIN_TO_UUID(o.order_id) AS order_id, " +
		"o.cost, " +
		"o.created_at, " +
		"IFNULL(GROUP_CONCAT(CONCAT(BIN_TO_UUID(oi.menu_item_id), '=', oi.quantity)), '') AS items " +
		"FROM `order` o " +
		"LEFT JOIN order_item oi ON (o.order_id = oi.order_id) " +
		"WHERE o.deleted_at IS NULL " +
		"GROUP BY o.order_id")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]service.OrderInfo, 0)
	for rows.Next() {
		order, err := parseOrder(rows)
		if err != nil {
			return nil, err
		}

		orders = append(orders, *order)
	}

	return &service.OrdersList{Orders: orders}, nil
}

func (qs *orderQueryService) GetOrderInfo(id string) (*service.OrderInfo, error) {
	rows, err := qs.db.Query(""+
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
