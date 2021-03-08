package service

type OrderQueryService interface {
	GetOrders() (*OrdersList, error)
	GetOrderInfo(id string) (*OrderInfo, error)
}
