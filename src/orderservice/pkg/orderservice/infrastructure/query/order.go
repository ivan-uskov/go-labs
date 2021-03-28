package query

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"orderservice/pkg/orderservice/application/data"
	"orderservice/pkg/orderservice/application/query"
	"strconv"
	"strings"
	"time"
)

type orderQueryService struct {
	db *sql.DB
}

func NewOrderQueryService(db *sql.DB) query.OrderQueryService {
	return &orderQueryService{db: db}
}

func parseMenuItems(itemsStr string) ([]data.MenuItem, error) {
	if len(itemsStr) == 0 {
		return make([]data.MenuItem, 0), nil
	}

	items := strings.Split(itemsStr, ",")
	result := make([]data.MenuItem, len(items))
	for i, pairStr := range items {
		pair := strings.Split(pairStr, "=")
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid pair: %s", pairStr)
		}

		quantity, err := strconv.Atoi(pair[1])
		if err != nil {
			return nil, fmt.Errorf("invalid quantity: %s", pair[1])
		}

		result[i] = data.MenuItem{ID: pair[0], Quantity: quantity}
	}

	return result, nil
}

func parseOrder(r *sql.Rows) (*data.OrderInfo, error) {
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

	return &data.OrderInfo{
		ID:        orderId,
		MenuItems: menuItems,
		OrderedAt: createdAt,
		Cost:      cost,
	}, nil
}

func (qs *orderQueryService) GetOrders() (*data.OrdersList, error) {
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
		log.Error(err)
		return nil, data.InternalError
	}
	defer rows.Close()

	orders := make([]data.OrderInfo, 0)
	for rows.Next() {
		order, err := parseOrder(rows)
		if err != nil {
			log.Error(err)
			return nil, data.InternalError
		}

		orders = append(orders, *order)
	}

	return &data.OrdersList{Orders: orders}, nil
}

func (qs *orderQueryService) GetOrderInfo(id string) (*data.OrderInfo, error) {
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
		log.Error(err)
		return nil, data.InternalError
	}
	defer rows.Close()

	if rows.Next() {
		order, err := parseOrder(rows)
		if err != nil {
			log.Error(err)
			return nil, data.InternalError
		}

		return order, nil
	}

	return nil, nil // not found
}
