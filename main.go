package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/jeebster/cryptofolio/asset"
	"golang.org/x/text/message"
)

const PROGRAM_TITLE = "Cryptocurrency Portfolio Balances"

func buildQueryString(cryptocurrencies []asset.Cryptocurrency) string {
	queryParamsBuilder := strings.Builder{}
	queryParamsBuilder.WriteString("?ids=")
	for _, c := range cryptocurrencies {
		queryParamsBuilder.WriteString(c.Id + ",")
	}

	return queryParamsBuilder.String()
}

func parseManifest(filepath string) ([]asset.Cryptocurrency, error) {
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

	manifestDataCryptocurrencies, err := parseManifest(*manifestPtr)

	if err != nil {
		log.Fatal(
			printer.Printf("Error parsing JSON manifest file: %s", err),
		)
	}

	printer.Println("Fetching price data from coincap.io...\n")

	queryParams := buildQueryString(manifestDataCryptocurrencies)
	apiDataCryptocurrencies, err := asset.FetchAssets(queryParams)
	if err != nil {
		log.Fatal(
			printer.Printf("Error fetching data from coincap.io: %s", err),
		)
	}

	asset.MergeAssetsData(manifestDataCryptocurrencies, apiDataCryptocurrencies)
	if err != nil {
		log.Fatal(
			printer.Printf("Error parsing data from coincap.io: %s", err),
		)
	}

	for i, cc := range manifestDataCryptocurrencies {
		cryptoRef := &manifestDataCryptocurrencies[i]
		err := cryptoRef.SetBalance()
		if err != nil {
			log.Fatal(
				printer.Printf(
					"Error setting %s balance: %s",
					cc.Symbol,
					err,
				),
			)
		}

		priceAsFloat, err := cc.PriceAsFloat()
		if err != nil {
			log.Fatal(
				printer.Printf(
					"Error converting %s price string to float: %s",
					cc.Symbol,
					err,
				),
			)
		}

		printer.Printf(
			"%f %s at $%f USD = $%f\n",
			cc.Quantity,
			strings.ToUpper(
				cc.Symbol,
			),
			priceAsFloat,
			cryptoRef.GetBalance(),
		)
	}

	printer.Println()
	printer.Printf("Total portfolio balance: $%f\n", asset.CalculateTotalAssetsBalance(manifestDataCryptocurrencies))
}
