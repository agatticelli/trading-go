package broker

import (
	"github.com/agatticelli/trading-common-types"
)

// Re-export common types for backward compatibility
type (
	Side        = types.Side
	OrderType   = types.OrderType
	OrderStatus = types.OrderStatus
	TimeInForce = types.TimeInForce
	WorkingType = types.WorkingType

	Balance      = types.Balance
	Position     = types.Position
	Order        = types.Order
	OrderRequest = types.OrderRequest

	StopLossConfig   = types.StopLossConfig
	TakeProfitConfig = types.TakeProfitConfig
	TrailingConfig   = types.TrailingConfig
)

// Re-export constants
const (
	SideLong  = types.SideLong
	SideShort = types.SideShort

	OrderTypeMarket       = types.OrderTypeMarket
	OrderTypeLimit        = types.OrderTypeLimit
	OrderTypeStop         = types.OrderTypeStop
	OrderTypeTakeProfit   = types.OrderTypeTakeProfit
	OrderTypeTrailingStop = types.OrderTypeTrailingStop

	OrderStatusNew             = types.OrderStatusNew
	OrderStatusPartiallyFilled = types.OrderStatusPartiallyFilled
	OrderStatusFilled          = types.OrderStatusFilled
	OrderStatusCanceled        = types.OrderStatusCanceled
	OrderStatusRejected        = types.OrderStatusRejected
	OrderStatusExpired         = types.OrderStatusExpired

	TimeInForceGTC = types.TimeInForceGTC
	TimeInForceIOC = types.TimeInForceIOC
	TimeInForceFOK = types.TimeInForceFOK

	WorkingTypeMark = types.WorkingTypeMark
	WorkingTypeLast = types.WorkingTypeLast
)
