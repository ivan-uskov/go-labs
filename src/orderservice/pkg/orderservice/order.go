package orderservice

type Order struct {
	ID        string     `json:"id"`
	MenuItems []MenuItem `json:"menuItems"`
}

type OrderDetails struct {
	Order
	Time int `json:"orderedAtTimestamp"`
	Cost int `json:"cost"`
}

type OrdersList struct {
	Orders []Order `json:"orders"`
}

type MenuItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}
