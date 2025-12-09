package main

import (
	"context"
	"fmt"

	"github.com/agatticelli/trading-go/broker"
	"github.com/agatticelli/trading-go/bingx"
)

// This example demonstrates different order types
// NOTE: This is example code. Do not run without understanding the implications!
func main() {
	fmt.Println("=== Order Placement Examples ===\n")
	fmt.Println("⚠️  This is example code showing the API")
	fmt.Println("   Uncomment sections to actually place orders\n")

	// Example client setup
	apiKey := "your-api-key"
	secretKey := "your-secret-key"
	client := bingx.NewClient(apiKey, secretKey, true) // Demo mode
	ctx := context.Background()

	// Example 1: Market Order
	fmt.Println("Example 1: Market Order (immediate execution)")
	marketOrder := &broker.OrderRequest{
		Symbol: "BTC-USDT",
		Side:   broker.SideLong,
		Type:   broker.OrderTypeMarket,
		Size:   0.001,
	}
	fmt.Printf("  %+v\n\n", marketOrder)

	// Uncomment to actually place:
	// result, err := client.PlaceOrder(ctx, marketOrder)

	// Example 2: Limit Order
	fmt.Println("Example 2: Limit Order (wait for price)")
	limitOrder := &broker.OrderRequest{
		Symbol:      "BTC-USDT",
		Side:        broker.SideLong,
		Type:        broker.OrderTypeLimit,
		Size:        0.001,
		Price:       45000.0, // Buy at $45k
		TimeInForce: broker.TimeInForceGTC,
	}
	fmt.Printf("  %+v\n\n", limitOrder)

	// Example 3: Stop Loss Order
	fmt.Println("Example 3: Stop Loss Order (triggers below entry)")
	stopOrder := &broker.OrderRequest{
		Symbol:     "BTC-USDT",
		Side:       broker.SideShort, // Close LONG = SHORT
		Type:       broker.OrderTypeStop,
		Size:       0.001,
		StopPrice:  44000.0, // Trigger at $44k
		ReduceOnly: true,    // Only close position
	}
	fmt.Printf("  %+v\n\n", stopOrder)

	// Example 4: Take Profit Order
	fmt.Println("Example 4: Take Profit Order (triggers above entry)")
	takeProfitOrder := &broker.OrderRequest{
		Symbol:     "BTC-USDT",
		Side:       broker.SideShort,
		Type:       broker.OrderTypeTakeProfit,
		Size:       0.001,
		StopPrice:  46000.0,
		ReduceOnly: true,
	}
	fmt.Printf("  %+v\n\n", takeProfitOrder)

	// Example 5: Limit Order with Stop Loss and Take Profit
	fmt.Println("Example 5: Bracket Order (entry + TP + SL)")
	bracketOrder := &broker.OrderRequest{
		Symbol:      "BTC-USDT",
		Side:        broker.SideLong,
		Type:        broker.OrderTypeLimit,
		Size:        0.001,
		Price:       45000.0,
		TimeInForce: broker.TimeInForceGTC,
		StopLoss: &broker.StopLossConfig{
			TriggerPrice: 44500.0,
			WorkingType:  broker.WorkingTypeMark,
		},
		TakeProfit: &broker.TakeProfitConfig{
			TriggerPrice: 46000.0,
			WorkingType:  broker.WorkingTypeMark,
		},
	}
	fmt.Printf("  Entry: $%.2f\n", bracketOrder.Price)
	fmt.Printf("  Stop Loss: $%.2f\n", bracketOrder.StopLoss.TriggerPrice)
	fmt.Printf("  Take Profit: $%.2f\n\n", bracketOrder.TakeProfit.TriggerPrice)

	// Example 6: Trailing Stop Order
	fmt.Println("Example 6: Trailing Stop (dynamic stop loss)")
	trailingOrder := &broker.OrderRequest{
		Symbol:     "BTC-USDT",
		Side:       broker.SideShort,
		Type:       broker.OrderTypeTrailingStop,
		Size:       0.001,
		ReduceOnly: true,
		Trailing: &broker.TrailingConfig{
			ActivationPrice: 46000.0, // Start trailing at $46k
			CallbackRate:    0.01,    // Trail by 1%
		},
	}
	fmt.Printf("  Activation: $%.2f\n", trailingOrder.Trailing.ActivationPrice)
	fmt.Printf("  Callback Rate: %.1f%%\n", trailingOrder.Trailing.CallbackRate*100)
	fmt.Println()

	fmt.Println("💡 Key Points:")
	fmt.Println("   - Market orders execute immediately at current price")
	fmt.Println("   - Limit orders wait for your specified price")
	fmt.Println("   - Stop orders trigger when price crosses threshold")
	fmt.Println("   - Bracket orders place entry + TP + SL atomically")
	fmt.Println("   - Trailing stops follow price to lock in profits")
	fmt.Println()
	fmt.Println("⚠️  Always use demo mode for testing!")

	_ = ctx    // Suppress unused warning
	_ = client // Suppress unused warning
}
