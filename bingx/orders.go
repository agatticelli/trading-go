package bingx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gattimassimo/trading-go/broker"
)

// PlaceOrder places a new order
func (c *Client) PlaceOrder(ctx context.Context, order *broker.OrderRequest) (*broker.Order, error) {
	// Convert broker types to BingX types
	side := "BUY"
	positionSide := "LONG"
	if order.Side == broker.SideShort {
		side = "SELL"
		positionSide = "SHORT"
	}

	// Build BingX order request
	params := map[string]string{
		"symbol":       order.Symbol,
		"side":         side,
		"positionSide": positionSide,
		"type":         string(order.Type),
		"quantity":     fmt.Sprintf("%.8f", order.Size),
	}

	// Add optional parameters
	if order.Price > 0 {
		params["price"] = fmt.Sprintf("%.8f", order.Price)
	}
	if order.StopPrice > 0 {
		params["stopPrice"] = fmt.Sprintf("%.8f", order.StopPrice)
	}
	if order.TimeInForce != "" {
		params["timeInForce"] = string(order.TimeInForce)
	} else if order.Type == broker.OrderTypeLimit {
		params["timeInForce"] = "GTC" // Default for limit orders
	}
	if order.ReduceOnly {
		params["reduceOnly"] = "true"
	}

	// Execute request
	body, err := c.makeRequest(ctx, "POST", EndpointPlaceOrder, params)
	if err != nil {
		return nil, err
	}

	var response OrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, broker.NewBrokerError("bingx", "PARSE_ERROR", "Failed to parse order response", err)
	}

	if response.Code != APISuccessCode {
		return nil, broker.NewBrokerError("bingx", fmt.Sprintf("API_%d", response.Code), response.Msg, nil)
	}

	// Convert response to broker.Order
	price, _ := strconv.ParseFloat(response.Data.Price, 64)
	size, _ := strconv.ParseFloat(response.Data.Quantity, 64)

	var brokerSide broker.Side
	if response.Data.PositionSide == "LONG" {
		brokerSide = broker.SideLong
	} else {
		brokerSide = broker.SideShort
	}

	return &broker.Order{
		ID:          fmt.Sprintf("%d", response.Data.OrderId),
		Symbol:      response.Data.Symbol,
		Side:        brokerSide,
		Type:        broker.OrderType(response.Data.Type),
		Status:      broker.OrderStatus(response.Data.Status),
		Size:        size,
		Price:       price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// GetOrders retrieves open orders
func (c *Client) GetOrders(ctx context.Context, filter *broker.OrderFilter) ([]*broker.Order, error) {
	params := make(map[string]string)
	if filter != nil && filter.Symbol != "" {
		params["symbol"] = filter.Symbol
	}

	body, err := c.makeRequest(ctx, "GET", EndpointOpenOrders, params)
	if err != nil {
		return nil, err
	}

	var response OpenOrdersResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, broker.NewBrokerError("bingx", "PARSE_ERROR", "Failed to parse orders response", err)
	}

	if response.Code != APISuccessCode {
		return nil, broker.NewBrokerError("bingx", fmt.Sprintf("API_%d", response.Code), response.Msg, nil)
	}

	var orders []*broker.Order
	for _, o := range response.Data {
		// Determine side
		var side broker.Side
		if o.PositionSide == "LONG" {
			side = broker.SideLong
		} else {
			side = broker.SideShort
		}

		// Apply filter if specified
		if filter != nil && filter.Side != nil && *filter.Side != broker.OrderStatus(o.Status) {
			continue
		}

		// Parse fields
		size, _ := strconv.ParseFloat(o.Quantity, 64)
		price, _ := strconv.ParseFloat(o.Price, 64)
		stopPrice, _ := strconv.ParseFloat(o.StopPrice, 64)
		filledSize, _ := strconv.ParseFloat(o.ExecutedQty, 64)
		avgPrice, _ := strconv.ParseFloat(o.AvgPrice, 64)

		orders = append(orders, &broker.Order{
			ID:            fmt.Sprintf("%d", o.OrderId),
			ClientOrderID: o.ClientOrderID,
			Symbol:        o.Symbol,
			Side:          side,
			Type:          broker.OrderType(o.Type),
			Status:        broker.OrderStatus(o.Status),
			Size:          size,
			Price:         price,
			StopPrice:     stopPrice,
			FilledSize:    filledSize,
			AveragePrice:  avgPrice,
			TimeInForce:   broker.TimeInForce(o.TimeInForce),
			CreatedAt:     time.Unix(o.Time/1000, 0),
			UpdatedAt:     time.Unix(o.UpdateTime/1000, 0),
		})
	}

	return orders, nil
}

// CancelOrder cancels a specific order
func (c *Client) CancelOrder(ctx context.Context, symbol string, orderID string) error {
	params := map[string]string{
		"symbol":  symbol,
		"orderId": orderID,
	}

	body, err := c.makeRequest(ctx, "DELETE", EndpointPlaceOrder, params)
	if err != nil {
		return err
	}

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return broker.NewBrokerError("bingx", "PARSE_ERROR", "Failed to parse cancel response", err)
	}

	if response.Code != APISuccessCode {
		return broker.NewBrokerError("bingx", fmt.Sprintf("API_%d", response.Code), response.Msg, nil)
	}

	return nil
}

// CancelAllOrders cancels all orders for a symbol (or all symbols if empty)
func (c *Client) CancelAllOrders(ctx context.Context, symbol string) error {
	params := make(map[string]string)
	if symbol != "" {
		params["symbol"] = symbol
	}

	body, err := c.makeRequest(ctx, "DELETE", EndpointCancelAll, params)
	if err != nil {
		return err
	}

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return broker.NewBrokerError("bingx", "PARSE_ERROR", "Failed to parse cancel all response", err)
	}

	if response.Code != APISuccessCode {
		return broker.NewBrokerError("bingx", fmt.Sprintf("API_%d", response.Code), response.Msg, nil)
	}

	return nil
}
