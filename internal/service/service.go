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

func (s *Service) Process() error {
	for _, str := range s.strategies {
		b, err := s.provider.Balance()
		if err != nil {
			return err
		}
		for _, st := range str.Steps {
			o, err := s.provider.IsOpenOrder(str.Pair, st.Type)
			if err != nil {
				return err
			}
			if b > 0 && !o {
				err := s.provider.AddOrder(str.Pair, st.Type, "limit", "", "")
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type Trade struct {
	airtable.Record
	Fields models.Trade
}

func (s *Service) Trades() error {
	trades := s.client.Table("Trades")

	//get all existing db trades
	tdb := []Trade{}
	trades.List(&tdb, nil)

	t, err := s.provider.Trades()
	if err != nil {
		return err
	}
	for _, trade := range t {
		if exists(trade, tdb) {
			continue
		}
		tt := &Trade{
			Fields: trade,
		}
		err = trades.Create(tt)
		if err != nil {
			return err
		}
	}
	return nil
}

func exists(trade models.Trade, trades []Trade) bool {
	for _, t := range trades {
		if t.Fields.OrderID == trade.OrderID {
			return true
		}
	}
	return false
}
