package bingx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gattimassimo/trading-go/broker"
)

// GetBalance retrieves account balance
func (c *Client) GetBalance(ctx context.Context) (*broker.Balance, error) {
	body, err := c.makeRequest(ctx, "GET", EndpointBalance, nil)
	if err != nil {
		return nil, err
	}

	var response BalanceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, broker.NewBrokerError("bingx", "PARSE_ERROR", "Failed to parse balance response", err)
	}

	if response.Code != APISuccessCode {
		return nil, broker.NewBrokerError("bingx", fmt.Sprintf("API_%d", response.Code), response.Msg, nil)
	}

	if len(response.Data) == 0 {
		return nil, broker.NewBrokerError("bingx", "NO_DATA", "No balance data returned", nil)
	}

	// Get USDT balance (assuming first entry is USDT)
	data := response.Data[0]

	total, _ := strconv.ParseFloat(data.Equity, 64)
	available, _ := strconv.ParseFloat(data.AvailableMargin, 64)
	inUse, _ := strconv.ParseFloat(data.UsedMargin, 64)
	unrealizedPnL, _ := strconv.ParseFloat(data.UnrealizedProfit, 64)
	realizedPnL, _ := strconv.ParseFloat(data.RealisedProfit, 64)

	return &broker.Balance{
		Asset:         data.Asset,
		Total:         total,
		Available:     available,
		InUse:         inUse,
		UnrealizedPnL: unrealizedPnL,
		RealizedPnL:   realizedPnL,
		Timestamp:     time.Now(),
	}, nil
}
