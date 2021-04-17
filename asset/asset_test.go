package asset

import (
	"testing"
)

func generateBitcoinFactory() Cryptocurrency {
	return Cryptocurrency{
		Id:       "bitcoin",
		Name:     "Bitcoin",
		Symbol:   "btc",
		Quantity: 10.00,
		PriceUsd: "60000.00",
	}
}

func TestMergeAssetsData(t *testing.T) {
	bitcoinFactory := generateBitcoinFactory()

	manifestData := []Cryptocurrency{
		{
			Id:       "bitcoin",
			Quantity: 10,
		},
	}

	apiData := []Cryptocurrency{
		bitcoinFactory,
	}

	MergeAssetsData(manifestData, apiData)

	sourceCrypto := apiData[0]
	targetCrypto := manifestData[0]

	if targetCrypto.Name != sourceCrypto.Name ||
		targetCrypto.Symbol != sourceCrypto.Symbol ||
		targetCrypto.PriceUsd != sourceCrypto.PriceUsd {
		t.Errorf(
			"Merge target data does not match source data. Target %v, Source %v",
			targetCrypto,
			sourceCrypto,
		)
	}
}

func TestPriceAsFloat(t *testing.T) {
	cc := generateBitcoinFactory()

	priceAsFloat, err := cc.PriceAsFloat()
	if err != nil {
		t.Errorf("Error calculating PriceAsFloat() for %s", cc.Name)
	}

	var iface interface{} = priceAsFloat
	_, isFloat64 := iface.(float64)

	if !isFloat64 {
		t.Errorf("Error converting PriceUsd %s to float", cc.PriceUsd)
	}
}

func TestSetBalance(t *testing.T) {
	cc := generateBitcoinFactory()

	initialBalance := cc.GetBalance()

	if initialBalance != 0 {
		t.Errorf("Balance should not be initialized, value is %f", initialBalance)
	}

	cc.SetBalance()

	if cc.GetBalance() == 0 {
		t.Errorf("Balance should not be zero")
	}
}

func TestCalculateTotalAssetsBalance(t *testing.T) {
	assets := []Cryptocurrency{
		generateBitcoinFactory(),
		{
			Id:       "ethereum",
			Name:     "Ethereum",
			Symbol:   "eth",
			Quantity: 10.00,
			PriceUsd: "2500.00",
		},
		{
			Id:       "cardano",
			Name:     "Cardano",
			Symbol:   "ada",
			Quantity: 10.00,
			PriceUsd: "1.50",
		},
	}

	for i, _ := range assets {
		ccRef := &assets[i]

		ccRef.SetBalance()

		priceAsFloat, err := ccRef.PriceAsFloat()
		if err != nil {
			t.Errorf("Error calculating PriceAsFloat() for %s", ccRef.Name)
		}

		balanceCalculation := priceAsFloat * ccRef.Quantity

		if ccRef.GetBalance() != balanceCalculation {
			t.Errorf(
				"balanceCalculation %f does not equal %s GetBalance() %f",
				balanceCalculation,
				ccRef.Name,
				ccRef.GetBalance(),
			)
		}
	}
}
