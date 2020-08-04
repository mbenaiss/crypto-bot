package service

import (
	"github.com/brianloveswords/airtable"
	"github.com/mbenaiss/crypto-bot/internal/provider"
	"github.com/mbenaiss/crypto-bot/models"
)

type Service struct {
	provider   provider.Provider
	strategies []models.Strategy
	client     airtable.Client
}

func New(c airtable.Client, p provider.Provider, s ...models.Strategy) *Service {
	return &Service{
		provider:   p,
		strategies: s,
		client:     c,
	}
}
