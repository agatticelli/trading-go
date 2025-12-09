# Adding a New Broker to trading-go

This guide walks you through implementing support for a new cryptocurrency exchange.

## Overview

Adding a new broker involves:
1. Creating a new package for the broker
2. Implementing the `broker.Broker` interface
3. Normalizing exchange-specific types to common types
4. Handling authentication and API requests
5. Error handling and edge cases
6. Testing your implementation

**Estimated time**: 4-8 hours for a basic implementation

---

## Step 1: Create Broker Package

Create a new directory for your broker:

```bash
mkdir -p /path/to/trading-go/youexchange
cd youexchange
```

Create the following files:
- `client.go` - Main broker implementation
- `types.go` - Exchange-specific types
- `auth.go` - Authentication/signing logic
- `positions.go` - Position-related operations
- `orders.go` - Order-related operations
- `balance.go` - Balance operations
- `market.go` - Market data operations

---

## Step 2: Understand the Broker Interface

Your broker must implement this interface:

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

---

## Step 3: Implement the Client

### 3.1: Basic Client Structure

**File: `youexchange/client.go`**

```go
package youexchange

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strconv"
    "time"

    "github.com/agatticelli/trading-go/broker"
    "github.com/agatticelli/trading-common-types"
)

// Client implements the broker.Broker interface for YourExchange
type Client struct {
    apiKey    string
    secretKey string
    baseURL   string
    client    *http.Client
}

// NewClient creates a new YourExchange client
// demo: true for testnet, false for production
func NewClient(apiKey, secretKey string, demo bool) *Client {
    baseURL := "https://api.yourexchange.com"
    if demo {
        baseURL = "https://testnet.yourexchange.com"
    }

    return &Client{
        apiKey:    apiKey,
        secretKey: secretKey,
        baseURL:   baseURL,
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// Name returns the broker name
func (c *Client) Name() string {
    return "YourExchange"
}

// SupportedFeatures returns capabilities
func (c *Client) SupportedFeatures() broker.Features {
    return broker.Features{
        TrailingStop:     true,
        MultipleTP:       false,
        BracketOrders:    true,
        MaxLeverage:      125,
        ReduceOnlyOrders: true,
    }
}
```

### 3.2: Authentication

**File: `youexchange/auth.go`**

```go
package youexchange

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "net/url"
    "sort"
    "strconv"
    "time"
)

// signRequest signs a request according to YourExchange's requirements
// Different exchanges use different signing methods (HMAC-SHA256, RSA, etc.)
func (c *Client) signRequest(params url.Values) string {
    // Add timestamp (most exchanges require this)
    timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
    params.Set("timestamp", timestamp)

    // Sort parameters (many exchanges require alphabetical order)
    keys := make([]string, 0, len(params))
    for key := range params {
        keys = append(keys, key)
    }
    sort.Strings(keys)

    // Build query string
    var query string
    for i, key := range keys {
        if i > 0 {
            query += "&"
        }
        query += key + "=" + params.Get(key)
    }

    // Sign with HMAC-SHA256
    h := hmac.New(sha256.New, []byte(c.secretKey))
    h.Write([]byte(query))
    signature := hex.EncodeToString(h.Sum(nil))

    return signature
}

// makeRequest makes an authenticated request to the API
func (c *Client) makeRequest(method, endpoint string, params url.Values) ([]byte, error) {
    // Add API key
    params.Set("apiKey", c.apiKey)

    // Sign request
    signature := c.signRequest(params)
    params.Set("signature", signature)

    // Build URL
    reqURL := c.baseURL + endpoint
    if method == "GET" {
        reqURL += "?" + params.Encode()
    }

    // Create request
    req, err := http.NewRequest(method, reqURL, nil)
    if err != nil {
        return nil, err
    }

    // Set headers
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("X-API-KEY", c.apiKey)

    // Execute request
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Check for API errors
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
    }

    return body, nil
}
```

### 3.3: Exchange-Specific Types

**File: `youexchange/types.go`**

```go
package youexchange

// YourExchange API response structures
// These are exchange-specific and will be converted to broker types

type BalanceResponse struct {
    Code int            `json:"code"`
    Msg  string         `json:"msg"`
    Data []BalanceData  `json:"data"`
}

type BalanceData struct {
    Asset            string `json:"asset"`
    Balance          string `json:"balance"`
    AvailableMargin  string `json:"available"`
    UnrealizedProfit string `json:"unrealizedPnl"`
}

type PositionResponse struct {
    Code int             `json:"code"`
    Msg  string          `json:"msg"`
    Data []PositionData  `json:"data"`
}

type PositionData struct {
    Symbol           string `json:"symbol"`
    Side             string `json:"positionSide"`  // "LONG" or "SHORT"
    PositionAmt      string `json:"positionAmt"`
    EntryPrice       string `json:"entryPrice"`
    MarkPrice        string `json:"markPrice"`
    LiquidationPrice string `json:"liquidationPrice"`
    Leverage         string `json:"leverage"`
    UnrealizedPnl    string `json:"unrealizedPnl"`
    Margin           string `json:"positionMargin"`
}

// Add more types as needed for orders, etc.
```

