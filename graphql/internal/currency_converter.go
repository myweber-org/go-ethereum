
package main

import (
	"fmt"
	"os"
	"strconv"
)

const usdToEurRate = 0.85

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run currency_converter.go <amount_in_usd>")
		os.Exit(1)
	}

	usdAmount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		os.Exit(1)
	}

	if usdAmount < 0 {
		fmt.Println("Amount cannot be negative")
		os.Exit(1)
	}

	eurAmount := usdAmount * usdToEurRate
	fmt.Printf("%.2f USD = %.2f EUR (Rate: 1 USD = %.2f EUR)\n", usdAmount, eurAmount, usdToEurRate)
}