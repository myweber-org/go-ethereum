package main

import (
	"fmt"
	"os"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

type CurrencyConverter struct {
	rates map[string]map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]map[string]float64),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	if c.rates[from] == nil {
		c.rates[from] = make(map[string]float64)
	}
	c.rates[from][to] = rate
	
	if c.rates[to] == nil {
		c.rates[to] = make(map[string]float64)
	}
	c.rates[to][from] = 1.0 / rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}
	
	if c.rates[from] == nil {
		return 0, fmt.Errorf("no rates found for currency: %s", from)
	}
	
	rate, exists := c.rates[from][to]
	if !exists {
		return 0, fmt.Errorf("no conversion rate from %s to %s", from, to)
	}
	
	return amount * rate, nil
}

func (c *CurrencyConverter) AvailableCurrencies() []string {
	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("USD", "GBP", 0.73)
	converter.AddRate("EUR", "JPY", 130.0)
	
	amount := 100.0
	fromCurrency := "USD"
	toCurrency := "EUR"
	
	result, err := converter.Convert(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
	
	fmt.Println("Available currencies:")
	for _, currency := range converter.AvailableCurrencies() {
		fmt.Printf("- %s\n", currency)
	}
}