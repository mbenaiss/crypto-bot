package provider

import (
	"strings"

	"github.com/mbenaiss/crypto-bot/models"
)

type ProviderName string

const (
	Kraken  ProviderName = "kraken"
	Binance ProviderName = "binance"
	Mock    ProviderName = "mock"
)

func ToProviderName(name string) ProviderName {
	return ProviderName(strings.ToLower(strings.TrimSpace(name)))
}

type Provider interface {
	Balance() (float64, error)
	IsOpenOrder(pair string, orderType string) (bool, error)
	AddOrder(pair, direction, orderType, price, volume string) error
	Trades() ([]models.Trade, error)
	ReadFromFile(filename string) ([]models.Trade, error)
	Name() ProviderName
	IsTradable() bool
}
