package asset

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

const COINCAP_API_BASE_URL string = "https://api.coincap.io/v2/assets"

type APIResponse struct {
	Data []Cryptocurrency `json:"data"`
}

type Cryptocurrency struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	PriceUsd string  `json:"priceUsd"` // USD
	balance  float64 `json:"-"`        // USD
}

func (c *Cryptocurrency) GetBalance() float64 {
	return c.balance
}

func (c *Cryptocurrency) SetBalance() error {
	priceAsFloat, err := c.PriceAsFloat()
	if err != nil {
		return err
	}

	c.balance = priceAsFloat * c.Quantity

	return nil
}

func (c *Cryptocurrency) PriceAsFloat() (float64, error) {
	result, err := strconv.ParseFloat(c.PriceUsd, 64)
	if err != nil {
		return 0.0, err
	}

	return result, nil
}

func FetchAssets(queryString string) ([]Cryptocurrency, error) {
	response, err := http.Get(COINCAP_API_BASE_URL + queryString)

	if err != nil {
		return []Cryptocurrency{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []Cryptocurrency{}, err
	}

	var parsedResponse APIResponse

	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return []Cryptocurrency{}, err
	}

	return parsedResponse.Data, nil
}

func getCryptocurrencyById(id string, cryptocurrencies []Cryptocurrency) *Cryptocurrency {
	var cryptoRef *Cryptocurrency
	for i, cc := range cryptocurrencies {
		if cc.Id == id {
			cryptoRef = &cryptocurrencies[i]
		}
	}

	return cryptoRef
}

func MergeAssetsData(manifestData []Cryptocurrency, apiData []Cryptocurrency) {
	for i, cc := range manifestData {
		targetCrypto := &manifestData[i]
		sourceCrypto := getCryptocurrencyById(cc.Id, apiData)

		targetCrypto.Name = sourceCrypto.Name
		targetCrypto.Symbol = sourceCrypto.Symbol
		targetCrypto.PriceUsd = sourceCrypto.PriceUsd
	}
}

func CalculateTotalAssetsBalance(assets []Cryptocurrency) float64 {
	totalBalance := 0.0

	for _, c := range assets {
		totalBalance += c.GetBalance()
	}

	return totalBalance
}
