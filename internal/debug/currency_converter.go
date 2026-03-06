package main

import (
	"fmt"
)

const usdToEurRate = 0.85

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	amounts := []float64{100.0, 250.5, 75.25}
	
	for _, usd := range amounts {
		eur := ConvertUSDToEUR(usd)
		fmt.Printf("USD %.2f = EUR %.2f\n", usd, eur)
	}
}