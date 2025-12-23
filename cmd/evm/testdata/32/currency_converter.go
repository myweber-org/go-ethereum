
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
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	converter := &CurrencyConverter{
		rates: make(map[string]float64),
	}
	
	converter.AddRate(USD, EUR, 0.92)
	converter.AddRate(USD, GBP, 0.79)
	converter.AddRate(USD, JPY, 148.50)
	converter.AddRate(EUR, USD, 1.09)
	converter.AddRate(EUR, GBP, 0.86)
	converter.AddRate(EUR, JPY, 161.41)
	converter.AddRate(GBP, USD, 1.27)
	converter.AddRate(GBP, EUR, 1.16)
	converter.AddRate(GBP, JPY, 187.97)
	converter.AddRate(JPY, USD, 0.0067)
	converter.AddRate(JPY, EUR, 0.0062)
	converter.AddRate(JPY, GBP, 0.0053)
	
	return converter
}

func (c *CurrencyConverter) AddRate(from, to Currency, rate float64) {
	key := fmt.Sprintf("%s:%s", from, to)
	c.rates[key] = rate
}

func (c *CurrencyConverter) GetRate(from, to Currency) (float64, error) {
	if from == to {
		return 1.0, nil
	}
	
	key := fmt.Sprintf("%s:%s", from, to)
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}
	
	return rate, nil
}

func (c *CurrencyConverter) Convert(amount float64, from, to Currency) (float64, error) {
	rate, err := c.GetRate(from, to)
	if err != nil {
		return 0, err
	}
	
	converted := amount * rate
	return math.Round(converted*100) / 100, nil
}

func (c *CurrencyConverter) ConvertWithPrecision(amount float64, from, to Currency, precision int) (float64, error) {
	rate, err := c.GetRate(from, to)
	if err != nil {
		return 0, err
	}
	
	converted := amount * rate
	multiplier := math.Pow(10, float64(precision))
	return math.Round(converted*multiplier) / multiplier, nil
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
	
	converted, err = converter.ConvertWithPrecision(amount, GBP, USD, 4)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.4f %s\n", amount, GBP, converted, USD)
}