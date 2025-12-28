package main

import (
	"fmt"
)

const usdToEurRate = 0.85

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	amounts := []float64{100.0, 50.0, 25.5}
	
	for _, amount := range amounts {
		converted := ConvertUSDToEUR(amount)
		fmt.Printf("$%.2f USD = €%.2f EUR\n", amount, converted)
	}
}