# Cryptofolio

A program to display your cryptocurrency balances in fiat currency - best for HODLers

## Requirements
If building from source, ensure [Golang](https://golang.org/dl/) is installed on your system

## Caveats
Currently...

This program leverages the [Coincap API](https://docs.coincap.io/#ee30bea9-bb6b-469d-958a-d3e35d442d7a) for markets data reference.

Only USD fiat currency is supported.

All input data is sourced from a user supplied manifest file in JSON format (e.g. `YOURFILENAME.json`) via the `-manifest` argument.

## Configuration
Edit your JSON manifest file with your asset data. The manifest must include an array of objects. Object properties include:

* `id (string, required)`: [Coincap API](https://docs.coincap.io/#ee30bea9-bb6b-469d-958a-d3e35d442d7a) asset identifier
* `quantity (float, required)`: quantity of the asset you currently hold

Example: `manifest.json`

```json
[
  {
    "id": "bitcoin",
    "quantity": 10.0
  },
  {
    "id": "bitcoin-cash",
    "quantity": 10.0
  },  
  {
    "id": "ethereum",
    "quantity": 10.0
  },
  {
    "id": "cardano",
    "quantity": 10.0
  },  
  {
    "id": "xrp",
    "quantity": 10.0
  }
]
```

## Usage
You can either compile the program from source, or download a pre-compiled binary specific to your operating system 

### Build
From the repository root directory, run the following shell command:

```go
go build
```

### Execute
Invoke the binary relative to your operating system via shell:

### Linux/MacOS
```sh
./cryptofolio -manifest=/path/to/your/manifest.json
```

### Windows
```sh
.\cryptofolio.exe -manifest=C:\path\to\your\manifest.json
```