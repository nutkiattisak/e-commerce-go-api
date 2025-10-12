.PHONY: run dev tidy

SWAG_BIN := $(shell go env GOPATH)/bin/swag

swagger:
	@echo "generate swagger docs..."
	@echo "using $(SWAG_BIN)"
	@$(SWAG_BIN) fmt
	@$(SWAG_BIN) init -g main.go -o docs

run:
	go run .

dev:
	air || go run .

tidy:
	go mod tidy
