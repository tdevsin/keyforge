# Variables with folder names and binary names
BINARY_NAME=keyforge
OUTPUT_BINARY=bin/$(BINARY_NAME)
PROTO_DIR=proto
PROTO_OUT=internal/proto

# Default target
# .PHONY is used to make sure that the target is not a file
.PHONY: all
all: build

# Build target
.PHONY: build
build:
	go build -o $(OUTPUT_BINARY) main.go

# Run targets

# Run the server using the CLI
.PHONY: run-server
run-server:
	$(OUTPUT_BINARY) start

# Run the CLI with custom arguments
.PHONY: run-cli
run-cli:
	$(OUTPUT_BINARY) $(ARGS)

# Generate Go code from proto files
.PHONY: proto
proto:
	protoc -I $(PROTO_DIR) $(PROTO_DIR)/*.proto --go_out=$(PROTO_OUT) --go_opt=paths=source_relative --go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative

# Initialize Cobra CLI
.PHONY: cobra-init
cobra-init:
	cobra-cli init

# Generate the CLI commands via Cobra CLI generator
.PHONY: gen-cmd
gen-cmd:
	cobra-cli add $(CMD)

# Test
.PHONY: test
test:
	go test ./... -v

# Clean the existing binaries
.PHONY: clean
clean:
	rm -rf bin
