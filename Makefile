CLIENT_VERSION = v1.0.0
DOC_PORT = 6060

SERVER_BINARY_NAME = server
SERVER_PACKAGE_PATH = ./cmd/server

CLIENT_BINARY_NAME = client
CLIENT_PACKAGE_PATH = ./cmd/client

# OS detection
ifeq ($(OS), Windows_NT)
	DATE = $(shell date /t)
else
	DATE = $(shell date +'%d/%m/%Y')
endif

CLIENT_LDFLAGS = "-X 'main.buildVersion=v1.0.0' -X 'main.buildDate=$(DATE)'"

# ====================
# HELPERS
# ====================

## help: show this help message
help:
	@echo 'usage: make <target> ...'
	@echo ''
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ====================
# QUALITY
# ====================

## tidy: format code and tidy mod file
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

# ====================
# DEVELOPMENT
# ====================

## test: run all tests
test:
	go test ./...

## test/cover: run all tests and display coverage
test/cover:
	go test ./... -coverprofile=/tmp/coverage.out
	go tool cover -html=/tmp/coverage.out

## server/build: build the server
server/build:
	go build -o /tmp/bin/$(SERVER_BINARY_NAME) $(SERVER_PACKAGE_PATH)

## server/run: run the server
server/run: server/build
	/tmp/bin/$(SERVER_BINARY_NAME)

## client/build: build the client
client/build:
	go build -ldflags $(CLIENT_LDFLAGS) -o /tmp/bin/$(CLIENT_BINARY_NAME) $(CLIENT_PACKAGE_PATH)

## client/run: run the client
client/run: client/build
	/tmp/bin/$(CLIENT_BINARY_NAME)

# ====================
# DOCUMENTATION
# ====================

## doc: generate documentation on http port
doc:
	@echo 'open http://localhost:$(DOC_PORT)/pkg/github.com/pavlegich/gophkeeper/?m=all'
	godoc -http=:$(DOC_PORT)

.PHONY: help tidy test test/cover doc