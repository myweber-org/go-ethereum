package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	usdToEurRate = 0.92
	usdToGbpRate = 0.79
)

func convertUSD(amount float64, targetCurrency string) (float64, error) {
	switch targetCurrency {
	case "EUR":
		return amount * usdToEurRate, nil
	case "GBP":
		return amount * usdToGbpRate, nil
	default:
		return 0, fmt.Errorf("unsupported currency: %s", targetCurrency)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run currency_converter.go <amount> <currency>")
		fmt.Println("Example: go run currency_converter.go 100 EUR")
		return
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		return
	}

	targetCurrency := os.Args[2]
	converted, err := convertUSD(amount, targetCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f USD = %.2f %s\n", amount, converted, targetCurrency)
}