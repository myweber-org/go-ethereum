package main

import (
	"fmt"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
)

type ExchangeRates struct {
	rates map[Currency]float64
}

func NewExchangeRates() *ExchangeRates {
	return &ExchangeRates{
		rates: map[Currency]float64{
			USD: 1.0,
			EUR: 0.85,
			GBP: 0.73,
			JPY: 110.0,
		},
	}
}

func (er *ExchangeRates) Convert(amount float64, from Currency, to Currency) (float64, error) {
	fromRate, okFrom := er.rates[from]
	toRate, okTo := er.rates[to]

	if !okFrom || !okTo {
		return 0, fmt.Errorf("unsupported currency")
	}

	baseAmount := amount / fromRate
	return baseAmount * toRate, nil
}

func (er *ExchangeRates) UpdateRate(currency Currency, rate float64) {
	er.rates[currency] = rate
}

func main() {
	rates := NewExchangeRates()

	converted, err := rates.Convert(100.0, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("100 USD = %.2f EUR\n", converted)

	rates.UpdateRate(JPY, 115.0)
	converted, err = rates.Convert(5000.0, JPY, GBP)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("5000 JPY = %.2f GBP\n", converted)
}