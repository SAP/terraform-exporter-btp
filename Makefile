default: build

SETENV=
ifeq ($(OS),Windows_NT)
	SETENV=set
endif

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -tags=all -timeout=900s -parallel=4 ./...

docs:
	go run main.go gendoc

.PHONY: build install lint fmt test docs
