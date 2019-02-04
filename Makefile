.PHONY: \
	all \
	dep \
	depup \
	update \
	build \
	install \
	lint \
	vet \
	fmt \
	fmtcheck \
	clean \
	pretest \
	test

VERSION := $(shell git describe --tags --abbrev=0)
SRCS = $(shell git ls-files '*.go')
PKGS = $(shell go list ./... | grep -v /vendor/)

all: dep build test

dep:
	go get -u github.com/golang/dep/...
	dep ensure -vendor-only

depup:
	go get -u github.com/golang/dep/...
	dep ensure -u

build: main.go dep
	go build -o pet $<

install: main.go dep
	go install

lint:
	@ go get -v github.com/golang/lint/golint
	$(foreach file,$(SRCS),golint $(file) || exit;)

vet:
	go vet $(PKGS) || exit;

fmt:
	gofmt -w $(SRCS)

fmtcheck:
	@ $(foreach file,$(SRCS),gofmt -s -l $(file);)

clean:
	go clean $(shell glide nv)

pretest: vet fmtcheck

test: pretest
	go install
	@ $(foreach pkg,$(PKGS), go test $(pkg) || exit;)
