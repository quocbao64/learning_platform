WIRE_DIR := ./internal/di

wire:
	cd $(WIRE_DIR) && wire

test:
	go test -v ./...

mockery:
	mockery

swag:
	swag init -g cmd/api/main.go -o docs