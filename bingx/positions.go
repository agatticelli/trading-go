package bingx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gattimassimo/trading-go/broker"
)

// GetPositions retrieves all open positions
func (c *Client) GetPositions(ctx context.Context, filter *broker.PositionFilter) ([]*broker.Position, error) {
	params := make(map[string]string)
	if filter != nil && filter.Symbol != "" {
		params["symbol"] = filter.Symbol
	}

	body, err := c.makeRequest(ctx, "GET", EndpointPositions, params)
	if err != nil {
		return nil, err
	}

	var response PositionsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, broker.NewBrokerError("bingx", "PARSE_ERROR", "Failed to parse positions response", err)
	}

	if response.Code != APISuccessCode {
		return nil, broker.NewBrokerError("bingx", fmt.Sprintf("API_%d", response.Code), response.Msg, nil)
	}

	var positions []*broker.Position
	for _, pos := range response.Data {
		// Parse position amount
		size, _ := strconv.ParseFloat(pos.PositionAmt, 64)

		// Skip positions with zero size
		if size == 0 {
			continue
		}

		// Determine side
		var side broker.Side
		if pos.PositionSide == "LONG" {
			side = broker.SideLong
		} else {
			side = broker.SideShort
		}

		// Apply filter if specified
		if filter != nil && filter.Side != nil && *filter.Side != side {
			continue
		}

		// Parse other fields
		entryPrice, _ := strconv.ParseFloat(pos.AvgPrice, 64)
		markPrice, _ := strconv.ParseFloat(pos.MarkPrice, 64)
		unrealizedPnL, _ := strconv.ParseFloat(pos.UnrealizedProfit, 64)
		realizedPnL, _ := strconv.ParseFloat(pos.RealisedProfit, 64)
		margin, _ := strconv.ParseFloat(pos.InitialMargin, 64)
		maintenanceMargin, _ := strconv.ParseFloat(pos.MaintenanceMargin, 64)

		// Parse leverage (can be string or number)
		leverage, err := pos.GetLeverageFloat()
		if err != nil {
			leverage = 0
		}

		// Parse liquidation price (can be string or number)
		liquidationPrice, err := pos.GetLiquidationPriceFloat()
		if err != nil {
			liquidationPrice = 0
		}

		positions = append(positions, &broker.Position{
			Symbol:            pos.Symbol,
			Side:              side,
			Size:              size,
			EntryPrice:        entryPrice,
			MarkPrice:         markPrice,
			LiquidationPrice:  liquidationPrice,
			Leverage:          int(leverage),
			UnrealizedPnL:     unrealizedPnL,
			RealizedPnL:       realizedPnL,
			Margin:            margin,
			MaintenanceMargin: maintenanceMargin,
			Timestamp:         time.Now(),
		})
	}

	return positions, nil
}

// GetPosition retrieves a single position by symbol
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
