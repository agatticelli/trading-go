package bingx

import (
	"encoding/json"
	"testing"
)

func TestPositionData_GetLeverageFloat(t *testing.T) {
	tests := []struct {
		name         string
		leverageJSON string
		want         float64
		wantErr      bool
	}{
		{
			name:         "Leverage as string",
			leverageJSON: `"10"`,
			want:         10.0,
			wantErr:      false,
		},
		{
			name:         "Leverage as number",
			leverageJSON: `25`,
			want:         25.0,
			wantErr:      false,
		},
		{
			name:         "Leverage as float string",
			leverageJSON: `"50.5"`,
			want:         50.5,
			wantErr:      false,
		},
		{
			name:         "Leverage as float number",
			leverageJSON: `125.0`,
			want:         125.0,
			wantErr:      false,
		},
		{
			name:         "High leverage",
			leverageJSON: `"125"`,
			want:         125.0,
			wantErr:      false,
		},
		{
			name:         "Invalid leverage - non-numeric string",
			leverageJSON: `"invalid"`,
			want:         0,
			wantErr:      true,
		},
		{
			name:         "Invalid leverage - null",
			leverageJSON: `null`,
			want:         0,
			wantErr:      true,
		},
		{
			name:         "Invalid leverage - boolean",
			leverageJSON: `true`,
			want:         0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := &PositionData{
				Leverage: json.RawMessage(tt.leverageJSON),
			}

			got, err := pos.GetLeverageFloat()

			if tt.wantErr {
				if err == nil {
					t.Error("GetLeverageFloat() error = nil, want error")
				}
				return
			}

			if err != nil {
				t.Errorf("GetLeverageFloat() error = %v, want nil", err)
				return
			}

			if got != tt.want {
				t.Errorf("GetLeverageFloat() = %.2f, want %.2f", got, tt.want)
			}
		})
	}
}

func TestPositionData_GetLiquidationPriceFloat(t *testing.T) {
	tests := []struct {
		name         string
		priceJSON    string
		want         float64
		wantErr      bool
	}{
		{
			name:      "Price as string",
			priceJSON: `"45000.50"`,
			want:      45000.50,
			wantErr:   false,
		},
		{
			name:      "Price as number",
			priceJSON: `42000.0`,
			want:      42000.0,
			wantErr:   false,
		},
		{
			name:      "Price as integer string",
			priceJSON: `"50000"`,
			want:      50000.0,
			wantErr:   false,
		},
		{
			name:      "Price as integer",
			priceJSON: `48000`,
			want:      48000.0,
			wantErr:   false,
		},
		{
			name:      "Empty string (no liquidation)",
			priceJSON: `""`,
			want:      0,
			wantErr:   false,
		},
		{
			name:      "Very high price",
			priceJSON: `"999999.99"`,
			want:      999999.99,
			wantErr:   false,
		},
		{
			name:      "Very low price",
			priceJSON: `"0.01"`,
			want:      0.01,
			wantErr:   false,
		},
		{
			name:      "Invalid price - non-numeric string",
			priceJSON: `"invalid"`,
			want:      0,
			wantErr:   true,
		},
		{
			name:      "Null price (no liquidation set)",
			priceJSON: `null`,
			want:      0,
			wantErr:   false,
		},
		{
			name:      "Invalid price - boolean",
			priceJSON: `false`,
			want:      0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := &PositionData{
				LiquidationPrice: json.RawMessage(tt.priceJSON),
			}

			got, err := pos.GetLiquidationPriceFloat()

			if tt.wantErr {
				if err == nil {
					t.Error("GetLiquidationPriceFloat() error = nil, want error")
				}
				return
			}

			if err != nil {
				t.Errorf("GetLiquidationPriceFloat() error = %v, want nil", err)
				return
			}

			if got != tt.want {
				t.Errorf("GetLiquidationPriceFloat() = %.2f, want %.2f", got, tt.want)
			}
		})
	}
}

func TestPositionData_GetLeverageFloat_RealWorldData(t *testing.T) {
	// Test with actual response formats from BingX API
	tests := []struct {
		name     string
		jsonData string
		want     float64
	}{
		{
			name: "BingX response with string leverage",
			jsonData: `{
				"symbol": "BTC-USDT",
				"leverage": "10",
				"positionAmt": "0.1"
			}`,
			want: 10.0,
		},
		{
			name: "BingX response with numeric leverage",
			jsonData: `{
				"symbol": "ETH-USDT",
				"leverage": 25,
				"positionAmt": "1.5"
			}`,
			want: 25.0,
		},
		{
			name: "High leverage position",
			jsonData: `{
				"symbol": "BTC-USDT",
				"leverage": "125",
				"positionAmt": "0.01"
			}`,
			want: 125.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pos PositionData
			if err := json.Unmarshal([]byte(tt.jsonData), &pos); err != nil {
				t.Fatalf("Failed to unmarshal test data: %v", err)
			}

			got, err := pos.GetLeverageFloat()
			if err != nil {
				t.Errorf("GetLeverageFloat() error = %v, want nil", err)
				return
			}

			if got != tt.want {
				t.Errorf("GetLeverageFloat() = %.2f, want %.2f", got, tt.want)
			}
		})
	}
}

