package service

import "time"

type MenuItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderInfo struct {
	ID        string     `json:"id"`
	MenuItems []MenuItem `json:"menuItems"`
	OrderedAt time.Time  `json:"orderedAtTimestamp"`
	Cost      int        `json:"cost"`
}

type OrdersList struct {
	Orders []OrderInfo `json:"orders"`
}
