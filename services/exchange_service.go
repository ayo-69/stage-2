package services

import (
	"encoding/json"
	"net/http"
	"time"
)

type ExchangeResponse struct {
	Result string             `json:"result"`
	Rates  map[string]float64 `json:"rates"`
}

var exchageCache = struct {
	Rates     map[string]float64
	Timestamp time.Time
}{}

func GetExchangeRates() (map[string]float64, error) {
	if time.Since(exchageCache.Timestamp) < time.Hour && exchageCache.Rates != nil {
		return exchageCache.Rates, nil
	}

	resp, err := http.Get("https://open.er-api.com/v6/latest/USD")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if data.Result != "success" {
		return nil, http.ErrBodyNotAllowed
	}

	exchageCache.Rates = data.Rates
	exchageCache.Timestamp = time.Now()
	return data.Rates, nil
}
