WIRE_DIR := ./internal/di

wire:
	cd $(WIRE_DIR) && wire

test:
	go test -v ./...

mockery:
	mockery