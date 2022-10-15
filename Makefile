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

GOSEC := $(GOBIN)/gosec
$(GOSEC):
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.13.1

fmt:
	@goimports -w .
	@gofmt -w .

lint: $(GOLANGCILINT)
	@golangci-lint run

test: $(RICHGO)
	@$(RICHGO) test -v ./...

security: $(GOSEC)
	@gosec -quiet ./...

check: fmt lint test security

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
