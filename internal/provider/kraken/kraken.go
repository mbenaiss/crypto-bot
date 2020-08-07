package kraken

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	krakenapi "github.com/beldur/kraken-go-api-client"

	"github.com/mbenaiss/crypto-bot/internal/provider"
	"github.com/mbenaiss/crypto-bot/models"
	"github.com/mbenaiss/crypto-bot/pkg/csv"
)

type kraken struct {
	api   *krakenapi.KrakenApi
	asset string
	csv   *csv.Client
}

var csvColumns = []string{"time", "type", "asset", "amount", "fee", "balance"}

func New(key, secret string, asset string) provider.Provider {
	api := krakenapi.New(key, secret)
	c := csv.New(',', csvColumns)
	return &kraken{
		api:   api,
		asset: asset,
		csv:   c,
	}
}

func (k *kraken) Name() provider.ProviderName {
	return provider.Kraken
}

func (k *kraken) IsTradable() bool {
	return true
}

func (k *kraken) Balance() (float64, error) {
	args := map[string]string{
		"asset": k.asset,
	}
	a, err := k.api.TradeBalance(args)
	if err != nil {
		return 0, err
	}
	orders, err := k.api.OpenOrders(args)
	if err != nil {
		return 0, err
	}
	total := 0.0
	for _, o := range orders.Open {
		if o.Description.Type == "buy" {
			v, err := strconv.ParseFloat(o.Volume, 64)
			if err != nil {
				return 0, err
			}
			p, err := strconv.ParseFloat(o.Description.PrimaryPrice, 64)
			if err != nil {
				return 0, err
			}
			total += p * v
		}
	}
	return a.FreeMargin - total, nil
}

func (k *kraken) IsOpenOrder(p string, t string) (bool, error) {
	orders, err := k.api.OpenOrders(map[string]string{
		"asset": k.asset,
	})
	if err != nil {
		return false, err
	}
	for _, o := range orders.Open {
		if o.Description.AssetPair == p && o.Description.Type == t {
			return true, nil
		}
	}
	return false, nil
}

func (k *kraken) AddOrder(pair, direction, orderType, price, volume string) error {
	resp, err := k.api.AddOrder(pair, direction, orderType, volume, map[string]string{
		"asset": k.asset,
		"price": price,
	})
	if err != nil {
		return err
	}
	if len(resp.TransactionIds) <= 0 {
		return fmt.Errorf("unable to add order")
	}
	return nil
}

func (k *kraken) Price(pair string) (float64, error) {
	_, err := k.api.Ticker()
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (k *kraken) Trades() ([]models.Trade, error) {
	h, err := k.api.TradesHistory(0, 0, map[string]string{
		"asset": k.asset,
	})
	if err != nil {
		return nil, err
	}
	res := []models.Trade{}
	for id, o := range h.Trades {
		amount := o.Cost
		if o.Type == "sell" {
			amount = amount * -1
		}
		res = append(res, models.Trade{
			OrderID: id,
			Crypto:  getCrypto(o.AssetPair),
			Time:    time.Unix(int64(o.Time), 0),
			Type:    models.OrderType(o.Type),
			Price:   o.Price,
			Amount:  amount,
			Fee:     o.Fee,
			Volume:  o.Volume,
		})
	}
	return res, nil
}

func (k *kraken) ReadFromFile(filename string) ([]models.Trade, error) {
	r, err := k.csv.Read(filename)
	if err != nil {
		return nil, err
	}
	trades := make([]models.Trade, 0, len(r))
	for _, trade := range r {
		fee, err := strconv.ParseFloat(trade["fee"], 64)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(trade["balance"], 64)
		if err != nil {
			return nil, err
		}
		price, err := strconv.ParseFloat(trade["amount"], 64)
		if err != nil {
			return nil, err
		}
		ti, err := time.Parse(time.RFC3339, trade["time"])
		if err != nil {
			return nil, err
		}
		t := models.Buy //trade["type"]
		trades = append(trades, models.Trade{
			Crypto:  trade["asset"],
			Amount:  amount,
			Fee:     fee,
			OrderID: trade["amount"],
			Price:   price,
			Time:    ti,
			Type:    t,
		})
	}
	return trades, nil
}

func getCrypto(pair string) string {
	r := strings.ReplaceAll(pair, "ZEUR", "")
	r = strings.ReplaceAll(r, "EUR", "")
	if strings.HasPrefix(r, "X") {
		return strings.TrimPrefix(r, "X")
	}
	return r
}