---

## Step 4: Implement Core Methods

### 4.1: GetBalance

**File: `youexchange/balance.go`**

```go
package youexchange

import (
    "context"
    "encoding/json"
    "net/url"
    "strconv"
    "time"

    "github.com/agatticelli/trading-go/broker"
)

func (c *Client) GetBalance(ctx context.Context) (*broker.Balance, error) {
    params := url.Values{}

    body, err := c.makeRequest("GET", "/v1/account/balance", params)
    if err != nil {
        return nil, err
    }

    var resp BalanceResponse
    if err := json.Unmarshal(body, &resp); err != nil {
        return nil, err
    }

    if resp.Code != 0 {
        return nil, broker.NewBrokerError(c.Name(), strconv.Itoa(resp.Code), resp.Msg, broker.ErrAPIError)
    }

    if len(resp.Data) == 0 {
        return nil, broker.ErrInsufficientBalance
    }

    // Parse balances (handle string to float conversion)
    data := resp.Data[0]
    balance, _ := strconv.ParseFloat(data.Balance, 64)
    available, _ := strconv.ParseFloat(data.AvailableMargin, 64)
    unrealizedPnL, _ := strconv.ParseFloat(data.UnrealizedProfit, 64)

    return &broker.Balance{
        Asset:         data.Asset,
        Total:         balance,
        Available:     available,
        InUse:         balance - available,
        UnrealizedPnL: unrealizedPnL,
        Timestamp:     time.Now(),
    }, nil
}
```

### 4.2: GetPositions

**File: `youexchange/positions.go`**

```go
package youexchange

import (
    "context"
    "encoding/json"
    "net/url"
    "strconv"
    "time"

    "github.com/agatticelli/trading-go/broker"
    "github.com/agatticelli/trading-common-types"
)

func (c *Client) GetPositions(ctx context.Context, filter *broker.PositionFilter) ([]*broker.Position, error) {
    params := url.Values{}

    if filter != nil && filter.Symbol != "" {
        params.Set("symbol", filter.Symbol)
    }

    body, err := c.makeRequest("GET", "/v1/positions", params)
    if err != nil {
        return nil, err
    }

    var resp PositionResponse
    if err := json.Unmarshal(body, &resp); err != nil {
        return nil, err
    }

    if resp.Code != 0 {
        return nil, broker.NewBrokerError(c.Name(), strconv.Itoa(resp.Code), resp.Msg, broker.ErrAPIError)
    }

    positions := make([]*broker.Position, 0, len(resp.Data))

    for _, data := range resp.Data {
        // Parse numeric fields
        size, _ := strconv.ParseFloat(data.PositionAmt, 64)
        if size == 0 {
            continue // Skip empty positions
        }

        entryPrice, _ := strconv.ParseFloat(data.EntryPrice, 64)
        markPrice, _ := strconv.ParseFloat(data.MarkPrice, 64)
        liqPrice, _ := strconv.ParseFloat(data.LiquidationPrice, 64)
        leverage, _ := strconv.Atoi(data.Leverage)
        unrealizedPnL, _ := strconv.ParseFloat(data.UnrealizedPnl, 64)
        margin, _ := strconv.ParseFloat(data.Margin, 64)

        // Normalize side
        var side types.Side
        if data.Side == "LONG" {
            side = types.SideLong
        } else {
            side = types.SideShort
        }

        // Apply filter if specified
        if filter != nil && filter.Side != nil && side != *filter.Side {
            continue
        }

        positions = append(positions, &broker.Position{
            Symbol:           data.Symbol,
            Side:             side,
            Size:             size,
            EntryPrice:       entryPrice,
            MarkPrice:        markPrice,
            LiquidationPrice: liqPrice,
            Leverage:         leverage,
            UnrealizedPnL:    unrealizedPnL,
            Margin:           margin,
            Timestamp:        time.Now(),
        })
    }

    return positions, nil
}

func (c *Client) GetPosition(ctx context.Context, symbol string) (*broker.Position, error) {
    filter := &broker.PositionFilter{Symbol: symbol}
    positions, err := c.GetPositions(ctx, filter)
    if err != nil {
        return nil, err
    }

    if len(positions) == 0 {
        return nil, broker.ErrPositionNotFound
    }

    return positions[0], nil
}
```

