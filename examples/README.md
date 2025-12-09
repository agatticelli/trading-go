# Trading-Go Examples

This directory contains examples demonstrating how to use trading-go to interact with cryptocurrency exchanges.

## Prerequisites

### Get API Credentials

**For Testing (Recommended):**
1. Visit https://bingx.com/en-us/demo/
2. Create a demo account
3. Generate API keys (Settings → API Management)
4. Export as environment variables:
   ```bash
   export BINGX_API_KEY="your-demo-api-key"
   export BINGX_SECRET_KEY="your-demo-secret-key"
   ```

**For Production:**
- Get keys from: https://bingx.com/en-us/account/api/
- ⚠️ **Never commit API keys to source control**
- Use demo mode for testing first!

## Running the Examples

### 1. Basic Operations
View account balance, positions, and orders.

```bash
go run basic_operations.go
```

**What it demonstrates:**
- Creating a BingX client
- Checking account balance
- Getting open positions with PnL
- Listing open orders
- Fetching current market prices

**Sample Output:**
```
✅ Account Balance:
   Asset: USDT
   Available: $1000.00
   In Use: $0.00
   Unrealized PnL: $0.00

✅ Current Price: $45234.56

✅ No open positions
✅ No open orders
```

### 2. Place Orders
Examples of different order types (code examples, not live execution).

```bash
go run place_orders.go
```

**What it demonstrates:**
- Market orders (immediate execution)
- Limit orders (wait for price)
- Stop loss orders
- Take profit orders
- Bracket orders (entry + TP + SL)
- Trailing stop orders

**Note:** This example shows API usage but doesn't execute orders by default. Uncomment specific sections to actually place orders.

### 3. Error Handling
Robust error handling patterns for production code.

```bash
go run error_handling.go
```

**What it demonstrates:**
- Typed error handling with `errors.Is()`
- Context cancellation
- Network error handling
- Retry logic for transient failures
- Different error types:
  - `ErrInsufficientBalance`
  - `ErrInvalidSymbol`
  - `ErrInvalidPrice`
  - `ErrRateLimited`
  - `ErrNetworkError`

## Common Patterns

### Creating a Client

```go
import (
    "github.com/agatticelli/trading-go/bingx"
)

// Demo mode (recommended for testing)
client := bingx.NewClient(apiKey, secretKey, true)

// Production mode
client := bingx.NewClient(apiKey, secretKey, false)
```

### Placing a Simple Order

```go
order := &broker.OrderRequest{
    Symbol: "BTC-USDT",
    Side:   broker.SideLong,
    Type:   broker.OrderTypeMarket,
    Size:   0.001,
}

result, err := client.PlaceOrder(ctx, order)
if err != nil {
    log.Fatal(err)
}
```

### Placing Order with TP/SL

```go
order := &broker.OrderRequest{
    Symbol: "BTC-USDT",
    Side:   broker.SideLong,
    Type:   broker.OrderTypeLimit,
    Size:   0.001,
    Price:  45000.0,
    StopLoss: &broker.StopLossConfig{
        TriggerPrice: 44500.0,
    },
    TakeProfit: &broker.TakeProfitConfig{
        TriggerPrice: 46000.0,
    },
}
```

### Setting a Trailing Stop

```go
order := &broker.OrderRequest{
    Symbol:     "BTC-USDT",
    Side:       broker.SideShort,  // Close a LONG position
    Type:       broker.OrderTypeTrailingStop,
    Size:       0.001,
    ReduceOnly: true,
    Trailing: &broker.TrailingConfig{
        ActivationPrice: 46000.0,  // Start trailing at $46k
        CallbackRate:    0.01,     // Trail by 1%
    },
}
```

### Error Handling

```go
result, err := client.PlaceOrder(ctx, order)
if err != nil {
    switch {
    case errors.Is(err, broker.ErrInsufficientBalance):
        fmt.Println("Not enough balance")
    case errors.Is(err, broker.ErrInvalidSymbol):
        fmt.Println("Invalid trading pair")
    default:
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Safety Tips

### For Testing
- ✅ **Always start with demo mode**
- ✅ Use small position sizes
- ✅ Test thoroughly before live trading
- ✅ Understand each order type

### For Production
- ⚠️ **Never commit API keys**
- ⚠️ Use environment variables for credentials
- ⚠️ Implement proper error handling
- ⚠️ Add retry logic for network errors
- ⚠️ Monitor rate limits
- ⚠️ Log all operations
- ⚠️ Test with small amounts first

## Order Types Quick Reference

| Type | Description | When to Use |
|------|-------------|-------------|
| **Market** | Execute immediately at current price | Quick entry/exit |
| **Limit** | Execute at specific price or better | Wait for better price |
| **Stop** | Trigger when price crosses threshold | Stop loss protection |
| **Take Profit** | Close position at profit target | Lock in profits |
| **Trailing Stop** | Dynamic stop that follows price | Maximize profits |

## Supported Symbols

Common perpetual futures symbols on BingX:
- `BTC-USDT`
- `ETH-USDT`
- `BNB-USDT`
- `SOL-USDT`
- `XRP-USDT`
- And many more...

Check BingX documentation for complete list.

## Troubleshooting

### "API credentials not found"
- Set environment variables: `BINGX_API_KEY` and `BINGX_SECRET_KEY`
- Verify keys are correct (no extra spaces)

### "Invalid symbol"
- Use format: `BTC-USDT` (not `BTCUSDT` or `BTC/USDT`)
- Check symbol exists on BingX

### "Insufficient balance"
- Check available balance with `GetBalance()`
- Reduce position size or add funds

### "Rate limited"
- Slow down requests
- Implement exponential backoff
- Check BingX rate limits

## Further Reading

- [Main README](../README.md) - Complete API documentation
- [BingX API Docs](https://bingx-api.github.io/docs/) - Official documentation
- [MIGRATION_STATUS.md](../../trading-cli/MIGRATION_STATUS.md) - Architecture overview
