override RELEASE="$(git tag -l --points-at HEAD)"
override COMMIT="$(shell git rev-parse --short HEAD)"
override BUILD_TIME="$(shell date -u '+%Y-%m-%dT%H:%M:%S')"
override VERSION=v1

all:

build:
	go build -ldflags="-s -w" -v -o dist/artisync-api ./cmd/artisync-api/...

install:
	go install -v ./cmd/...

