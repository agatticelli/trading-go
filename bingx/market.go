package bingx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gattimassimo/trading-go/broker"
)

// GetCurrentPrice retrieves current market price for a symbol
func (c *Client) GetCurrentPrice(ctx context.Context, symbol string) (float64, error) {
	params := map[string]string{
		"symbol": symbol,
	}

	body, err := c.makeRequest(ctx, "GET", EndpointPrice, params)
	if err != nil {
		return 0, err
	}

	var response PriceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, broker.NewBrokerError("bingx", "PARSE_ERROR", "Failed to parse price response", err)
	}

	if response.Code != APISuccessCode {
		return 0, broker.NewBrokerError("bingx", fmt.Sprintf("API_%d", response.Code), response.Msg, nil)
	}

	price, err := strconv.ParseFloat(response.Data.Price, 64)
	if err != nil {
		return 0, broker.NewBrokerError("bingx", "PARSE_ERROR", "Failed to parse price value", err)
	}

	return price, nil
}

// SetLeverage sets leverage for a symbol
func (c *Client) SetLeverage(ctx context.Context, symbol string, side string, leverage int) error {
	params := map[string]string{
		"symbol":   symbol,
		"side":     side,
		"leverage": strconv.Itoa(leverage),
	}

	body, err := c.makeRequest(ctx, "POST", EndpointLeverage, params)
	if err != nil {
		return err
	}

	var response LeverageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return broker.NewBrokerError("bingx", "PARSE_ERROR", "Failed to parse leverage response", err)
	}

	if response.Code != APISuccessCode {
		return broker.NewBrokerError("bingx", fmt.Sprintf("API_%d", response.Code), response.Msg, nil)
	}

	return nil
}
