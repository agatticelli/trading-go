# trading-go

Broker abstraction layer for cryptocurrency perpetual futures trading. Provides a unified interface for multiple exchanges with support for advanced order types, position management, and real-time market data.

## Features

- **Unified Broker Interface**: Single API for multiple exchanges
- **BingX Implementation**: Full support for BingX perpetual futures (demo & live)
- **Advanced Order Types**: Market, Limit, Stop Loss, Take Profit, Trailing Stop
- **Position Management**: Open, close, modify positions with proper risk management
- **Bracket Orders**: Attach TP/SL when opening positions
- **Trailing Stops**: Dynamic stop loss that follows price
- **Demo Mode**: Test strategies safely with demo accounts
- **Type Safety**: Strongly-typed API with clear error handling

## Dependencies

- **[trading-common-types](https://github.com/agatticelli/trading-common-types)**: Shared type definitions (Position, Order, OrderRequest, etc.)

All types are re-exported for convenience, so you can use `broker.Position` or `types.Position` interchangeably - they're the same type!

## Installation

```bash
go get github.com/agatticelli/trading-go
```

## Quick Start

```go
import (
    "context"
    "fmt"

    "github.com/agatticelli/trading-go/broker"
    "github.com/agatticelli/trading-go/bingx"
)

// Create BingX client (demo mode)
client := bingx.NewClient(apiKey, secretKey, true)

// Get account balance
balance, err := client.GetBalance(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("Available: $%.2f\n", balance.Available)

// Get current price
price, err := client.GetCurrentPrice(context.Background(), "BTC-USDT")
fmt.Printf("BTC Price: $%.2f\n", price)

// Place a limit order with TP/SL
order := &broker.OrderRequest{
    Symbol:   "BTC-USDT",
    Side:     broker.SideLong,
    Type:     broker.OrderTypeLimit,
    Size:     0.001,
    Price:    45000.0,
    StopLoss: &broker.StopLossConfig{
        TriggerPrice: 44500.0,
        WorkingType:  broker.WorkingTypeMark,
    },
    TakeProfit: &broker.TakeProfitConfig{
        TriggerPrice: 46000.0,
        WorkingType:  broker.WorkingTypeMark,
    },
}

result, err := client.PlaceOrder(context.Background(), order)
if err != nil {
    panic(err)
}
fmt.Printf("Order placed: %s\n", result.ID)
```

For complete working examples, see the [examples/](examples/) directory.

## Architecture

trading-go is part of a 5-module trading system:

```
trading-go (v0.1.0)     → Broker abstraction (this module)
    ↓
trading-cli             → CLI orchestrator
```

**Key Design Decisions:**

1. **Broker Interface**: Abstract interface that any exchange can implement
2. **Normalized Types**: Exchange-specific data normalized to common types
3. **Zero Dependencies**: Uses only Go standard library
4. **Demo Mode Support**: Safe testing without real money

## Broker Interface

All broker implementations must satisfy this interface:

```go
type Broker interface {
    // Account operations
    GetBalance(ctx context.Context) (*Balance, error)

    // Position operations
    GetPositions(ctx context.Context, filter *PositionFilter) ([]*Position, error)
    GetPosition(ctx context.Context, symbol string) (*Position, error)

    // Order operations
    PlaceOrder(ctx context.Context, order *OrderRequest) (*Order, error)
    GetOrders(ctx context.Context, filter *OrderFilter) ([]*Order, error)
    CancelOrder(ctx context.Context, symbol string, orderID string) error
    CancelAllOrders(ctx context.Context, symbol string) error

    // Market data
    GetCurrentPrice(ctx context.Context, symbol string) (float64, error)

    // Configuration
    SetLeverage(ctx context.Context, symbol string, side string, leverage int) error

    // Metadata
    Name() string
    SupportedFeatures() Features
}
```

## Core Types

### Balance
```go
type Balance struct {
    Asset         string
    Total         float64
    Available     float64
    InUse         float64
    UnrealizedPnL float64
    RealizedPnL   float64
    Timestamp     time.Time
}
```

### Position
```go
type Position struct {
    Symbol            string
    Side              Side     // LONG or SHORT
    Size              float64
    EntryPrice        float64
    MarkPrice         float64
    LiquidationPrice  float64
    Leverage          int
    UnrealizedPnL     float64
    RealizedPnL       float64
    Margin            float64
    Timestamp         time.Time
}
```

### OrderRequest
```go
type OrderRequest struct {
    Symbol      string
    Side        Side        // LONG or SHORT
    Type        OrderType   // MARKET, LIMIT, STOP, etc.
    Size        float64
    Price       float64     // For LIMIT orders
    StopPrice   float64     // For STOP orders
    ReduceOnly  bool        // Only reduce position

    // Advanced features (optional)
    StopLoss   *StopLossConfig
    TakeProfit *TakeProfitConfig
    Trailing   *TrailingConfig
}
```

### Order
```go
type Order struct {
    ID            string
    Symbol        string
    Side          Side
    Type          OrderType
    Status        OrderStatus  // NEW, FILLED, CANCELED, etc.
    Size          float64
    Price         float64
    FilledSize    float64
    AveragePrice  float64
    ReduceOnly    bool
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### Advanced Configs

**Stop Loss:**
```go
type StopLossConfig struct {
    TriggerPrice float64
    OrderPrice   float64     // 0 for market order
    WorkingType  WorkingType // MARK_PRICE or LAST_PRICE
}
```

**Take Profit:**
```go
type TakeProfitConfig struct {
    TriggerPrice float64
    OrderPrice   float64
    WorkingType  WorkingType
}
```

**Trailing Stop:**
```go
type TrailingConfig struct {
    ActivationPrice float64  // Price where trailing starts
    CallbackRate    float64  // Trailing % (0.005 = 0.5%)
}
```

## BingX Implementation

### Features Supported
- ✅ Demo and Production modes
- ✅ Account balance
- ✅ Position management
- ✅ Market, Limit, Stop, Take Profit orders
- ✅ Trailing stop orders
- ✅ Bracket orders (TP/SL on entry)
- ✅ Reduce-only orders
- ✅ Max leverage: 125x
- ✅ HMAC-SHA256 authentication

### Creating a Client

```go
import "github.com/agatticelli/trading-go/bingx"

// Demo mode (safe testing)
demoClient := bingx.NewClient(apiKey, secretKey, true)

// Production mode (real trading)
liveClient := bingx.NewClient(apiKey, secretKey, false)
```

### API Credentials

Get your API keys from:
- **Demo**: https://bingx.com/en-us/demo/
- **Production**: https://bingx.com/en-us/account/api/

⚠️ **Security**: Never commit API keys to source control. Use environment variables.

## Common Operations

### Check Balance
```go
balance, err := client.GetBalance(ctx)
fmt.Printf("Available: $%.2f\n", balance.Available)
fmt.Printf("In Use: $%.2f\n", balance.InUse)
fmt.Printf("Unrealized PnL: $%.2f\n", balance.UnrealizedPnL)
```

### Get Open Positions
```go
positions, err := client.GetPositions(ctx, nil)
for _, pos := range positions {
    fmt.Printf("%s %s: %.4f @ $%.2f (PnL: $%.2f)\n",
        pos.Symbol, pos.Side, pos.Size, pos.EntryPrice, pos.UnrealizedPnL)
}
```

### Place Market Order
```go
order := &broker.OrderRequest{
    Symbol: "BTC-USDT",
    Side:   broker.SideLong,
    Type:   broker.OrderTypeMarket,
    Size:   0.001,
}
result, err := client.PlaceOrder(ctx, order)
```

### Place Limit Order with TP/SL
```go
order := &broker.OrderRequest{
    Symbol: "ETH-USDT",
    Side:   broker.SideLong,
    Type:   broker.OrderTypeLimit,
    Size:   0.1,
    Price:  3000.0,
    StopLoss: &broker.StopLossConfig{
        TriggerPrice: 2900.0,
    },
    TakeProfit: &broker.TakeProfitConfig{
        TriggerPrice: 3200.0,
    },
}
result, err := client.PlaceOrder(ctx, order)
```

### Set Trailing Stop
```go
order := &broker.OrderRequest{
    Symbol:     "BTC-USDT",
    Side:       broker.SideLong,
    Type:       broker.OrderTypeTrailingStop,
    Size:       0.001,
    ReduceOnly: true,
    Trailing: &broker.TrailingConfig{
        ActivationPrice: 46000.0,
        CallbackRate:    0.01, // 1% trailing
    },
}
result, err := client.PlaceOrder(ctx, order)
```

### Cancel Orders
```go
// Cancel specific order
err := client.CancelOrder(ctx, "BTC-USDT", orderID)

// Cancel all orders for symbol
err := client.CancelAllOrders(ctx, "BTC-USDT")
```

## Error Handling

trading-go uses typed errors for common failure cases:

```go
result, err := client.PlaceOrder(ctx, order)
if err != nil {
    switch {
    case errors.Is(err, broker.ErrInsufficientBalance):
        fmt.Println("Not enough balance")
    case errors.Is(err, broker.ErrInvalidSymbol):
        fmt.Println("Invalid trading pair")
    case errors.Is(err, broker.ErrOrderNotFound):
        fmt.Println("Order doesn't exist")
    default:
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Implementing a Custom Broker

To add support for a new exchange:

1. Create a new package (e.g., `binance/`)
2. Implement the `broker.Broker` interface
3. Normalize exchange-specific types to `broker` types

```go
package myexchange

import "github.com/agatticelli/trading-go/broker"

type Client struct {
    apiKey    string
    secretKey string
}

func NewClient(apiKey, secretKey string) *Client {
    return &Client{apiKey: apiKey, secretKey: secretKey}
}

func (c *Client) GetBalance(ctx context.Context) (*broker.Balance, error) {
    // 1. Call exchange API
    // 2. Parse response
    // 3. Normalize to broker.Balance
    return &broker.Balance{ /* ... */ }, nil
}

// Implement remaining interface methods...
```

## Examples

See the [examples/](examples/) directory for complete working code:

- **basic_operations.go**: Balance, positions, orders
- **place_orders.go**: Market, limit, stop orders
- **bracket_orders.go**: Opening positions with TP/SL
- **trailing_stops.go**: Dynamic stop loss management
- **error_handling.go**: Robust error handling patterns

## Dependencies

**None** - Uses only Go standard library:
- `context` - Request cancellation
- `net/http` - HTTP client
- `crypto/hmac` - API authentication
- `encoding/json` - JSON parsing

## Testing

```bash
# Run tests
go test ./...

# Run with demo account
export BINGX_API_KEY="your-demo-key"
export BINGX_SECRET_KEY="your-demo-secret"
go run examples/basic_operations.go
```

## Supported Exchanges

| Exchange | Status | Features |
|----------|--------|----------|
| BingX    | ✅ Complete | All features |
| Binance  | 🚧 Planned | - |
| Bybit    | 🚧 Planned | - |

## License

MIT
