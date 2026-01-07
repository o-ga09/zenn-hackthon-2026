.PHONY: lint
lint:
	go vet ./...
	golangci-lint run

.PHONY: test
test:
	go test -race -parallel 1 ./...

.PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out

.PHONY: generate
generate:
	go generate ./...

# DATABASE_URL is expected to be set in the environment
ifndef DATABASE_URL
$(warning DATABASE_URL is not set. Using default local development settings)
DATABASE_URL ?= user:P@ssw0rd@tcp(127.0.0.1:3306)/develop_tavinikkiy?parseTime=true
export DATABASE_URL
endif

.PHONY: migrate
migrate-up:
	$(eval MIGRATE_CMD := go run cmd/migration/main.go)
	$(MIGRATE_CMD) -command up

.PHONY: create-migrate
create-migrate:
	$(eval MIGRATE_CMD := go run cmd/migration/main.go)
	$(MIGRATE_CMD) -command new -name $(name)

.PHONY: seed
seed:
	$(eval MIGRATE_CMD := go run cmd/migration/main.go)
	$(MIGRATE_CMD) -command seed

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down -v

.PHONY: build
compose-build:
	docker compose up -d --build

.PHONY: check
check: lint test test-coverage migrate-up migrate-status migrate-down seed run build compose-up compose-down compose-build compose-start encode-image compose-stop
