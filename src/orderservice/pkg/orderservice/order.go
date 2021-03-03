package orderservice

type order struct {
	ID        string     `json:"id"`
	MenuItems []menuItem `json:"menuItems"`
}

type orderDetails struct {
	order
	Time int `json:"orderedAtTimestamp"`
	Cost int `json:"cost"`
}

type ordersList struct {
	Orders []order `json:"orders"`
}

type menuItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}
