package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/jeebster/cryptofolio/asset"
	"golang.org/x/text/message"
)

const PROGRAM_TITLE = "Cryptocurrency Portfolio Balances"

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

	fileContents, err := os.ReadFile(filepath)
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
	var pathEx string
	os := runtime.GOOS
	switch os {
	case "windows":
		pathEx = `C:\Users\guest\documents\manifest.json`
	case "darwin":
		pathEx = "/Users/guest/documents/manifest.json"
	case "linux":
		pathEx = "/home/guest/documents/manifest.json"
	}

	manifestFlageDescription := fmt.Sprintf("absolute path to your manifest JSON file, e.g. %s", pathEx)

	manifestPtr := flag.String("manifest", "", manifestFlageDescription)
	flag.Parse()

	if len(*manifestPtr) == 0 {
		log.Fatal("You must supply a manifest JSON file via the -manifest flag. Run the program with the -h flag for more information.")
	}

	printer := message.NewPrinter(message.MatchLanguage("en"))

	programSubtitleBuilder := strings.Builder{}
	for i := 0; i < len(PROGRAM_TITLE)-1; i++ {
		programSubtitleBuilder.WriteRune('_')
	}

	printer.Println(PROGRAM_TITLE)
	printer.Println(programSubtitleBuilder.String(), "\n")

	cryptocurrencies, err := parseManifestToCryptocurrencyAssets(*manifestPtr)

	if err != nil {
		log.Fatal(
			printer.Printf("Error parsing JSON manifest file: %s", err),
		)
	}

	var wg sync.WaitGroup

	printer.Println("Fetching price data from coincap.io...\n")

	for i, _ := range cryptocurrencies {
		// second element returned from range is a *COPY* of underlying array element
		// https://tour.golang.org/moretypes/7

		wg.Add(1)
		go fetchDataWorker(&wg, &cryptocurrencies[i], printer)
	}

	wg.Wait()
	printer.Println()
	printer.Printf("Total portfolio balance: $%f\n", asset.CalculateTotalAssetsBalance(cryptocurrencies))
}
