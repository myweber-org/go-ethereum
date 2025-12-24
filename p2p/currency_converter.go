
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

type ExchangeRate struct {
	From Currency
	To   Currency
	Rate float64
}

type CurrencyConverter struct {
	rates []ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: []ExchangeRate{
			{USD, EUR, 0.92},
			{USD, GBP, 0.79},
			{USD, JPY, 149.5},
			{EUR, USD, 1.09},
			{EUR, GBP, 0.86},
			{EUR, JPY, 162.5},
			{GBP, USD, 1.27},
			{GBP, EUR, 1.16},
			{GBP, JPY, 189.2},
			{JPY, USD, 0.0067},
			{JPY, EUR, 0.0062},
			{JPY, GBP, 0.0053},
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from Currency, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	for _, rate := range c.rates {
		if rate.From == from && rate.To == to {
			return math.Round(amount*rate.Rate*100) / 100, nil
		}
	}

	return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
}

func (c *CurrencyConverter) AddRate(from Currency, to Currency, rate float64) {
	c.rates = append(c.rates, ExchangeRate{from, to, rate})
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	converted, err := converter.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, converted, EUR)

	converted, err = converter.Convert(amount, EUR, JPY)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, EUR, converted, JPY)

	converter.AddRate(USD, CAD, 1.35)
	converted, err = converter.Convert(amount, USD, CAD)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, converted, CAD)
}
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
	rates map[string]ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]ExchangeRate),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	key := fmt.Sprintf("%s:%s", from, to)
	c.rates[key] = ExchangeRate{
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

	key := fmt.Sprintf("%s:%s", from, to)
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}

	return amount * rate.Rate, nil
}

func (c *CurrencyConverter) GetRateCount() int {
	return len(c.rates)
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.5)
	converter.AddRate("GBP", "USD", 1.38)
	
	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	fmt.Printf("Total exchange rates stored: %d\n", converter.GetRateCount())
}