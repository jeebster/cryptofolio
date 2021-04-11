package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jeebster/cryptofolio/asset"
	"golang.org/x/text/message"
)

const MANIFEST_FILENAME = "cryptocurrencies.json"

func fetchDataWorker(wg *sync.WaitGroup, c *asset.Cryptocurrency, printer *message.Printer) {
	defer wg.Done()

	err := c.FetchData()
	if err != nil {
		log.Fatal(
			printer.Printf("Error fetching data from Coincap.io: %s", err),
		)
	}

	printer.Printf(
		"%f %s at $%f USD = $%f\n",
		c.Quantity,
		strings.ToUpper(
			c.Symbol,
		),
		c.Price,
		c.Balance,
	)
}

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

	var wg sync.WaitGroup

	printer.Printf("Fetching price data from Coincap.io...\n\n")

	for i, _ := range cryptocurrencies {
		// second element returned from range is a *COPY* of underlying array element
		// https://tour.golang.org/moretypes/7

		wg.Add(1)
		go fetchDataWorker(&wg, &cryptocurrencies[i], printer)
	}

	wg.Wait()
	printer.Printf("\nTotal portfolio balance: $%f", asset.CalculateTotalAssetsBalance(cryptocurrencies))
}
