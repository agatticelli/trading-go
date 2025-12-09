# trading-go

Broker abstraction layer for cryptocurrency exchanges.

## Features

- Unified interface for multiple exchanges
- BingX implementation (demo & live)
- Support for perpetual futures trading
- Advanced order types (limit, stop, take profit, trailing)

## Installation

```bash
go get github.com/gattimassimo/trading-go
```

## Usage

```go
import (
    "context"
    "github.com/gattimassimo/trading-go/broker"
    "github.com/gattimassimo/trading-go/bingx"
)

// Create BingX client
client := bingx.NewClient(apiKey, secretKey, false) // false = live mode

// Get balance
balance, err := client.GetBalance(context.Background())

// Get positions
positions, err := client.GetPositions(context.Background(), nil)

// Place order
order := &broker.OrderRequest{
    Symbol: "ETH-USDT",
    Side:   broker.SideLong,
    Type:   broker.OrderTypeLimit,
    Size:   0.5,
    Price:  3950.0,
}
result, err := client.PlaceOrder(context.Background(), order)
```

## Supported Brokers

- [x] BingX (live & demo)
- [ ] Binance (planned)
- [ ] Bybit (planned)

## License

MIT
