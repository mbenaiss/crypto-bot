package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/brianloveswords/airtable"
	"github.com/mbenaiss/crypto-bot/cmd/server"
	"github.com/mbenaiss/crypto-bot/config"
	"github.com/mbenaiss/crypto-bot/internal/provider/kraken"
	"github.com/mbenaiss/crypto-bot/internal/service"
	"github.com/mbenaiss/crypto-bot/models"
)

var (
	g errgroup.Group
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
	svc := service.New(client, k, models.Strategy{})

	srv := server.New(c)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		g.Go(func() error {
			err := srv.StartHTTP(svc)
			if err != nil && err != http.ErrServerClosed {
				zap.S().Panicf("an error was occurred", err)
				<-quit
			}
			return err
		})

		g.Go(func() error {
			err := srv.StartHealthz()
			if err != nil && err != http.ErrServerClosed {
				zap.S().Panicf("an error was occurred", err)
				<-quit
			}
			return err
		})

		if err := g.Wait(); err != nil {
			zap.S().Panicf("an error was occurred", err)
			<-quit
		}
	}()

	<-quit
	zap.S().Info("Shutdown Server ...")
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Panicf("Server Shutdown:", err)
	}
	zap.S().Info("Server exiting")
}

func Ticker(c *config.Config, ser *service.Service) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	ticker := time.NewTicker(time.Duration(c.TradeTicker) * time.Second)
	go func() {
		for range ticker.C {
			err := ser.Trades()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	wg.Wait()
}
