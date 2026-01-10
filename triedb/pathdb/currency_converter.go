package main

import (
	"fmt"
	"sync"
)

type ExchangeRate struct {
	BaseCurrency    string
	TargetCurrency  string
	Rate            float64
}

type CurrencyConverter struct {
	rates map[string]map[string]float64
	mu    sync.RWMutex
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]map[string]float64),
	}
}

func (c *CurrencyConverter) AddRate(base, target string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.rates[base]; !exists {
		c.rates[base] = make(map[string]float64)
	}
	c.rates[base][target] = rate

	if _, exists := c.rates[target]; !exists {
		c.rates[target] = make(map[string]float64)
	}
	c.rates[target][base] = 1.0 / rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if from == to {
		return amount, nil
	}

	if targetRates, exists := c.rates[from]; exists {
		if rate, ok := targetRates[to]; ok {
			return amount * rate, nil
		}
	}

	return 0, fmt.Errorf("conversion rate not available from %s to %s", from, to)
}

func (c *CurrencyConverter) GetAllRates(base string) map[string]float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	rates := make(map[string]float64)
	if baseRates, exists := c.rates[base]; exists {
		for currency, rate := range baseRates {
			rates[currency] = rate
		}
	}
	return rates
}

func main() {
	converter := NewCurrencyConverter()

	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("USD", "GBP", 0.73)
	converter.AddRate("EUR", "JPY", 130.0)

	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	}

	usdRates := converter.GetAllRates("USD")
	fmt.Println("USD exchange rates:")
	for currency, rate := range usdRates {
		fmt.Printf("  %s: %.4f\n", currency, rate)
	}
}