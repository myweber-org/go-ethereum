
package main

import (
	"fmt"
	"os"
	"strconv"
)

const usdToEurRate = 0.85

func convertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run currency_converter.go <amount_in_usd>")
		os.Exit(1)
	}

	usdAmount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Error: Invalid amount '%s'. Please provide a valid number.\n", os.Args[1])
		os.Exit(1)
	}

	if usdAmount < 0 {
		fmt.Println("Error: Amount cannot be negative.")
		os.Exit(1)
	}

	eurAmount := convertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR (Rate: 1 USD = %.2f EUR)\n", usdAmount, eurAmount, usdToEurRate)
}
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

func (er *ExchangeRates) Convert(amount float64, from, to Currency) (float64, error) {
	fromRate, okFrom := er.rates[from]
	toRate, okTo := er.rates[to]

	if !okFrom || !okTo {
		return 0, fmt.Errorf("unsupported currency")
	}

	baseAmount := amount / fromRate
	convertedAmount := baseAmount * toRate

	return math.Round(convertedAmount*100) / 100, nil
}

func (er *ExchangeRates) UpdateRate(currency Currency, rate float64) {
	if rate > 0 {
		er.rates[currency] = rate
	}
}

func main() {
	rates := NewExchangeRates()

	amount := 100.0
	result, err := rates.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)

	rates.UpdateRate(EUR, 0.88)
	newResult, _ := rates.Convert(amount, USD, EUR)
	fmt.Printf("Updated rate: %.2f %s = %.2f %s\n", amount, USD, newResult, EUR)
}