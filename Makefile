default: build

SETENV=
ifeq ($(OS),Windows_NT)
	SETENV=set
endif

ifeq ($(OS),Windows_NT)
	BINARY_NAME=btptf.exe
	GOBIN_PATH=$(if $(GOBIN),$(GOBIN),$(shell powershell -Command go env GOPATH)\bin)
	BINARY_PATH=$(GOBIN_PATH)\$(BINARY_NAME)
else
	BINARY_NAME=btptf
	GOBIN_PATH=$(if $(GOBIN),$(GOBIN),$(shell go env GOPATH)/bin)
	BINARY_PATH=$(GOBIN_PATH)/$(BINARY_NAME)
endif

MAIN_PACKAGE=main.go

build:
	go build -v ./...

install: build
	go build -o $(BINARY_PATH) $(MAIN_PACKAGE)

lint:
	golangci-lint run

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -tags=all -timeout=900s -parallel=4 ./...

docs:
	go run main.go gendoc -s "abc"

.PHONY: build install lint fmt test docs
