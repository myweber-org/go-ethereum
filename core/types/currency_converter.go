
package main

import (
	"fmt"
)

const usdToEurRate = 0.92

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	amounts := []float64{100.0, 250.0, 50.0}
	
	fmt.Println("USD to EUR Conversion")
	fmt.Println("=====================")
	
	for _, usd := range amounts {
		eur := ConvertUSDToEUR(usd)
		fmt.Printf("$%.2f USD = €%.2f EUR\n", usd, eur)
	}
}