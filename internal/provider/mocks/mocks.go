package mocks

import (
	"github.com/mbenaiss/crypto-bot/internal/provider"
	"github.com/mbenaiss/crypto-bot/models"
)

type mock struct{}

func New() provider.Provider {
	return &mock{}
}

func (m *mock) Balance() (float64, error) {
	return 1000.0, nil
}

func (m *mock) IsOpenOrder(p string, t string) (bool, error) {
	return false, nil
}

func (m *mock) AddOrder(pair, direction, orderType, price, volume string) error {
	return nil
}

func (m *mock) Trades() ([]models.Trade, error) {
	return nil, nil
}

func (m *mock) ReadFromFile(filename string) ([]models.Trade, error) {
	return nil, nil
}

func (m *mock) Name() provider.ProviderName {
	return provider.Mock
}

func (m *mock) IsTradable() bool {
	return false
}
