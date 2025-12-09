package main

import (
	"context"
	"fmt"
	"os"

	"github.com/agatticelli/trading-go/bingx"
	"github.com/agatticelli/trading-go/broker"
)

// This example demonstrates basic broker operations:
// - Creating a client
// - Checking account balance
// - Getting open positions
// - Getting open orders
// - Checking current market price
func main() {
	fmt.Println("=== Basic Trading Operations Example ===\n")

	// Get API credentials from environment
	apiKey := os.Getenv("BINGX_API_KEY")
	secretKey := os.Getenv("BINGX_SECRET_KEY")

	if apiKey == "" || secretKey == "" {
		fmt.Println("⚠️  API credentials not found in environment")
		fmt.Println("   Set BINGX_API_KEY and BINGX_SECRET_KEY to run this example")
		fmt.Println()
		fmt.Println("Example usage:")
		fmt.Println("  export BINGX_API_KEY=\"your-api-key\"")
		fmt.Println("  export BINGX_SECRET_KEY=\"your-secret-key\"")
		fmt.Println("  go run basic_operations.go")
		fmt.Println()
		fmt.Println("Get demo API keys from: https://bingx.com/en-us/demo/")
		return
	}

	// Create BingX client in demo mode
	fmt.Println("Creating BingX client (demo mode)...")
	client := bingx.NewClient(apiKey, secretKey, true)
	fmt.Printf("✅ Client created: %s\n", client.Name())
	fmt.Printf("   Max Leverage: %dx\n", client.SupportedFeatures().MaxLeverage)
	fmt.Printf("   Trailing Stops: %v\n", client.SupportedFeatures().TrailingStop)
	fmt.Println()

	ctx := context.Background()

	// 1. Get account balance
	fmt.Println("📊 Fetching account balance...")
	balance, err := client.GetBalance(ctx)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("✅ Account Balance:")
	fmt.Printf("   Asset: %s\n", balance.Asset)
	fmt.Printf("   Total: $%.2f\n", balance.Total)
	fmt.Printf("   Available: $%.2f\n", balance.Available)
	fmt.Printf("   In Use: $%.2f\n", balance.InUse)
	fmt.Printf("   Unrealized PnL: $%.2f\n", balance.UnrealizedPnL)
	fmt.Printf("   Realized PnL: $%.2f\n", balance.RealizedPnL)
	fmt.Println()

	// 2. Get current market price
	symbol := "BTC-USDT"
	fmt.Printf("📈 Fetching current price for %s...\n", symbol)
	price, err := client.GetCurrentPrice(ctx, symbol)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Current Price: $%.2f\n", price)
	}
	fmt.Println()

	// 3. Get open positions
	fmt.Println("📋 Fetching open positions...")
	positions, err := client.GetPositions(ctx, nil)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(positions) == 0 {
		fmt.Println("✅ No open positions")
	} else {
		fmt.Printf("✅ Found %d open position(s):\n", len(positions))
		for i, pos := range positions {
			fmt.Printf("\n   Position #%d:\n", i+1)
			fmt.Printf("     Symbol: %s\n", pos.Symbol)
			fmt.Printf("     Side: %s\n", pos.Side)
			fmt.Printf("     Size: %.4f\n", pos.Size)
			fmt.Printf("     Entry Price: $%.2f\n", pos.EntryPrice)
			fmt.Printf("     Mark Price: $%.2f\n", pos.MarkPrice)
			fmt.Printf("     Leverage: %dx\n", pos.Leverage)
			fmt.Printf("     Unrealized PnL: $%.2f\n", pos.UnrealizedPnL)
			fmt.Printf("     Liquidation Price: $%.2f\n", pos.LiquidationPrice)
		}
	}
	fmt.Println()

	// 4. Get open orders
	fmt.Println("📝 Fetching open orders...")
	orders, err := client.GetOrders(ctx, &broker.OrderFilter{})
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(orders) == 0 {
		fmt.Println("✅ No open orders")
	} else {
		fmt.Printf("✅ Found %d open order(s):\n", len(orders))
		for i, order := range orders {
			fmt.Printf("\n   Order #%d:\n", i+1)
			fmt.Printf("     ID: %s\n", order.ID)
			fmt.Printf("     Symbol: %s\n", order.Symbol)
			fmt.Printf("     Side: %s\n", order.Side)
			fmt.Printf("     Type: %s\n", order.Type)
			fmt.Printf("     Status: %s\n", order.Status)
			fmt.Printf("     Size: %.4f\n", order.Size)
			fmt.Printf("     Price: $%.2f\n", order.Price)
			if order.StopPrice > 0 {
				fmt.Printf("     Stop Price: $%.2f\n", order.StopPrice)
			}
			fmt.Printf("     Reduce Only: %v\n", order.ReduceOnly)
		}
	}

	fmt.Println("\n💡 This example shows basic read operations")
	fmt.Println("   See other examples for placing/canceling orders")
}
