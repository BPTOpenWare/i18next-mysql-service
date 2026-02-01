.DEFAULT_GOAL := build

PROTO_DIR := proto
PROTO_SRC := $(PROTO_DIR)/resources/resources.proto
OPENAPI_DIR := Docs
BIN_DIR := bin
BINARY := i18next-mysql-service

# Where Go installs tools
GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

.PHONY: help tools proto build clean

help:
	@echo "Available targets:"
	@echo "  tools   - Install protoc plugins (go, go-grpc, grpc-gateway, openapi)"
	@echo "  proto   - Generate Go stubs, gRPC-Gateway, and OpenAPI"
	@echo "  build   - Generate protos then build the Go binary"
	@echo "  clean   - Remove build artifacts"

tools:
	@echo "Installing protoc plugins to $(GOBIN)"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapi@latest

proto:
	@command -v protoc >/dev/null 2>&1 || { echo "Error: protoc is not installed."; exit 1; }
	@mkdir -p $(OPENAPI_DIR)
	protoc -I $(PROTO_DIR) \
	  --go_out $(PROTO_DIR) --go_opt paths=source_relative \
	  --go-grpc_out $(PROTO_DIR) --go-grpc_opt paths=source_relative \
	  --grpc-gateway_out $(PROTO_DIR) --grpc-gateway_opt paths=source_relative \
	  --openapi_out=$(OPENAPI_DIR) \
	  $(PROTO_SRC)

build: proto
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) .

clean:
	@rm -rf $(BIN_DIR)
	@rm -rf $(OPENAPI_DIR)
	@find $(PROTO_DIR) -name "*.pb.go" -type f -delete
	@find $(PROTO_DIR) -name "*.pb.gw.go" -type f -delete
