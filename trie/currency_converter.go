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
}