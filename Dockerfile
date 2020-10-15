FROM golang:latest as builder

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

ARG RELEASE=unset
ARG COMMIT=unset
ARG BUILD_TIME=unset
ENV PROJECT=/go/src/github.com/artisync/artisync/internal

WORKDIR /go/src/github.com/artisync
COPY migrations /var/artisync/migrations
COPY cmd cmd
COPY internal internal

RUN go build -v -a \
    -installsuffix cgo \
    -gcflags "all=-trimpath=$(GOPATH)" \
    -ldflags '-linkmode external -extldflags "-static" -s -w \
       -X ${PROJECT}/version.Release=${RELEASE} \
       -X ${PROJECT}/version.Commit=${COMMIT} \
       -X ${PROJECT}/version.BuildTime=${BUILD_TIME}"' \
    -o /usr/local/bin/artisync-api ./cmd/artisync-api/...

FROM alpine:latest

RUN addgroup -S artisync-api && adduser -S artisync-api -G artisync-api
USER artisync-api
WORKDIR /home/artisync-api

COPY --from=builder --chown=artisync-api:artisync-api /var/artisync/migrations /var/artisync/migrations
COPY --from=builder --chown=artisync-api:artisync-api /usr/local/bin/artisync-api /usr/local/bin/artisync-api

ENTRYPOINT ["/usr/local/bin/artisync-api"]
CMD []