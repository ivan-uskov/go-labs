package query

import "orderservice/pkg/orderservice/application/data"

type OrderQueryService interface {
	GetOrders() (*data.OrdersList, error)
	GetOrderInfo(id string) (*data.OrderInfo, error)
}
