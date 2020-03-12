package main

import (
	"log"
	"time"

	"github.com/brianloveswords/airtable"
	"github.com/mbenaiss/crypto-bot/config"
	"github.com/mbenaiss/crypto-bot/internal/provider/kraken"
	"github.com/mbenaiss/crypto-bot/internal/service"
	"github.com/mbenaiss/crypto-bot/models"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatalf("unable to init config %+v", err)
	}
	client := airtable.Client{
		APIKey: c.AirtableKey,
		BaseID: c.AirtableBase,
	}
	k := kraken.New(c.KrakenKey, c.KrakenSecret, "ZEUR")

	ser := service.New(client, k, models.Strategy{})

	ticker := time.NewTicker(time.Duration(c.TradeTicker) * time.Second)
	go func() {
		for range ticker.C {
			err := ser.Trades()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}
