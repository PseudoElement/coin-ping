package api

type CoinMarketCapErrorResp struct {
	Status struct {
		Timestamp    string `json:"timestamp"`
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
		Elapsed      int    `json:"elapsed"`
		CreditCount  int    `json:"credit_count"`
	} `json:"status"`
}

type CoinMarketCapCoinQuotesLatestReq struct {
	Symbol string `json:"symbol"`
}

type CoinMarketCapCoinQuotesLatestResp struct {
	Data map[string][]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Quote  map[string]struct {
			Price float64 `json:"price"`
		} `json:"quote"`
	} `json:"data"`
}
