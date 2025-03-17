.PHONY: help run generate-graphql swagger fmt migrate-create migrate-up migrate-down shell

# Variables
API_NAME = msp-api
API_BIN = ./bin/api
MAIN_FILE = ./src/main.go
LOCAL_PORT = 8080
DOCKER_PORT = 3010

# Tools
GO = go

# Colors
COLOR_RESET = \033[0m
COLOR_GREEN = \033[32m
COLOR_YELLOW = \033[33m


help: ## Display available commands
	@echo "$(COLOR_GREEN)Available commands:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(COLOR_YELLOW)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'

ssh-be: ## SSH to the container
	@echo "$(COLOR_GREEN)SSH to the BACKEND...$(COLOR_RESET)"
	docker exec -it makeshop_payment_backend_1 sh

ssh-mysql: ## SSH to the container
	@echo "$(COLOR_GREEN)SSH to the MYSQL...$(COLOR_RESET)"
	docker exec -it makeshop_payment_mysql_1 sh

run: ## Run the API locally
	@echo "$(COLOR_GREEN)Running $(API_NAME) on http://localhost:$(DOCKER_PORT)...$(COLOR_RESET)"
	docker exec -it makeshop_payment_backend_1 go run $(MAIN_FILE)

generate-graphql: ## Generate GraphQL code from schema
	@echo "$(COLOR_GREEN)Generating GraphQL code...$(COLOR_RESET)"
	docker exec -it makeshop_payment_backend_1 go run github.com/99designs/gqlgen generate
	@echo "$(COLOR_GREEN)GraphQL code generated!$(COLOR_RESET)"

swagger: ## Generate Swagger documentation
	@echo "$(COLOR_GREEN)Generating Swagger documentation...$(COLOR_RESET)"
	docker exec -it makeshop_payment_backend_1 swag init -g src/main.go
	@echo "$(COLOR_GREEN)Swagger documentation generated!$(COLOR_RESET)"

fmt: ## Format code
	@echo "$(COLOR_GREEN)Formatting code...$(COLOR_RESET)"
	docker exec -it makeshop_payment_backend_1 $(GO) fmt ./...
	@echo "$(COLOR_GREEN)Formatting complete!$(COLOR_RESET)"

migrate-create: ## Run database migrations create
	@echo "$(COLOR_GREEN)Running create Migrations...$(COLOR_RESET)"
	@read -p "Enter migration name: " name; \
   	docker exec -it makeshop_payment_backend_1 goose -dir ./config/migrations create $$name sql
	@echo "$(COLOR_GREEN)Create Migrations complete!$(COLOR_RESET)"

migrate-up: ## Run database migrations up
	@echo "$(COLOR_GREEN)Running database migrations up...$(COLOR_RESET)"
	docker exec -it makeshop_payment_backend_1 goose -dir ./config/migrations mysql "root:rootpw@tcp(mysql:3306)/msp-db-dev" up
	@echo "$(COLOR_GREEN)Migrations complete!$(COLOR_RESET)"

migrate-down: ## Run database migrations down
	@echo "$(COLOR_GREEN)Running database migrations down...$(COLOR_RESET)"
	docker exec -it makeshop_payment_backend_1 goose -dir ./config/migrations mysql "root:rootpw@tcp(mysql:3306)/msp-db-dev" down
	@echo "$(COLOR_GREEN)Migrations complete!$(COLOR_RESET)"

shell: ## Run shell in the container
	@$(eval ARGS := $(filter-out $@,$(MAKECMDGOALS)))
	@echo "$(COLOR_GREEN)Running Shell...$(COLOR_RESET)"
	docker exec -it makeshop_payment_backend_1 go run main.go $(ARGS)
	@echo "$(COLOR_GREEN)Running Shell complete!$(COLOR_RESET)"