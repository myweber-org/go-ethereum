package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	apiURL    string
	rates     map[string]float64
	lastFetch time.Time
	cacheTTL  time.Duration
}

func NewCurrencyConverter(apiKey string) *CurrencyConverter {
	return &CurrencyConverter{
		apiURL:    fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/USD"),
		rates:     make(map[string]float64),
		cacheTTL:  30 * time.Minute,
		lastFetch: time.Time{},
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastFetch) < c.cacheTTL && len(c.rates) > 0 {
		return nil
	}

	resp, err := http.Get(c.apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var data ExchangeRateResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	c.rates = data.Rates
	c.lastFetch = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if err := c.fetchRates(); err != nil {
		return 0, err
	}

	fromRate, fromExists := c.rates[fromCurrency]
	toRate, toExists := c.rates[toCurrency]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unsupported currency: %s or %s", fromCurrency, toCurrency)
	}

	usdAmount := amount / fromRate
	convertedAmount := usdAmount * toRate

	return convertedAmount, nil
}

func main() {
	converter := NewCurrencyConverter("")

	amount := 100.0
	from := "EUR"
	to := "JPY"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
}package main

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
	key := from + "_" + to
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

	key := from + "_" + to
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}

	return amount * rate.Rate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	currencies := make(map[string]bool)
	for _, rate := range c.rates {
		currencies[rate.FromCurrency] = true
		currencies[rate.ToCurrency] = true
	}

	result := make([]string, 0, len(currencies))
	for currency := range currencies {
		result = append(result, currency)
	}
	return result
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.5)
	
	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	fmt.Printf("Supported currencies: %v\n", converter.GetSupportedCurrencies())
}