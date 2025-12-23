
package main

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	LastUpdated  time.Time
}

type CurrencyConverter struct {
	rates map[string]map[string]ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]map[string]ExchangeRate),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	if c.rates[from] == nil {
		c.rates[from] = make(map[string]ExchangeRate)
	}
	c.rates[from][to] = ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		LastUpdated:  time.Now(),
	}
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
		return 0, fmt.Errorf("conversion rate not available from %s to %s", from, to)
	}

	return amount * rate.Rate, nil
}

func (c *CurrencyConverter) GetRate(from, to string) (ExchangeRate, error) {
	if c.rates[from] == nil {
		return ExchangeRate{}, fmt.Errorf("no rates found for currency: %s", from)
	}

	rate, exists := c.rates[from][to]
	if !exists {
		return ExchangeRate{}, fmt.Errorf("rate not found from %s to %s", from, to)
	}

	return rate, nil
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.5)
	
	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	
	rate, err := converter.GetRate("USD", "JPY")
	if err != nil {
		fmt.Printf("Rate retrieval error: %v\n", err)
		return
	}
	
	fmt.Printf("Exchange rate from %s to %s: %.4f (updated: %v)\n",
		rate.FromCurrency, rate.ToCurrency, rate.Rate, rate.LastUpdated.Format(time.RFC3339))
}