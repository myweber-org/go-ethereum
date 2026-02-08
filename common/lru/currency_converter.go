
package main

import (
	"fmt"
	"math"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
)

type ExchangeRates struct {
	rates map[Currency]map[Currency]float64
}

func NewExchangeRates() *ExchangeRates {
	rates := map[Currency]map[Currency]float64{
		USD: {EUR: 0.85, GBP: 0.73, JPY: 110.0},
		EUR: {USD: 1.18, GBP: 0.86, JPY: 129.5},
		GBP: {USD: 1.37, EUR: 1.16, JPY: 150.8},
		JPY: {USD: 0.0091, EUR: 0.0077, GBP: 0.0066},
	}
	return &ExchangeRates{rates: rates}
}

func (er *ExchangeRates) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	rateMap, exists := er.rates[from]
	if !exists {
		return 0, fmt.Errorf("unsupported source currency: %s", from)
	}

	rate, exists := rateMap[to]
	if !exists {
		return 0, fmt.Errorf("conversion from %s to %s not supported", from, to)
	}

	converted := amount * rate
	return math.Round(converted*100) / 100, nil
}

func main() {
	converter := NewExchangeRates()

	amount := 100.0
	from := USD
	to := EUR

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
}