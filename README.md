# Cryptofolio

A program to display your cryptocurrency balances in fiat currency - best for HODLers

## Requirements
Ensure [Golang](https://golang.org/dl/) is installed on your system

## Caveats
Currently...

This program leverages the [Coincap API](https://docs.coincap.io/#ee30bea9-bb6b-469d-958a-d3e35d442d7a) for markets data reference

Only USD fiat currency is supported

All input data is sourced from an included manifest template file, `cryptocurrencies.json`. The top 5 assets by market cap are included for your convenience.

## Configuration
Edit the manifest template file `cryptocurrencies.json` with your asset data. Properties include:

* `id (string, required)`: [Coincap API](https://docs.coincap.io/#ee30bea9-bb6b-469d-958a-d3e35d442d7a) asset identifier
* `quantity (float, required)`: quantity of the asset you currently hold

## Build
From the repository root directory, run the following shell command:

```go
go build
```

## Execute
Invoke the executable relative to your operating system via shell:

### Linux/MacOS
```sh
./cryptofolio
```

### Windows
```sh
.\cryptofolio.exe
```