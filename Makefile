.PHONY=fmt
BINARY_NAME=water
GOPATH := $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin
INSTALL_DIR=~/.local/bin
BUILD_DIR=./build

GOLANGCILINT := $(GOBIN)/golangci-lint
$(GOLANGCILINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.0

RICHGO := $(GOBIN)/richgo
$(RICHGO):
	@go install github.com/kyoh86/richgo@v0.3.6

fmt:
	@goimports -w .
	@gofmt -w .

lint: $(GOLANGCILINT)
	@golangci-lint run

test: $(RICHGO)
	@$(RICHGO) test -v ./...

check: fmt lint test

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

install:
	test -f $(BUILD_DIR)/$(BINARY_NAME)
	mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)

uninstall:
	rm $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	@rm -rf $(BUILD_DIR)
