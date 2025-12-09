package broker

import (
	"context"
)

// Broker defines the interface all exchange implementations must satisfy
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

// Features describes broker capabilities
type Features struct {
	TrailingStop     bool
	MultipleTP       bool
	BracketOrders    bool
	MaxLeverage      int
	ReduceOnlyOrders bool
}

// PositionFilter for filtering positions
type PositionFilter struct {
	Symbol string
	Side   *Side // Filter by side (nil = all)
}

// OrderFilter for filtering orders
type OrderFilter struct {
	Symbol string
	Side   *OrderStatus
	Status *OrderStatus
}
