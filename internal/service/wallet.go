package service

import (
	"fmt"
	"os"
)

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

func (s *Service) ReadFromFile(filename string) error {
	err := s.provider.ReadFromFile(filename)
	if err != nil {
		return err
	}
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("unable to remove %s: %w", filename, err)
	}
	return nil
}
