package bingx

const (
	// Base URLs
	BaseURLProd = "https://open-api.bingx.com"
	BaseURLDemo = "https://open-api-vst.bingx.com"

	// BingX API endpoints
	EndpointBalance    = "/openApi/swap/v3/user/balance"
	EndpointPositions  = "/openApi/swap/v2/user/positions"
	EndpointPlaceOrder = "/openApi/swap/v2/trade/order"
	EndpointOpenOrders = "/openApi/swap/v2/trade/openOrders"
	EndpointCancelAll  = "/openApi/swap/v2/trade/allOpenOrders"
	EndpointLeverage   = "/openApi/swap/v2/trade/leverage"
	EndpointServerTime = "/openApi/swap/v2/server/time"
	EndpointPrice      = "/openApi/swap/v1/ticker/price"

	// API response codes
	APISuccessCode = 0
)
