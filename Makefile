# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

CURRENT_DIR=$(shell pwd)
LGOPATH=$(shell echo ${CURRENT_DIR}/../..)

BINARY_NAME=douban_spider

release_windows: clean release_build_windows
release_macos: clean release_build_macos
release_linux: clean release_build_linux

release_build_linux:
	export GOPATH=$(LGOPATH);$(GOBUILD) -o $(CURRENT_DIR)/bin/$(BINARY_NAME) -v main.go

release_build_macos:
	export GOPATH=$(LGOPATH);$(GOBUILD) -o $(CURRENT_DIR)/bin/$(BINARY_NAME) -v main.go

release_build_windows:
	export GOPATH=$(LGOPATH);$(GOBUILD) -o $(CURRENT_DIR)/bin/$(BINARY_NAME).exe -v main.go

clean:
	$(GOCLEAN)