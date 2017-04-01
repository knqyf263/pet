.PHONY: \
	all \
	glide \
	deps \
	update \
	build \
	install \
	lint \
	vet \
	fmt \
	clean

VERSION := $(shell git describe --tags --abbrev=0)

all: glide deps build 

glide:
	go get github.com/Masterminds/glide

deps: glide
	glide install

update: glide
	glide update

build: main.go deps
	go build -o pet $<

install: main.go deps
	go install


lint:
	@ go get -v github.com/golang/lint/golint
	golint $(shell glide nv)

vet:
	go vet $(shell glide nv)

fmt:
	go fmt $(shell glide nv)

clean:
	go clean $(shell glide nv)
