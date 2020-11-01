FROM golang:1-alpine as artisync

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ARG RELEASE=unset
ARG COMMIT=unset
ARG BUILD_TIME=unset
ENV PROJECT=github.com/musicmash/artisync

WORKDIR /go/src/github.com/artisync
COPY migrations /var/artisync/migrations
COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd cmd
COPY internal internal

RUN go build -v -a \
    -gcflags "all=-trimpath=${WORKDIR}" \
    -ldflags "-w -s \
       -X ${PROJECT}/internal/version.Release=${RELEASE} \
       -X ${PROJECT}/internal/version.Commit=${COMMIT} \
       -X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}" \
    -o /usr/local/bin/artisync-api ./cmd/artisync-api/...

RUN go build -v -a \
    -gcflags "all=-trimpath=${WORKDIR}" \
    -ldflags "-w -s \
       -X ${PROJECT}/internal/version.Release=${RELEASE} \
       -X ${PROJECT}/internal/version.Commit=${COMMIT} \
       -X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}" \
    -o /usr/local/bin/artisync-daily ./cmd/artisync-daily/...

RUN go build -v -a \
    -gcflags "all=-trimpath=${WORKDIR}" \
    -ldflags "-w -s \
       -X ${PROJECT}/internal/version.Release=${RELEASE} \
       -X ${PROJECT}/internal/version.Commit=${COMMIT} \
       -X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}" \
    -o /usr/local/bin/artisync-sync ./cmd/artisync-sync/...

FROM alpine:latest as artisync-api

RUN addgroup -S artisync-api && adduser -S artisync-api -G artisync-api
USER artisync-api
WORKDIR /home/artisync-api

COPY --from=musicmash/artisync-builder --chown=artisync-api:artisync-api /var/artisync/migrations /var/artisync/migrations
COPY --from=musicmash/artisync-builder --chown=artisync-api:artisync-api /usr/local/bin/artisync-api /usr/local/bin/artisync-api

ENTRYPOINT ["/usr/local/bin/artisync-api"]
CMD []

FROM alpine:latest as artisync-daily

RUN addgroup -S artisync-daily && adduser -S artisync-daily -G artisync-daily
USER artisync-daily
WORKDIR /home/artisync-daily

COPY --from=musicmash/artisync-builder --chown=artisync-daily:artisync-daily /var/artisync/migrations /var/artisync/migrations
COPY --from=musicmash/artisync-builder --chown=artisync-daily:artisync-daily /usr/local/bin/artisync-daily /usr/local/bin/artisync-daily

ENTRYPOINT ["/usr/local/bin/artisync-daily"]
CMD []

FROM alpine:latest as artisync-sync

RUN addgroup -S artisync-sync && adduser -S artisync-sync -G artisync-sync
USER artisync-sync
WORKDIR /home/artisync-sync

COPY --from=musicmash/artisync-builder --chown=artisync-sync:artisync-sync /var/artisync/migrations /var/artisync/migrations
COPY --from=musicmash/artisync-builder --chown=artisync-sync:artisync-sync /usr/local/bin/artisync-sync /usr/local/bin/artisync-sync

ENTRYPOINT ["/usr/local/bin/artisync-sync"]
CMD []