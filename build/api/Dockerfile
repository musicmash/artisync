FROM alpine:latest as artisync-api

RUN addgroup -S artisync-api && adduser -S artisync-api -G artisync-api
USER artisync-api
WORKDIR /home/artisync-api

COPY --from=musicmash/artisync-builder --chown=artisync-api:artisync-api /var/artisync/migrations /var/artisync/migrations
COPY --from=musicmash/artisync-builder --chown=artisync-api:artisync-api /usr/local/bin/artisync-api /usr/local/bin/artisync-api

ENTRYPOINT ["/usr/local/bin/artisync-api"]
CMD []