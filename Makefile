include .env

start:
	@echo "running api..."
	@CHEATCODES_DB_DSN=$(CHEATCODES_DB_DSN) go run ./cmd/api
