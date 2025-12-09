package main

import (
	"context"
	"fmt"
	"os"

	"github.com/agatticelli/trading-go/bingx"
	"github.com/agatticelli/trading-go/broker"
)

func main() {
	// Get credentials from environment or use test values
	apiKey := os.Getenv("BINGX_API_KEY")
	secretKey := os.Getenv("BINGX_SECRET_KEY")
	
	if apiKey == "" || secretKey == "" {
		fmt.Println("Please set BINGX_API_KEY and BINGX_SECRET_KEY")
		return
	}

	client := bingx.NewClient(apiKey, secretKey, true) // true = demo mode

	fmt.Println("Testing GetOrders...")
	orders, err := client.GetOrders(context.Background(), &broker.OrderFilter{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Success! Found %d orders\n", len(orders))
	for _, order := range orders {
		fmt.Printf("  - %s: %s %s %.4f @ %.2f\n", 
			order.ID, order.Symbol, order.Side, order.Size, order.Price)
	}
}
