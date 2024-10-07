# Makefile

# Database connection string
DB_URL := postgres://postgres:ad681789@192.168.0.106:5433/onlinestore?sslmode=disable

# Migration path
MIGRATION_PATH := ./database/migrations

# Command to run migrations
migrate-up:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" down

migrate-reset:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" drop
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up

# Help target to display available commands
help:
	@echo "Available commands:"
	@echo "  make migrate-up      Run database migrations up"
	@echo "  make migrate-down    Run database migrations down"
	@echo "  make migrate-reset    Drop and reapply all migrations"
	@echo "  make help           Show this help message"
