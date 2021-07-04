-include .env


VERSION := $(shell git describe --tags)
PROJECTNAME := $(shell basename "$(PWD)")
LDFLAGS=-ldflags "-X=main.Version=$(VERSION)"


## build: Complie Golang files
build:
	@echo "  >  Building binary..."
	go build $(LDFLAGS) ./cmd/replica

## buildimage: Build docker image
buildimage:
	@echo "  >  Building docker image..."
	docker build -t $(PROJECTNAME) .

## test: Run test with verbose
test:
	go test -v


.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
