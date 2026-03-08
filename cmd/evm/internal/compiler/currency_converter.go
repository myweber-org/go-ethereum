package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRateResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	apiEndpoint string
	client      *http.Client
	cache       map[string]ExchangeRateResponse
	lastUpdated time.Time
	cacheTTL    time.Duration
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		apiEndpoint: "https://api.exchangerate.host/latest",
		client:      &http.Client{Timeout: 10 * time.Second},
		cache:       make(map[string]ExchangeRateResponse),
		cacheTTL:    1 * time.Hour,
	}
}

func (c *CurrencyConverter) fetchRates(baseCurrency string) (ExchangeRateResponse, error) {
	if cached, exists := c.cache[baseCurrency]; exists && time.Since(c.lastUpdated) < c.cacheTTL {
		return cached, nil
	}

	url := fmt.Sprintf("%s?base=%s", c.apiEndpoint, baseCurrency)
	resp, err := c.client.Get(url)
	if err != nil {
		return ExchangeRateResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ExchangeRateResponse{}, err
	}

	var rates ExchangeRateResponse
	if err := json.Unmarshal(body, &rates); err != nil {
		return ExchangeRateResponse{}, err
	}

	c.cache[baseCurrency] = rates
	c.lastUpdated = time.Now()
	return rates, nil
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	rates, err := c.fetchRates(fromCurrency)
	if err != nil {
		return 0, err
	}

	rate, exists := rates.Rates[toCurrency]
	if !exists {
		return 0, fmt.Errorf("currency %s not supported", toCurrency)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies(baseCurrency string) ([]string, error) {
	rates, err := c.fetchRates(baseCurrency)
	if err != nil {
		return nil, err
	}

	currencies := make([]string, 0, len(rates.Rates)+1)
	currencies = append(currencies, rates.Base)
	for currency := range rates.Rates {
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	from := "USD"
	to := "EUR"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)

	currencies, err := converter.GetSupportedCurrencies("USD")
	if err != nil {
		fmt.Printf("Error fetching currencies: %v\n", err)
		return
	}

	fmt.Printf("Supported currencies based on USD: %v\n", currencies)
}