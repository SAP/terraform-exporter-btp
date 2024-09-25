default: build

SETENV=
ifeq ($(OS),Windows_NT)
	SETENV=set
endif

BINARY_NAME=btptf
MAIN_PACKAGE=main.go
GOBIN_PATH=$(if $(GOBIN),$(GOBIN),$(shell go env GOPATH)/bin)

build:
	go build -v ./...

install: build
	go build -o $(GOBIN_PATH)/$(BINARY_NAME) $(MAIN_PACKAGE)

lint:
	golangci-lint run

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -tags=all -timeout=900s -parallel=4 ./...

docs:
	go run main.go gendoc -s "abc"

.PHONY: build install lint fmt test docs
