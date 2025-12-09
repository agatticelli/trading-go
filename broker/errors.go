package broker

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidSymbol       = errors.New("invalid symbol")
	ErrInvalidPrice        = errors.New("invalid price")
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrPositionNotFound    = errors.New("position not found")
	ErrOrderNotFound       = errors.New("order not found")
	ErrLeverageTooHigh     = errors.New("leverage exceeds maximum")
	ErrAuthFailed          = errors.New("authentication failed")
	ErrRateLimited         = errors.New("rate limited")
	ErrAPIError            = errors.New("API error")
)

// BrokerError wraps exchange-specific errors
type BrokerError struct {
	Broker  string
	Code    string
	Message string
	Err     error
}

func (e *BrokerError) Error() string {
	return fmt.Sprintf("%s error [%s]: %s", e.Broker, e.Code, e.Message)
}

func (e *BrokerError) Unwrap() error {
	return e.Err
}

// NewBrokerError creates a new broker error
func NewBrokerError(broker, code, message string, err error) *BrokerError {
	return &BrokerError{
		Broker:  broker,
		Code:    code,
		Message: message,
		Err:     err,
	}
}
