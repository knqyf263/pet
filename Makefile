.PHONY: \
	dep \
	build \
	vet \
	test

dep:
	go mod download

build: main.go
	go build -o pet $<

test:
	go test ./...

vet:
	go vet