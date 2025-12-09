package bingx

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// BingX API response structures

type BalanceData struct {
	UserId           string `json:"userId"`
	Asset            string `json:"asset"`
	Balance          string `json:"balance"`
	Equity           string `json:"equity"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	RealisedProfit   string `json:"realisedProfit"`
	AvailableMargin  string `json:"availableMargin"`
	UsedMargin       string `json:"usedMargin"`
	FreezedMargin    string `json:"freezedMargin"`
	ShortUid         string `json:"shortUid"`
}

type BalanceResponse struct {
	Code int           `json:"code"`
	Data []BalanceData `json:"data"`
	Msg  string        `json:"msg"`
}

type PositionData struct {
	Symbol            string          `json:"symbol"`
	PositionSide      string          `json:"positionSide"`
	PositionAmt       string          `json:"positionAmt"`
	AvailableAmt      string          `json:"availableAmt"`
	UnrealizedProfit  string          `json:"unrealizedProfit"`
	RealisedProfit    string          `json:"realisedProfit"`
	InitialMargin     string          `json:"initialMargin"`
	MaintenanceMargin string          `json:"maintenanceMargin"`
	PositionValue     string          `json:"positionValue"`
	Leverage          json.RawMessage `json:"leverage"`          // Can be string or number
	IsolatedMargin    string          `json:"isolatedMargin"`
	AvgPrice          string          `json:"avgPrice"`
	MaxNotionalValue  string          `json:"maxNotionalValue"`
	BidNotional       string          `json:"bidNotional"`
	AskNotional       string          `json:"askNotional"`
	LiquidationPrice  json.RawMessage `json:"liquidationPrice"`  // Can be string or number
	MarkPrice         string          `json:"markPrice"`
}

// GetLeverageFloat returns the leverage as float64, handling both string and number formats
func (p *PositionData) GetLeverageFloat() (float64, error) {
	// First try to unmarshal as a string
	var leverageStr string
	if err := json.Unmarshal(p.Leverage, &leverageStr); err == nil {
		return strconv.ParseFloat(leverageStr, 64)
	}

	// If that fails, try to unmarshal as a number
	var leverageFloat float64
	if err := json.Unmarshal(p.Leverage, &leverageFloat); err == nil {
		return leverageFloat, nil
	}

	return 0, fmt.Errorf("unable to parse leverage: %s", string(p.Leverage))
}

// GetLiquidationPriceFloat returns the liquidation price as float64, handling both string and number formats
func (p *PositionData) GetLiquidationPriceFloat() (float64, error) {
	// First try to unmarshal as a string
	var priceStr string
	if err := json.Unmarshal(p.LiquidationPrice, &priceStr); err == nil {
		if priceStr == "" {
			return 0, nil
		}
		return strconv.ParseFloat(priceStr, 64)
	}

	// If that fails, try to unmarshal as a number
	var priceFloat float64
	if err := json.Unmarshal(p.LiquidationPrice, &priceFloat); err == nil {
		return priceFloat, nil
	}

	return 0, fmt.Errorf("unable to parse liquidation price: %s", string(p.LiquidationPrice))
}

type PositionsResponse struct {
	Code int            `json:"code"`
	Data []PositionData `json:"data"`
	Msg  string         `json:"msg"`
}

type BingXOrderRequest struct {
	Symbol       string `json:"symbol"`
	Side         string `json:"side"`         // BUY, SELL
	PositionSide string `json:"positionSide"` // LONG, SHORT
	Type         string `json:"type"`         // LIMIT, MARKET, STOP, TAKE_PROFIT
	Quantity     string `json:"quantity"`
	Price        string `json:"price,omitempty"`
	StopPrice    string `json:"stopPrice,omitempty"`
	TimeInForce  string `json:"timeInForce,omitempty"` // GTC, IOC, FOK
}

type OrderResponse struct {
	Code int `json:"code"`
	Data struct {
		OrderId      int64  `json:"orderId"`
		Symbol       string `json:"symbol"`
		Side         string `json:"side"`
		PositionSide string `json:"positionSide"`
		Type         string `json:"type"`
		Quantity     string `json:"origQty"`
		Price        string `json:"price"`
		Status       string `json:"status"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type OpenOrderData struct {
	OrderId           int64  `json:"orderId"`
	Symbol            string `json:"symbol"`
	Side              string `json:"side"`
	PositionSide      string `json:"positionSide"`
	Type              string `json:"type"`
	Quantity          string `json:"origQty"`
	Price             string `json:"price"`
	StopPrice         string `json:"stopPrice"`
	ExecutedQty       string `json:"executedQty"`
	AvgPrice          string `json:"avgPrice"`
	Status            string `json:"status"`
	TimeInForce       string `json:"timeInForce"`
	ClientOrderID     string `json:"clientOrderId"`
	WorkingType       string `json:"workingType"`
	Time              int64  `json:"time"`
	UpdateTime        int64  `json:"updateTime"`
}

type OpenOrdersResponse struct {
	Code int `json:"code"`
	Data struct {
		Orders []OpenOrderData `json:"orders"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type PriceResponse struct {
	Code int `json:"code"`
	Data struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type LeverageResponse struct {
	Code int `json:"code"`
	Data struct {
		Symbol   string `json:"symbol"`
		Leverage string `json:"leverage"`
	} `json:"data"`
	Msg string `json:"msg"`
}