### 4.3: PlaceOrder

**File: `youexchange/orders.go`**

```go
package youexchange

import (
    "context"
    "encoding/json"
    "fmt"
    "net/url"
    "strconv"
    "time"

    "github.com/agatticelli/trading-go/broker"
)

func (c *Client) PlaceOrder(ctx context.Context, order *broker.OrderRequest) (*broker.Order, error) {
    params := url.Values{}
    params.Set("symbol", order.Symbol)
    params.Set("side", normalizeSide(order.Side))
    params.Set("type", normalizeOrderType(order.Type))
    params.Set("quantity", strconv.FormatFloat(order.Size, 'f', -1, 64))

    if order.Price > 0 {
        params.Set("price", strconv.FormatFloat(order.Price, 'f', -1, 64))
    }

    if order.ReduceOnly {
        params.Set("reduceOnly", "true")
    }

    // Handle TP/SL if supported
    if order.StopLoss != nil {
        params.Set("stopPrice", strconv.FormatFloat(order.StopLoss.TriggerPrice, 'f', -1, 64))
        params.Set("stopLossWorkingType", string(order.StopLoss.WorkingType))
    }

    if order.TakeProfit != nil {
        params.Set("takeProfitPrice", strconv.FormatFloat(order.TakeProfit.TriggerPrice, 'f', -1, 64))
        params.Set("takeProfitWorkingType", string(order.TakeProfit.WorkingType))
    }

    body, err := c.makeRequest("POST", "/v1/order", params)
    if err != nil {
        return nil, err
    }

    var resp OrderResponse
    if err := json.Unmarshal(body, &resp); err != nil {
        return nil, err
    }

    if resp.Code != 0 {
        return nil, c.parseAPIError(resp.Code, resp.Msg)
    }

    // Convert to broker.Order
    return &broker.Order{
        ID:        resp.Data.OrderID,
        Symbol:    order.Symbol,
        Side:      order.Side,
        Type:      order.Type,
        Status:    broker.OrderStatusNew,
        Size:      order.Size,
        Price:     order.Price,
        CreatedAt: time.Now(),
    }, nil
}

// Helper: Normalize side to exchange format
func normalizeSide(side types.Side) string {
    if side == types.SideLong {
        return "BUY"
    }
    return "SELL"
}

// Helper: Normalize order type to exchange format
func normalizeOrderType(orderType broker.OrderType) string {
    switch orderType {
    case broker.OrderTypeMarket:
        return "MARKET"
    case broker.OrderTypeLimit:
        return "LIMIT"
    case broker.OrderTypeStop:
        return "STOP"
    case broker.OrderTypeTakeProfit:
        return "TAKE_PROFIT"
    case broker.OrderTypeTrailingStop:
        return "TRAILING_STOP_MARKET"
    default:
        return "MARKET"
    }
}
```

---

## Step 5: Error Handling

Map exchange error codes to broker errors:

```go
func (c *Client) parseAPIError(code int, message string) error {
    baseErr := broker.ErrAPIError

    switch code {
    case 100400, 200001:  // Insufficient balance
        baseErr = broker.ErrInsufficientBalance
    case 100401:  // Invalid symbol
        baseErr = broker.ErrInvalidSymbol
    case 100403:  // Invalid price
        baseErr = broker.ErrInvalidPrice
    case 100404:  // Invalid quantity
        baseErr = broker.ErrInvalidQuantity
    case 100001, 100002:  // Auth errors
        baseErr = broker.ErrAuthFailed
    case 429:  // Rate limit
        baseErr = broker.ErrRateLimited
    }

    return broker.NewBrokerError(c.Name(), strconv.Itoa(code), message, baseErr)
}
```

---

## Step 6: Implement Remaining Methods

Implement all remaining interface methods:

- `GetOrders()`
- `CancelOrder()`
- `CancelAllOrders()`
- `GetCurrentPrice()`
- `SetLeverage()`

Follow the same pattern as above:
1. Build request parameters
2. Make authenticated request
3. Parse response
4. Normalize to broker types
5. Handle errors

---

## Step 7: Testing

### Unit Tests

Create `youexchange/client_test.go`:

