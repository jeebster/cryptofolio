# Crediting the following article for Makefile assistance
# http://www.codershaven.com/multi-platform-makefile-for-go/

EXECUTABLE=cryptofolio
WINDOWS=$(EXECUTABLE).windows-amd64.exe
LINUX=$(EXECUTABLE).linux-amd64
DARWIN_AMD=$(EXECUTABLE).mac_os-amd64
DARWIN_ARM=$(EXECUTABLE).mac_os-arm64
VERSION=$(shell git describe --tags --always --long --dirty)

## Platform-specific builds
windows: $(WINDOWS)
linux: $(LINUX)
darwin-amd: $(DARWIN_AMD)
darwin-arm: $(DARWIN_ARM)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -i -v -o ./dist/$(VERSION)/$(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go
	tar -czvf ./dist/$(VERSION)/$(WINDOWS).tar.gz ./dist/$(VERSION)/$(WINDOWS)
	sha256sum ./dist/$(VERSION)/$(WINDOWS) > ./dist/$(VERSION)/sha256sums.txt
	sha256sum ./dist/$(VERSION)/$(WINDOWS).tar.gz > ./dist/$(VERSION)/sha256sums.txt
	rm -rf ./dist/$(VERSION)/$(WINDOWS)

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -i -v -o ./dist/$(VERSION)/$(LINUX) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go
	tar -czvf ./dist/$(VERSION)/$(LINUX).tar.gz ./dist/$(VERSION)/$(LINUX)
	sha256sum ./dist/$(VERSION)/$(LINUX) > ./dist/$(VERSION)/sha256sums.txt
	sha256sum ./dist/$(VERSION)/$(LINUX).tar.gz > ./dist/$(VERSION)/sha256sums.txt
	rm -rf ./dist/$(VERSION)/$(LINUX)

$(DARWIN_AMD):
	env GOOS=darwin GOARCH=amd64 go build -i -v -o ./dist/$(VERSION)/$(DARWIN_AMD) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go
	tar -czvf ./dist/$(VERSION)/$(DARWIN_AMD).tar.gz ./dist/$(VERSION)/$(DARWIN_AMD)
	sha256sum ./dist/$(VERSION)/$(DARWIN_AMD) > ./dist/$(VERSION)/sha256sums.txt
	sha256sum ./dist/$(VERSION)/$(DARWIN_AMD).tar.gz > ./dist/$(VERSION)/sha256sums.txt
	rm -rf ./dist/$(VERSION)/$(DARWIN_AMD)

$(DARWIN_ARM):
	env GOOS=darwin GOARCH=arm64 go build -i -v -o ./dist/$(VERSION)/$(DARWIN_ARM) -ldflags="-s -w -X main.version=$(VERSION)" ./main.go
	tar -czvf ./dist/$(VERSION)/$(DARWIN_ARM).tar.gz ./dist/$(VERSION)/$(DARWIN_ARM)
	sha256sum ./dist/$(VERSION)/$(DARWIN_ARM) > ./dist/$(VERSION)/sha256sums.txt
	sha256sum ./dist/$(VERSION)/$(DARWIN_ARM).tar.gz > ./dist/$(VERSION)/sha256sums.txt
	rm -rf ./dist/$(VERSION)/$(DARWIN_ARM)	

## Checksum file generation
sha256sums: $(SHA256SUMS)

$(SHA256SUMS):
	touch ./dist/$(VERSION)/sha256sums.txt

## Build everything
build: sha256sums windows linux darwin-amd darwin-arm
	@echo version: $(VERSION)

## Display available commands
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
