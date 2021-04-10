package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jeebster/cryptofolio/asset"
	"golang.org/x/text/message"
)

const MANIFEST_FILENAME = "cryptocurrencies.json"

func parseManifestToCryptocurrencyAssets(filepath string) ([]asset.Cryptocurrency, error) {
	var assetData []asset.Cryptocurrency

	fileContents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return make([]asset.Cryptocurrency, 0), err
	}

	err = json.Unmarshal(fileContents, &assetData)
	if err != nil {
		return make([]asset.Cryptocurrency, 0), err
	}

	return assetData, nil
}

func main() {
	printer := message.NewPrinter(message.MatchLanguage("en"))

	printer.Println("Cryptocurrency Portfolio Balances\n")

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(
			printer.Print("Error obtaining current executable: ", err),
		)
	}

	pwd := filepath.Dir(ex)
	manifestPath := printer.Sprintf("%s/%s", pwd, MANIFEST_FILENAME)

	crossPlatformManifestPath := filepath.FromSlash(manifestPath)

	cryptocurrencies, err := parseManifestToCryptocurrencyAssets(crossPlatformManifestPath)

	if err != nil {
		log.Fatal(
			printer.Printf("Error parsing JSON manifest file: %s", err),
		)
	}

	for i, c := range cryptocurrencies {
		// second element returned from range is a *COPY* of underlying array element
		// https://tour.golang.org/moretypes/7

		printer.Printf("Fetching price data for %s from Coincap.io...\n", c.Id)

		elementRef := &cryptocurrencies[i]
		err := elementRef.FetchData()
		if err != nil {
			log.Fatal(
				printer.Printf("Error fetching data from Coincap.io: %s", err),
			)
		}

		printer.Printf(
			"%f %s at $%f USD = $%f\n\n",
			c.Quantity,
			strings.ToUpper(
				elementRef.Symbol,
			),
			elementRef.Price,
			elementRef.Balance,
		)

		time.Sleep(1 * time.Second)
	}

	printer.Printf("Total portfolio balance: $%f", asset.CalculateTotalAssetsBalance(cryptocurrencies))
}
