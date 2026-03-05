
package main

import (
	"fmt"
	"os"
	"strconv"
)

const usdToEurRate = 0.85

func convertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run currency_converter.go <amount_in_usd>")
		return
	}

	amountStr := os.Args[1]
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Printf("Error: Invalid amount '%s'\n", amountStr)
		return
	}

	if amount < 0 {
		fmt.Println("Error: Amount cannot be negative")
		return
	}

	converted := convertUSDToEUR(amount)
	fmt.Printf("%.2f USD = %.2f EUR (Rate: 1 USD = %.2f EUR)\n", amount, converted, usdToEurRate)
}package main

import (
	"fmt"
)

const usdToEurRate = 0.92

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	usdAmount := 100.0
	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}
package main

import (
	"fmt"
)

func main() {
	const usdToEurRate = 0.85
	var usdAmount float64

	fmt.Print("Enter amount in USD: ")
	fmt.Scan(&usdAmount)

	eurAmount := usdAmount * usdToEurRate
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ExchangeRates struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
}

func fetchExchangeRates(baseCurrency string) (*ExchangeRates, error) {
	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", baseCurrency)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rates ExchangeRates
	err = json.Unmarshal(body, &rates)
	if err != nil {
		return nil, err
	}

	return &rates, nil
}

func convertCurrency(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	rates, err := fetchExchangeRates(fromCurrency)
	if err != nil {
		return 0, err
	}

	rate, exists := rates.Rates[toCurrency]
	if !exists {
		return 0, fmt.Errorf("currency %s not found in exchange rates", toCurrency)
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

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	result, err := convertCurrency(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}