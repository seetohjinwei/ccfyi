MAIN=.
NAME=redis
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
	go test -v -race -buildvcs ./...

## build: build the application
.PHONY: build
build:
	go build -o $(BUILD_DIR)$(NAME) ${MAIN}

## run: run the application
.PHONY: run
run: build
	$(BUILD_DIR)$(NAME)

.PHONE: clean
clean:
	rm -rf $(BUILD_DIR)
