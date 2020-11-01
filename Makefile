override RELEASE="$(git tag -l --points-at HEAD)"
override COMMIT="$(shell git rev-parse --short HEAD)"
override BUILD_TIME="$(shell date -u '+%Y-%m-%dT%H:%M:%S')"

all:

build:
	go build -ldflags="-s -w" -v -o dist/artisync-api ./cmd/artisync-api/...
	go build -ldflags="-s -w" -v -o dist/artisync-daily ./cmd/artisync-daily/...
	go build -ldflags="-s -w" -v -o dist/artisync-sync ./cmd/artisync-sync/...

test t:
	go test -v ./internal/...

lint l:
	bash ./scripts/revive.sh
	bash ./scripts/golangci-lint.sh

run-api:
	go run -v ./cmd/artisync-api/... --config ./artisync.example.yml

run-daily:
	go run -v ./cmd/artisync-daily/... --config ./artisync.example.yml

run-sync:
	go run -v ./cmd/artisync-sync/... --config ./artisync.example.yml

compose:
	docker-compose up -d --build

exec-sources:
	docker exec -it artisync.sources bash

image:
	docker build --file ./Dockerfile             --tag musicmash/artisync-builder:latest .
	docker build --file ./build/api/Dockerfile   --tag musicmash/artisync-api:latest .
	docker build --file ./build/daily/Dockerfile --tag musicmash/artisync-daily:latest .
	docker build --file ./build/sync/Dockerfile  --tag musicmash/artisync-sync:latest .

ensure-go-migrate-installed:
	bash ./scripts/install-go-migrate.sh

db-generate:
	sqlc generate

# show latest applied migration
db-status: ensure-go-migrate-installed
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose version

# apply migration up
db-up: ensure-go-migrate-installed
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

# apply migration down
db-down: ensure-go-migrate-installed
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down
