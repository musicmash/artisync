override RELEASE="$(git tag -l --points-at HEAD)"
override COMMIT="$(shell git rev-parse --short HEAD)"
override BUILD_TIME="$(shell date -u '+%Y-%m-%dT%H:%M:%S')"

all:

build:
	go build -ldflags="-s -w" -v -o dist/artisync-api ./cmd/artisync-api/...

test t:
	go test -v ./internal/...

run:
	go run ./cmd/artisync-api/...

image:
	docker build \
		--build-arg RELEASE=${RELEASE} \
		--build-arg COMMIT=${COMMIT} \
		--build-arg BUILD_TIME=${BUILD_TIME} \
		-t "musicmash/artisync-api:latest" .
