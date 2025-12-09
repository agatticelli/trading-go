package broker

import "time"

// Balance represents account balance information
type Balance struct {
	Asset         string
	Total         float64 // Total balance
	Available     float64 // Available for trading
	InUse         float64 // Currently in positions
	UnrealizedPnL float64
	RealizedPnL   float64
	Timestamp     time.Time
}

// Position represents an open trading position
type Position struct {
	Symbol            string
	Side              Side
	Size              float64 // Position size (positive)
	EntryPrice        float64
	MarkPrice         float64
	LiquidationPrice  float64
	Leverage          int
	UnrealizedPnL     float64
	RealizedPnL       float64
	Margin            float64
	MaintenanceMargin float64
	Timestamp         time.Time
}

// Order represents a trading order
type Order struct {
	ID            string
	ClientOrderID string
	Symbol        string
	Side          Side
	Type          OrderType
	Status        OrderStatus
	Size          float64
	Price         float64 // Limit price (0 for market)
	StopPrice     float64 // Trigger price (for stop orders)
	FilledSize    float64
	AveragePrice  float64
	ReduceOnly    bool
	TimeInForce   TimeInForce
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// OrderRequest represents a request to place an order
type OrderRequest struct {
	Symbol      string
	Side        Side
	Type        OrderType
	Size        float64
	Price       float64 // Required for LIMIT orders
	StopPrice   float64 // Required for STOP/TAKE_PROFIT orders
	TimeInForce TimeInForce
	ReduceOnly  bool

	// Advanced order features (BingX-specific, optional)
	StopLoss   *StopLossConfig
	TakeProfit *TakeProfitConfig
	Trailing   *TrailingConfig
}

// StopLossConfig for attaching SL to orders
type StopLossConfig struct {
	TriggerPrice float64
	OrderPrice   float64     // Limit price (0 for market)
	WorkingType  WorkingType // MARK_PRICE or LAST_PRICE
}

// TakeProfitConfig for attaching TP to orders
type TakeProfitConfig struct {
	TriggerPrice float64
	OrderPrice   float64
	WorkingType  WorkingType
}

// TrailingConfig for trailing stop orders
type TrailingConfig struct {
	ActivationPrice float64 // Price where trailing starts
	CallbackRate    float64 // Trailing percentage (0.005 = 0.5%)
}

// Side represents position direction
type Side string

const (
	SideLong  Side = "LONG"
	SideShort Side = "SHORT"
)

// OrderType represents order type
type OrderType string

const (
	OrderTypeMarket       OrderType = "MARKET"
	OrderTypeLimit        OrderType = "LIMIT"
	OrderTypeStop         OrderType = "STOP"
	OrderTypeTakeProfit   OrderType = "TAKE_PROFIT"
	OrderTypeTrailingStop OrderType = "TRAILING_STOP_MARKET"
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

// TimeInForce represents time in force
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC" // Good Till Cancel
	TimeInForceIOC TimeInForce = "IOC" // Immediate or Cancel
	TimeInForceFOK TimeInForce = "FOK" // Fill or Kill
)

// WorkingType represents working price type
type WorkingType string

const (
	WorkingTypeMark WorkingType = "MARK_PRICE"
	WorkingTypeLast WorkingType = "LAST_PRICE"
)
