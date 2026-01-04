
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
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run currency_converter.go <amount_in_usd>")
		os.Exit(1)
	}

	usdAmount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Error: Invalid amount '%s'. Please provide a valid number.\n", os.Args[1])
		os.Exit(1)
	}

	if usdAmount < 0 {
		fmt.Println("Error: Amount cannot be negative.")
		os.Exit(1)
	}

	eurAmount := convertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR (Rate: 1 USD = %.2f EUR)\n", usdAmount, eurAmount, usdToEurRate)
}