
package main

import (
	"fmt"
	"os"
	"strconv"
)

type ExchangeRate struct {
	Currency string
	Rate     float64
}

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: map[string]float64{
			"USD": 1.0,
			"EUR": 0.85,
			"GBP": 0.73,
			"JPY": 110.0,
			"CAD": 1.25,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	fromRate, ok := c.rates[fromCurrency]
	if !ok {
		return 0, fmt.Errorf("unsupported source currency: %s", fromCurrency)
	}

	toRate, ok := c.rates[toCurrency]
	if !ok {
		return 0, fmt.Errorf("unsupported target currency: %s", toCurrency)
	}

	if fromRate == 0 {
		return 0, fmt.Errorf("invalid exchange rate for currency: %s", fromCurrency)
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
}

func (c *CurrencyConverter) AddRate(currency string, rate float64) {
	c.rates[currency] = rate
}

func (c *CurrencyConverter) ListCurrencies() []string {
	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter()

	if len(os.Args) < 4 {
		fmt.Println("Usage: go run currency_converter.go <amount> <from_currency> <to_currency>")
		fmt.Println("Available currencies:", converter.ListCurrencies())
		return
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		return
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	result, err := converter.Convert(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}