package broker

import (
	"errors"
	"testing"
)

func TestBrokerError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *BrokerError
		wantMsg string
	}{
		{
			name: "BingX authentication error",
			err: &BrokerError{
				Broker:  "BingX",
				Code:    "100001",
				Message: "Invalid API key",
				Err:     ErrAuthFailed,
			},
			wantMsg: "BingX error [100001]: Invalid API key",
		},
		{
			name: "Generic API error",
			err: &BrokerError{
				Broker:  "TestBroker",
				Code:    "500",
				Message: "Internal server error",
				Err:     ErrAPIError,
			},
			wantMsg: "TestBroker error [500]: Internal server error",
		},
		{
			name: "Rate limit error",
			err: &BrokerError{
				Broker:  "BingX",
				Code:    "429",
				Message: "Too many requests",
				Err:     ErrRateLimited,
			},
			wantMsg: "BingX error [429]: Too many requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("BrokerError.Error() = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestBrokerError_Unwrap(t *testing.T) {
	tests := []struct {
		name    string
		err     *BrokerError
		wantErr error
	}{
		{
			name: "Unwrap auth error",
			err: &BrokerError{
				Broker:  "BingX",
				Code:    "100001",
				Message: "Invalid API key",
				Err:     ErrAuthFailed,
			},
			wantErr: ErrAuthFailed,
		},
		{
			name: "Unwrap insufficient balance",
			err: &BrokerError{
				Broker:  "BingX",
				Code:    "100400",
				Message: "Not enough balance",
				Err:     ErrInsufficientBalance,
			},
			wantErr: ErrInsufficientBalance,
		},
		{
			name: "Unwrap nil error",
			err: &BrokerError{
				Broker:  "BingX",
				Code:    "200",
				Message: "Success",
				Err:     nil,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Unwrap(); got != tt.wantErr {
				t.Errorf("BrokerError.Unwrap() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestNewBrokerError(t *testing.T) {
	tests := []struct {
		name       string
		broker     string
		code       string
		message    string
		err        error
		wantBroker string
		wantCode   string
		wantMsg    string
		wantErr    error
	}{
		{
			name:       "Create auth error",
			broker:     "BingX",
			code:       "100001",
			message:    "Invalid signature",
			err:        ErrAuthFailed,
			wantBroker: "BingX",
			wantCode:   "100001",
			wantMsg:    "Invalid signature",
			wantErr:    ErrAuthFailed,
		},
		{
			name:       "Create position not found error",
			broker:     "TestExchange",
			code:       "404",
			message:    "Position not found",
			err:        ErrPositionNotFound,
			wantBroker: "TestExchange",
			wantCode:   "404",
			wantMsg:    "Position not found",
			wantErr:    ErrPositionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBrokerError(tt.broker, tt.code, tt.message, tt.err)

			if got.Broker != tt.wantBroker {
				t.Errorf("Broker = %q, want %q", got.Broker, tt.wantBroker)
			}
			if got.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", got.Code, tt.wantCode)
			}
			if got.Message != tt.wantMsg {
				t.Errorf("Message = %q, want %q", got.Message, tt.wantMsg)
			}
			if got.Err != tt.wantErr {
				t.Errorf("Err = %v, want %v", got.Err, tt.wantErr)
			}
		})
	}
}

func TestStandardErrors(t *testing.T) {
	// Test that standard errors are defined
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrInvalidSymbol", ErrInvalidSymbol, "invalid symbol"},
		{"ErrInvalidPrice", ErrInvalidPrice, "invalid price"},
		{"ErrInvalidQuantity", ErrInvalidQuantity, "invalid quantity"},
		{"ErrInsufficientBalance", ErrInsufficientBalance, "insufficient balance"},
		{"ErrPositionNotFound", ErrPositionNotFound, "position not found"},
		{"ErrOrderNotFound", ErrOrderNotFound, "order not found"},
		{"ErrLeverageTooHigh", ErrLeverageTooHigh, "leverage exceeds maximum"},
		{"ErrAuthFailed", ErrAuthFailed, "authentication failed"},
		{"ErrRateLimited", ErrRateLimited, "rate limited"},
		{"ErrAPIError", ErrAPIError, "API error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s is nil", tt.name)
			}
			if tt.err.Error() != tt.msg {
				t.Errorf("%s.Error() = %q, want %q", tt.name, tt.err.Error(), tt.msg)
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test that errors.Is works correctly with BrokerError
	tests := []struct {
		name      string
		err       error
		target    error
		wantMatch bool
	}{
		{
			name: "Match wrapped auth error",
			err: &BrokerError{
				Broker:  "BingX",
				Code:    "100001",
				Message: "Invalid API key",
				Err:     ErrAuthFailed,
			},
			target:    ErrAuthFailed,
			wantMatch: true,
		},
		{
			name: "Don't match different error",
			err: &BrokerError{
				Broker:  "BingX",
				Code:    "100001",
				Message: "Invalid API key",
				Err:     ErrAuthFailed,
			},
			target:    ErrRateLimited,
			wantMatch: false,
		},
		{
			name: "Match with standard error",
			err: &BrokerError{
				Broker:  "BingX",
				Code:    "100400",
				Message: "Not enough balance",
				Err:     ErrInsufficientBalance,
			},
			target:    ErrInsufficientBalance,
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errors.Is(tt.err, tt.target)
			if got != tt.wantMatch {
				t.Errorf("errors.Is() = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}
