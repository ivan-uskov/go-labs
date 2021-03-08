package service

type Order struct {
	ID        string     `json:"id"`
	MenuItems []MenuItem `json:"menuItems"`
}

type OrderInfo struct {
	Order
	Time int `json:"orderedAtTimestamp"`
	Cost int `json:"cost"`
}

type OrdersList struct {
	Orders []Order `json:"getOrdersList"`
}

type MenuItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}
