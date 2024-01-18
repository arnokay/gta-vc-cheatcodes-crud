include .env

start:
	@echo "running api:"
	@CHEATCODES_DB_DSN=$(CHEATCODES_DB_DSN) go run ./cmd/api
migrate-up:
	@echo "migrate up:"
	@migrate -path=./migrations -database=$(CHEATCODES_DB_DSN) up
migrate-force:
	@echo "migrate force:"
	@migrate -path=./migrations -database=$(CHEATCODES_DB_DSN) force $(ARGS)
migrate-goto:
	@echo "migrate goto:"
	@migrate -path=./migrations -database=$(CHEATCODES_DB_DSN) goto $(ARGS)
migrate-down:
	@echo "migrate down:"
	@migrate -path=./migrations -database=$(CHEATCODES_DB_DSN) down
