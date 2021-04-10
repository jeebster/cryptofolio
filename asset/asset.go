package asset

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

const COINCAP_API_BASE_URL string = "https://api.coincap.io/v2/assets"

type Cryptocurrency struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"priceUsd"` // USD
	Balance  float64 `json:"-"`        // USD
}

func (c *Cryptocurrency) FetchData() error {
	response, err := http.Get(COINCAP_API_BASE_URL + "/" + c.Id)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var parsedResponse map[string]interface{}
	err = json.Unmarshal([]byte(body), &parsedResponse)
	if err != nil {
		return err
	}

	data := parsedResponse["data"].(map[string]interface{})

	c.Name = data["name"].(string)
	c.Symbol = data["symbol"].(string)

	price, err := strconv.ParseFloat(
		data["priceUsd"].(string),
		64,
	)

	if err != nil {
		return err
	}

	c.Price = price
	c.Balance = price * c.Quantity

	return nil
}

func CalculateTotalAssetsBalance(assets []Cryptocurrency) float64 {
	totalBalance := 0.0

	for _, c := range assets {
		totalBalance += c.Balance
	}

	return totalBalance
}
