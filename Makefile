include .env
export

.PHONY: gen gen-code gen-jet migrate

gen:
	go generate ./...

gen-code:
	go generate ./core/code/

gen-jet:
	go run github.com/go-jet/jet/v2/cmd/jet \
		-dsn="postgresql://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable" \
		-path=./database

migrate:
	atlas schema apply \
		--url "postgresql://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable" \
		--to "file://database/schema/" \
		--dev-url "postgresql://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(TEMP_DATABASE_NAME)?sslmode=disable" \
