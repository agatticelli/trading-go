package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/agatticelli/trading-go/broker"
	"github.com/agatticelli/trading-go/bingx"
)

// This example demonstrates robust error handling
func main() {
	fmt.Println("=== Error Handling Example ===\n")

	client := bingx.NewClient("demo-key", "demo-secret", true)
	ctx := context.Background()

	// Example 1: Handling typed errors
	fmt.Println("Example 1: Typed Error Handling\n")

	order := &broker.OrderRequest{
		Symbol: "BTC-USDT",
		Side:   broker.SideLong,
		Type:   broker.OrderTypeLimit,
		Size:   100.0, // Intentionally large to trigger insufficient balance
		Price:  45000.0,
	}

	result, err := client.PlaceOrder(ctx, order)
	if err != nil {
		handleOrderError(err)
	} else {
		fmt.Printf("✅ Order placed: %s\n", result.ID)
	}
	fmt.Println()

	// Example 2: Context cancellation
	fmt.Println("Example 2: Context Cancellation\n")

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err = client.GetBalance(cancelCtx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Println("❌ Operation canceled by user")
		} else {
			fmt.Printf("❌ Error: %v\n", err)
		}
	}
	fmt.Println()

	// Example 3: API errors
	fmt.Println("Example 3: API Error Handling\n")

	_, err = client.GetCurrentPrice(ctx, "INVALID-SYMBOL")
	if err != nil {
		switch {
		case errors.Is(err, broker.ErrInvalidSymbol):
			fmt.Println("❌ Invalid trading symbol")
			fmt.Println("   Use format: BTC-USDT, ETH-USDT, etc.")
		case errors.Is(err, broker.ErrAPIError):
			fmt.Println("❌ API error occurred")
			fmt.Println("   Check exchange status and try again")
		default:
			fmt.Printf("❌ Unexpected error: %v\n", err)
		}
	}
	fmt.Println()

	// Example 4: Retry logic
	fmt.Println("Example 4: Retry Logic Pattern\n")
	demonstrateRetryPattern(client)

	fmt.Println("\n💡 Best Practices:")
	fmt.Println("   - Always check errors before using results")
	fmt.Println("   - Use typed errors (errors.Is) for specific handling")
	fmt.Println("   - Implement retry logic for transient failures")
	fmt.Println("   - Use context timeouts for long operations")
	fmt.Println("   - Log errors with sufficient context")
}

func handleOrderError(err error) {
	switch {
	case errors.Is(err, broker.ErrInsufficientBalance):
		fmt.Println("❌ Insufficient balance")
		fmt.Println("   Deposit more funds or reduce position size")

	case errors.Is(err, broker.ErrInvalidPrice):
		fmt.Println("❌ Invalid price")
		fmt.Println("   Check price is within allowed range")

	case errors.Is(err, broker.ErrInvalidQuantity):
		fmt.Println("❌ Invalid order size")
		fmt.Println("   Check minimum/maximum order size for symbol")

	case errors.Is(err, broker.ErrInvalidSymbol):
		fmt.Println("❌ Invalid symbol")
		fmt.Println("   Ensure symbol is in correct format (e.g., BTC-USDT)")

	case errors.Is(err, broker.ErrOrderNotFound):
		fmt.Println("❌ Order not found")
		fmt.Println("   Order may have been filled or canceled")

	case errors.Is(err, broker.ErrRateLimited):
		fmt.Println("❌ Rate limited")
		fmt.Println("   Too many requests - slow down and retry")

	case errors.Is(err, broker.ErrAPIError):
		fmt.Println("❌ API error")
		fmt.Println("   Check exchange status and retry")

	default:
		fmt.Printf("❌ Unexpected error: %v\n", err)
		fmt.Println("   Contact support if error persists")
	}
}

func demonstrateRetryPattern(client *bingx.Client) {
	maxRetries := 3
	var balance *broker.Balance
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("Attempt %d/%d... ", attempt, maxRetries)

		balance, err = client.GetBalance(context.Background())
		if err == nil {
			fmt.Println("✅ Success")
			fmt.Printf("Balance: $%.2f\n", balance.Available)
			return
		}

		if !isRetryable(err) {
			fmt.Println("❌ Non-retryable error")
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("⚠️  Retryable error, trying again...")
	}

	fmt.Println("❌ Max retries exceeded")
}

func isRetryable(err error) bool {
	return errors.Is(err, broker.ErrAPIError) ||
		errors.Is(err, broker.ErrRateLimited) ||
		errors.Is(err, context.DeadlineExceeded)
}
