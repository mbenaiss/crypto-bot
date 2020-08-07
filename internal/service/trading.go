package service

import (
	"fmt"

	"github.com/brianloveswords/airtable"
	"github.com/mbenaiss/crypto-bot/internal/provider"
	"github.com/mbenaiss/crypto-bot/models"
)

type Trade struct {
	airtable.Record
	Fields models.Trade
}

func (s *Service) Trades(name provider.ProviderName) error {
	trades := s.db.Table("Trades")

	p, err := s.getProviderFromName(name)
	if err != nil {
		return err
	}

	if p.IsTradable() {
		return fmt.Errorf("this provider is not tradable")
	}

	//get all existing db trades
	tdb := []Trade{}
	err = trades.List(&tdb, nil)
	if err != nil {
		return err
	}

	t, err := p.Trades()
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
