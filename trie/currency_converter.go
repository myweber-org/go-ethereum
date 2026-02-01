package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strconv"
)

type ExchangeRates struct {
    Rates map[string]float64 `json:"rates"`
    Base  string             `json:"base"`
    Date  string             `json:"date"`
}

func fetchExchangeRates(baseCurrency string) (*ExchangeRates, error) {
    url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", baseCurrency)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var rates ExchangeRates
    if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
        return nil, err
    }
    return &rates, nil
}

func convertCurrency(amount float64, from, to string) (float64, error) {
    rates, err := fetchExchangeRates(from)
    if err != nil {
        return 0, err
    }

    rate, exists := rates.Rates[to]
    if !exists {
        return 0, fmt.Errorf("currency %s not found", to)
    }

    return amount * rate, nil
}

func main() {
    if len(os.Args) != 4 {
        fmt.Println("Usage: currency_converter <amount> <from_currency> <to_currency>")
        fmt.Println("Example: currency_converter 100 USD EUR")
        os.Exit(1)
    }

    amount, err := strconv.ParseFloat(os.Args[1], 64)
    if err != nil {
        fmt.Printf("Invalid amount: %v\n", err)
        os.Exit(1)
    }

    from := os.Args[2]
    to := os.Args[3]

    result, err := convertCurrency(amount, from, to)
    if err != nil {
        fmt.Printf("Conversion error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
}package main

import (
	"fmt"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

type CurrencyConverter struct {
	rates []ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: []ExchangeRate{
			{"USD", "EUR", 0.92},
			{"EUR", "USD", 1.09},
			{"USD", "JPY", 149.50},
			{"JPY", "USD", 0.0067},
			{"GBP", "USD", 1.27},
			{"USD", "GBP", 0.79},
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	for _, rate := range c.rates {
		if rate.FromCurrency == fromCurrency && rate.ToCurrency == toCurrency {
			return amount * rate.Rate, nil
		}
	}

	return 0, fmt.Errorf("conversion rate not found for %s to %s", fromCurrency, toCurrency)
}

func (c *CurrencyConverter) AddRate(fromCurrency, toCurrency string, rate float64) {
	c.rates = append(c.rates, ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
	})
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)

	converter.AddRate("EUR", "JPY", 162.50)
	result2, err := converter.Convert(50.0, "EUR", "JPY")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f EUR = %.2f JPY\n", 50.0, result2)
}