include .env

ifeq ($(POSTGRES_SETUP_TEST),)
    POSTGRES_SETUP_TEST := user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) host=$(DB_HOST) port=$(DB_PORT) sslmode=disable
endif

PROJECT_ROOT := $(CURDIR)
MIGRATION_FOLDER := $(PROJECT_ROOT)/migrations

.PHONY: migration-create
migration-create:
		goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: migration-up
migration-up:
		goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: migration-down
migration-down:
		goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: test-environment
test-environment: 
		docker-compose up -d
		echo "Waiting for containers to start..." && sleep 3
		make migration-up

.PHONY: test-integrationDB
test-integrationDB: 
		make test-environment 
		go test -tags=integrationDB ./tests
		make clean-test-data
		docker-compose down 

.PHONY: test-integrationHandler
test-integrationHandler: 
		make test-environment 
		go test -tags=integrationHandler ./tests
		make clean-test-data
		docker-compose down

.PHONY: test-unit
test-unit:
		go test ./internal/pkg/server

.PHONY: clean-test-data
clean-test-data:
		goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down-to 0
