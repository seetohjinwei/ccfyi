MAIN=.
NAME=git
BUILD_DIR=output/

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

## tidy: format code + tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

## test: run all tests
.PHONY: test
test:
	LOG=ERROR go test -v -race -buildvcs ./...

## test: run only the short tests
.PHONY: testshort
testshort:
	LOG=ERROR go test -v -race -buildvcs -short ./...

## build: build the application
.PHONY: build
build:
	go build -o $(BUILD_DIR)$(NAME) ${MAIN}

## run: run the application
.PHONY: run
run: build
	LOG=TRACE $(BUILD_DIR)$(NAME)

## prod: run the application with almost no logs
.PHONY: prod
prod: build
	LOG=PANIC $(BUILD_DIR)$(NAME)

.PHONE: clean
clean:
	rm -rf $(BUILD_DIR)