```go
package youexchange

import (
    "testing"
)

func TestNormalizeSide(t *testing.T) {
    tests := []struct {
        input types.Side
        want  string
    }{
        {types.SideLong, "BUY"},
        {types.SideShort, "SELL"},
    }

    for _, tt := range tests {
        got := normalizeSide(tt.input)
        if got != tt.want {
            t.Errorf("normalizeSide(%v) = %q, want %q", tt.input, got, tt.want)
        }
    }
}

// Add more tests...
```

### Integration Testing

Test with exchange's testnet:

```go
func TestGetBalance(t *testing.T) {
    apiKey := os.Getenv("YOUREXCHANGE_API_KEY")
    secretKey := os.Getenv("YOUREXCHANGE_SECRET_KEY")

    if apiKey == "" || secretKey == "" {
        t.Skip("API credentials not set")
    }

    client := NewClient(apiKey, secretKey, true)  // demo mode
    ctx := context.Background()

    balance, err := client.GetBalance(ctx)
    if err != nil {
        t.Fatalf("GetBalance() error = %v", err)
    }

    if balance.Asset != "USDT" {
        t.Errorf("Asset = %q, want %q", balance.Asset, "USDT")
    }
}
```

---

## Step 8: Documentation

Create a README for your broker:

**File: `youexchange/README.md`**

```markdown
# YourExchange Broker Implementation

Implementation of the broker.Broker interface for YourExchange.

## Features

- ✅ Account balance
- ✅ Position management
- ✅ Market, Limit, Stop orders
- ✅ Trailing stops
- ✅ Bracket orders (TP/SL)
- ✅ Demo mode support
- ✅ Max leverage: 125x

## Usage

```go
import "github.com/agatticelli/trading-go/yourexchange"

// Create client
client := yourexchange.NewClient(apiKey, secretKey, true)  // true = demo

// Get balance
balance, err := client.GetBalance(ctx)

// Place order
order := &broker.OrderRequest{
    Symbol: "BTC-USDT",
    Side:   broker.SideLong,
    Type:   broker.OrderTypeMarket,
    Size:   0.001,
}
result, err := client.PlaceOrder(ctx, order)
```

## API Credentials

Get your API keys from:
- Demo: https://testnet.yourexchange.com/account/api
- Production: https://yourexchange.com/account/api

## Rate Limits

- Public endpoints: 1200/min
- Private endpoints: 600/min
```

---

## Best Practices

### 1. Use Context for Cancellation

```go
func (c *Client) GetBalance(ctx context.Context) (*broker.Balance, error) {
    // Check context before expensive operations
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    // Make request...
}
```

### 2. Handle String/Number Type Flexibility

Many exchanges return numbers as strings:

```go
// Bad
size := data.Size  // might be string

// Good
size, err := strconv.ParseFloat(data.Size, 64)
if err != nil {
    return nil, fmt.Errorf("invalid size: %w", err)
}
```

### 3. Retry on Network Errors

```go
func (c *Client) makeRequestWithRetry(method, endpoint string, params url.Values) ([]byte, error) {
    maxRetries := 3
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        body, err := c.makeRequest(method, endpoint, params)
        if err == nil {
            return body, nil
        }
        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1))
    }

    return nil, lastErr
}
```

### 4. Log API Requests (Debug Mode)

```go
if c.debug {
    log.Printf("[%s] %s %s?%s", c.Name(), method, endpoint, params.Encode())
}
```

---

## Checklist

Before submitting your broker implementation:

- [ ] All interface methods implemented
- [ ] Error codes properly mapped
- [ ] Types normalized correctly
- [ ] Authentication working
- [ ] Demo mode supported
- [ ] Unit tests written
- [ ] Integration tests passing
- [ ] Documentation complete
- [ ] Examples provided
- [ ] Code reviewed

---

## Common Pitfalls

### 1. Timestamp Format

Different exchanges use different formats:
- Unix milliseconds: `1234567890123`
- Unix seconds: `1234567890`
- ISO 8601: `2024-01-01T00:00:00Z`

### 2. Symbol Format

Normalize to consistent format:
- Exchange: `BTCUSDT`
- Broker: `BTC-USDT`

### 3. Side Naming

Different exchanges use different conventions:
- BUY/SELL
- LONG/SHORT
- 1/-1

Always normalize to `types.SideLong` / `types.SideShort`.

### 4. Order States

Map exchange-specific order states to broker states:
- NEW, PARTIALLY_FILLED, FILLED, CANCELED, EXPIRED, REJECTED

---

## Need Help?

- Check the BingX implementation as a reference
- Read exchange's API documentation carefully
- Test thoroughly on testnet first
- Ask questions in discussions

---

**Good luck implementing your broker!** 🚀
