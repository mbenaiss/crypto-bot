package models

import "time"

type Trade struct {
	OrderID string    `json:"OrderID"`
	Crypto  string    `json:"Crypto"`
	Time    time.Time `json:"Date"`
	Type    OrderType `json:"Type"`
	Price   float64   `json:"Price"`
	Amount  float64   `json:"Amount"`
	Fee     float64   `json:"Fee"`
	Volume  float64   `json:"Volume"`
}

type OrderType string

const (
	Buy  OrderType = "buy"
	Sell OrderType = "sell"
)
