package service

import (
	"fmt"
	"os"

	"github.com/mbenaiss/crypto-bot/internal/provider"
)

func (s *Service) Process(providerName provider.ProviderName) error {
	p, err := s.getProviderFromName(providerName)
	if err != nil {
		return err
	}
	for _, str := range s.strategies {
		b, err := p.Balance()
		if err != nil {
			return err
		}
		for _, st := range str.Steps {
			o, err := p.IsOpenOrder(str.Pair, st.Type)
			if err != nil {
				return err
			}
			if b > 0 && !o {
				err := p.AddOrder(str.Pair, st.Type, "limit", "", "")
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Service) ReadFromFile(filename string, name provider.ProviderName) error {
	p, err := s.getProviderFromName(name)
	if err != nil {
		return err
	}
	trades, err := p.ReadFromFile(filename)
	if err != nil {
		return err
	}
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("unable to remove %s: %w", filename, err)
	}
	fmt.Println(trades)
	return nil
}
