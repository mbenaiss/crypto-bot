package service

import (
	"fmt"

	"github.com/brianloveswords/airtable"

	"github.com/mbenaiss/crypto-bot/internal/provider"
	"github.com/mbenaiss/crypto-bot/internal/provider/kraken"
	"github.com/mbenaiss/crypto-bot/models"
)

type Service struct {
	providers  []provider.Provider
	strategies []models.Strategy
	db         airtable.Client
}

func New(c airtable.Client, p []provider.Provider, s ...models.Strategy) *Service {
	return &Service{
		providers:  p,
		strategies: s,
		db:         c,
	}
}

func (s *Service) AddProvider(pr models.Provider) error {
	c, err := s.getProviderFromName(provider.ToProviderName(pr.Name))
	if err != nil {
		return err
	}

	var newProvider provider.Provider
	switch c.Name() {
	case provider.Kraken:
		newProvider = kraken.New(pr.Key, pr.Secret, "")
	}
	s.providers = append(s.providers, newProvider)
	return nil
}

func (s *Service) getProviderFromName(name provider.ProviderName) (provider.Provider, error) {
	for _, p := range s.providers {
		if p.Name() == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("unable to find provider")
}
