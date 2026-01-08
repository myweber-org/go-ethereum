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
	rates map[string]float64
	mu    sync.RWMutex
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]float64),
	}
}

func (c *CurrencyConverter) AddRate(base, target string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := base + ":" + target
	c.rates[key] = rate
}

func (c *CurrencyConverter) Convert(amount float64, base, target string) (float64, error) {
	if base == target {
		return amount, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	key := base + ":" + target
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", base, target)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) GetAllRates() []ExchangeRate {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var rates []ExchangeRate
	for key, rate := range c.rates {
		var er ExchangeRate
		fmt.Sscanf(key, "%s:%s", &er.BaseCurrency, &er.TargetCurrency)
		er.Rate = rate
		rates = append(rates, er)
	}
	return rates
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.0)
	
	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	
	allRates := converter.GetAllRates()
	fmt.Println("Available exchange rates:")
	for _, rate := range allRates {
		fmt.Printf("%s to %s: %.4f\n", rate.BaseCurrency, rate.TargetCurrency, rate.Rate)
	}
}