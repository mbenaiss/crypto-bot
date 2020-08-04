package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	KrakenKey    string `envconfig:"KRAKEN_KEY"`
	KrakenSecret string `envconfig:"KRAKEN_SECRET"`
	AirtableKey  string `envconfig:"AIRTABLE_KEY"`
	AirtableBase string `envconfig:"AIRTABLE_BASE"`
	TradeTicker  int    `envconfig:"TRADE_TICKER"`
	HttpPort     int    `envconfig:"HTTP_PORT"`
	HealthzPort  int    `envconfig:"HEALTHZ_PORT"`
}

//New retrun new instance of Config
func New() (*Config, error) {
	var c Config
	err := godotenv.Load()
	if err != nil {
		log.Println("unable to load .env file")
	}
	err = envconfig.Process("", &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
