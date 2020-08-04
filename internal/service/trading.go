package service

import (
	"github.com/brianloveswords/airtable"
	"github.com/mbenaiss/crypto-bot/models"
)

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
