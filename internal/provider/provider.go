package provider

import "github.com/mbenaiss/crypto-bot/models"

type ProviderName string

const (
	Kraken  ProviderName = "kraken"
	Binance ProviderName = "binance"
)

type Provider interface {
	Balance() (float64, error)
	IsOpenOrder(pair string, orderType string) (bool, error)
	AddOrder(pair, direction, orderType, price, volume string) error
	Trades() ([]models.Trade, error)
	ReadFromFile(filename string) error
}
