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

	if c.rates[base] == nil {
		c.rates[base] = make(map[string]float64)
	}
	c.rates[base][target] = rate

	// Add inverse rate
	if c.rates[target] == nil {
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

	if rates, ok := c.rates[from]; ok {
		if rate, ok := rates[to]; ok {
			return amount * rate, nil
		}
	}

	return 0, fmt.Errorf("no conversion rate found from %s to %s", from, to)
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter()

	// Add some sample exchange rates
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("USD", "JPY", 110.0)
	converter.AddRate("EUR", "GBP", 0.86)

	// Perform conversions
	amounts := []float64{100.0, 250.0, 500.0}
	conversions := []struct {
		from string
		to   string
	}{
		{"USD", "EUR"},
		{"EUR", "GBP"},
		{"JPY", "USD"},
	}

	for i, amount := range amounts {
		conv := conversions[i%len(conversions)]
		result, err := converter.Convert(amount, conv.from, conv.to)
		if err != nil {
			fmt.Printf("Error converting %.2f %s to %s: %v\n", amount, conv.from, conv.to, err)
		} else {
			fmt.Printf("%.2f %s = %.2f %s\n", amount, conv.from, result, conv.to)
		}
	}

	// Show supported currencies
	fmt.Printf("\nSupported currencies: %v\n", converter.GetSupportedCurrencies())
}