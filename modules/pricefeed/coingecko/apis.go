package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/forbole/bdjuno/v4/types"
)

const (
	coingeckoRequestTimeout  = 10 * time.Second
	coingeckoMaxResponseSize = 2 << 20 // 2 MiB
)

var coingeckoHTTPClient = &http.Client{
	Timeout: coingeckoRequestTimeout,
}

// GetCoinsList allows to fetch from the remote APIs the list of all the supported tokens
func GetCoinsList() (coins Tokens, err error) {
	err = queryCoinGecko("/coins/list", &coins)
	return coins, err
}

// GetTokensPrices queries the remote APIs to get the token prices of all the tokens having the given ids
func GetTokensPrices(ids []string) ([]types.TokenPrice, error) {
	var prices []MarketTicker
	query := fmt.Sprintf("/coins/markets?vs_currency=usd&ids=%s", strings.Join(ids, ","))
	err := queryCoinGecko(query, &prices)
	if err != nil {
		return nil, err
	}

	return ConvertCoingeckoPrices(prices), nil
}

func ConvertCoingeckoPrices(prices []MarketTicker) []types.TokenPrice {
	tokenPrices := make([]types.TokenPrice, len(prices))
	for i, price := range prices {
		tokenPrices[i] = types.NewTokenPrice(
			price.Symbol,
			price.CurrentPrice,
			int64(math.Trunc(price.MarketCap)),
			price.LastUpdated,
		)
	}
	return tokenPrices
}

// queryCoinGecko queries the CoinGecko APIs for the given endpoint
func queryCoinGecko(endpoint string, ptr interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), coingeckoRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.coingecko.com/api/v3"+endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := coingeckoHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("coingecko returned non-success status: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(io.LimitReader(resp.Body, coingeckoMaxResponseSize))
	err = decoder.Decode(ptr)
	if err != nil {
		return fmt.Errorf("error while unmarshaling response body: %s", err)
	}

	return nil
}
