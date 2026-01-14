package main

import (
	"fmt"
)

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: map[string]float64{
			"USD_EUR": 0.85,
			"EUR_USD": 1.18,
			"USD_GBP": 0.73,
			"GBP_USD": 1.37,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	key := fmt.Sprintf("%s_%s", from, to)
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("conversion rate not available for %s to %s", from, to)
	}
	return amount * rate, nil
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	key := fmt.Sprintf("%s_%s", from, to)
	c.rates[key] = rate
}

func main() {
	converter := NewCurrencyConverter()
	
	usdAmount := 100.0
	eurAmount, err := converter.Convert(usdAmount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
	
	converter.AddRate("EUR_JPY", 130.5)
	jpyAmount, _ := converter.Convert(50.0, "EUR", "JPY")
	fmt.Printf("%.2f EUR = %.2f JPY\n", 50.0, jpyAmount)
}