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