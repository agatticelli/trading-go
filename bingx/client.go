package bingx

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/agatticelli/trading-go/broker"
)

// Client implements broker.Broker interface for BingX
type Client struct {
	apiKey     string
	secretKey  string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new BingX broker client
func NewClient(apiKey, secretKey string, demoMode bool) *Client {
	baseURL := BaseURLProd
	if demoMode {
		baseURL = BaseURLDemo
	}

	return &Client{
		apiKey:    apiKey,
		secretKey: secretKey,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the broker name
func (c *Client) Name() string {
	return "bingx"
}

// SupportedFeatures returns the features supported by BingX
func (c *Client) SupportedFeatures() broker.Features {
	return broker.Features{
		TrailingStop:     true,
		MultipleTP:       true,
		BracketOrders:    true,
		MaxLeverage:      125,
		ReduceOnlyOrders: true,
	}
}

// sign creates HMAC-SHA256 signature for API requests
func (c *Client) sign(params string) string {
	h := hmac.New(sha256.New, []byte(c.secretKey))
	h.Write([]byte(params))
	signature := hex.EncodeToString(h.Sum(nil))
	return signature
}

// makeRequest makes an HTTP request to BingX API
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, params map[string]string) ([]byte, error) {
	timestamp := time.Now().UnixMilli()

	// Add timestamp to parameters
	if params == nil {
		params = make(map[string]string)
	}
	params["timestamp"] = strconv.FormatInt(timestamp, 10)

	// Build query parameters (sorted keys)
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	queryString := values.Encode()

	// Create signature
	signature := c.sign(queryString)

	// Add signature to URL
	fullURL := fmt.Sprintf("%s%s?%s&signature=%s", c.baseURL, endpoint, queryString, signature)

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("X-BX-APIKEY", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, broker.NewBrokerError("bingx", "REQUEST_FAILED", "HTTP request failed", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, broker.NewBrokerError("bingx", "READ_FAILED", "Failed to read response", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, broker.NewBrokerError("bingx", "HTTP_ERROR",
			fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)), nil)
	}

	return body, nil
}
