package main

import (
	"fmt"
	"math"
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
	if _, exists := c.rates[from]; !exists {
		c.rates[from] = make(map[string]float64)
	}
	c.rates[from][to] = rate
	
	if _, exists := c.rates[to]; !exists {
		c.rates[to] = make(map[string]float64)
	}
	c.rates[to][from] = 1.0 / rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}
	
	if rates, exists := c.rates[from]; exists {
		if rate, exists := rates[to]; exists {
			return math.Round(amount*rate*100) / 100, nil
		}
	}
	
	return 0, fmt.Errorf("no conversion rate found from %s to %s", from, to)
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("USD", "JPY", 110.0)
	converter.AddRate("EUR", "GBP", 0.86)
	
	amount := 100.0
	
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	
	result, err = converter.Convert(amount, "USD", "JPY")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f JPY\n", amount, result)
	
	result, err = converter.Convert(amount, "EUR", "GBP")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f EUR = %.2f GBP\n", amount, result)
	
	fmt.Println("Supported currencies:", converter.GetSupportedCurrencies())
}