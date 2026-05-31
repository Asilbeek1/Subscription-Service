-include .env
export 
MIGRATION_DSN=postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable


show-dsn:
	@echo "MIGRATION_DSN=$(MIGRATION_DSN)"

migrate-up:
	migrate -path migrations -database "$(MIGRATION_DSN)" up

migrate-down:
	migrate -path migrations -database "$(MIGRATION_DSN)" down 1