func TestPositionData_GetLiquidationPriceFloat_RealWorldData(t *testing.T) {
	// Test with actual response formats from BingX API
	tests := []struct {
		name     string
		jsonData string
		want     float64
	}{
		{
			name: "Position with liquidation price string",
			jsonData: `{
				"symbol": "BTC-USDT",
				"liquidationPrice": "42000.50",
				"positionAmt": "0.1"
			}`,
			want: 42000.50,
		},
		{
			name: "Position with liquidation price number",
			jsonData: `{
				"symbol": "ETH-USDT",
				"liquidationPrice": 2800.0,
				"positionAmt": "1.5"
			}`,
			want: 2800.0,
		},
		{
			name: "Position with empty liquidation price",
			jsonData: `{
				"symbol": "BTC-USDT",
				"liquidationPrice": "",
				"positionAmt": "0.1"
			}`,
			want: 0.0,
		},
		{
			name: "Very precise liquidation price",
			jsonData: `{
				"symbol": "BTC-USDT",
				"liquidationPrice": "45123.456789",
				"positionAmt": "0.1"
			}`,
			want: 45123.456789,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pos PositionData
			if err := json.Unmarshal([]byte(tt.jsonData), &pos); err != nil {
				t.Fatalf("Failed to unmarshal test data: %v", err)
			}

			got, err := pos.GetLiquidationPriceFloat()
			if err != nil {
				t.Errorf("GetLiquidationPriceFloat() error = %v, want nil", err)
				return
			}

			if got != tt.want {
				t.Errorf("GetLiquidationPriceFloat() = %.6f, want %.6f", got, tt.want)
			}
		})
	}
}

func TestBalanceResponse_Unmarshal(t *testing.T) {
	// Test unmarshaling a complete balance response
	jsonData := `{
		"code": 0,
		"msg": "success",
		"data": [
			{
				"userId": "123456",
				"asset": "USDT",
				"balance": "1000.00",
				"equity": "1050.00",
				"unrealizedProfit": "50.00",
				"realisedProfit": "100.00",
				"availableMargin": "950.00",
				"usedMargin": "100.00"
			}
		]
	}`

	var resp BalanceResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
	if resp.Msg != "success" {
		t.Errorf("Msg = %q, want %q", resp.Msg, "success")
	}
	if len(resp.Data) != 1 {
		t.Fatalf("len(Data) = %d, want 1", len(resp.Data))
	}

	data := resp.Data[0]
	if data.Asset != "USDT" {
		t.Errorf("Asset = %q, want %q", data.Asset, "USDT")
	}
	if data.Balance != "1000.00" {
		t.Errorf("Balance = %q, want %q", data.Balance, "1000.00")
	}
}

func TestPositionsResponse_Unmarshal(t *testing.T) {
	// Test unmarshaling a complete positions response with mixed leverage/price formats
	jsonData := `{
		"code": 0,
		"msg": "success",
		"data": [
			{
				"symbol": "BTC-USDT",
				"positionSide": "LONG",
				"positionAmt": "0.1",
				"leverage": "10",
				"avgPrice": "45000.00",
				"markPrice": "46000.00",
				"liquidationPrice": "42000.50",
				"unrealizedProfit": "100.00"
			},
			{
				"symbol": "ETH-USDT",
				"positionSide": "SHORT",
				"positionAmt": "1.5",
				"leverage": 25,
				"avgPrice": "3000.00",
				"markPrice": "2950.00",
				"liquidationPrice": 3200.0,
				"unrealizedProfit": "75.00"
			}
		]
	}`

	var resp PositionsResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("len(Data) = %d, want 2", len(resp.Data))
	}

	// Test first position (string formats)
	pos1 := resp.Data[0]
	if pos1.Symbol != "BTC-USDT" {
		t.Errorf("Position 1 Symbol = %q, want %q", pos1.Symbol, "BTC-USDT")
	}
	lev1, err := pos1.GetLeverageFloat()
	if err != nil {
		t.Errorf("Position 1 GetLeverageFloat() error = %v", err)
	}
	if lev1 != 10.0 {
		t.Errorf("Position 1 Leverage = %.2f, want 10.00", lev1)
	}
	liq1, err := pos1.GetLiquidationPriceFloat()
	if err != nil {
		t.Errorf("Position 1 GetLiquidationPriceFloat() error = %v", err)
	}
	if liq1 != 42000.50 {
		t.Errorf("Position 1 LiquidationPrice = %.2f, want 42000.50", liq1)
	}

	// Test second position (numeric formats)
	pos2 := resp.Data[1]
	if pos2.Symbol != "ETH-USDT" {
		t.Errorf("Position 2 Symbol = %q, want %q", pos2.Symbol, "ETH-USDT")
	}
	lev2, err := pos2.GetLeverageFloat()
	if err != nil {
		t.Errorf("Position 2 GetLeverageFloat() error = %v", err)
	}
	if lev2 != 25.0 {
		t.Errorf("Position 2 Leverage = %.2f, want 25.00", lev2)
	}
	liq2, err := pos2.GetLiquidationPriceFloat()
	if err != nil {
		t.Errorf("Position 2 GetLiquidationPriceFloat() error = %v", err)
	}
	if liq2 != 3200.0 {
		t.Errorf("Position 2 LiquidationPrice = %.2f, want 3200.00", liq2)
	}
}
