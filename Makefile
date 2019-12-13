VERSION := $(shell git describe --tags --abbrev=0)

all: test build

build: fmt vet
	CGO_ENABLED=0 go build -a -o build/pet main.go

# Run tests against code
test: fmt vet
	go test ./...

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...
