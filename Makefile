PROTO_DIR := proto
PROTO_GEN_DIR := pb

.PHONY: all
all: proto-clean proto-gen deps build

.PHONY: proto-gen
proto-gen:
	@echo "Generating protocol buffer files..."
	mkdir -p $(PROTO_GEN_DIR)
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/ratelimit/v1/ratelimit.proto
	@echo "Proto files generated successfully in $(PROTO_GEN_DIR)/"

proto-clean:
	@echo "Cleaning generated protocol buffer files..."
	rm -rf $(PROTO_GEN_DIR)/*
	@echo "Cleaned generated proto files from $(PROTO_GEN_DIR)/"

.PHONY: build
build:
	@echo "Building rate limiter server..."
	go build -ldflags="-s -w" -o build/ ./cmd/rate-limiter
		
.PHONY: run
run:
	@echo "Starting rate limiter server..."
	go run ./cmd/rate-limiter/main.go

.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race -timeout 10s ./...

.PHONY: lint
lint:
	@echo "Linting..."
	go vet ./...

.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

.DEFAULT_GOAL := build

