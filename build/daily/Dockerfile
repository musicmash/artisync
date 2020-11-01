FROM alpine:latest as artisync-daily

RUN addgroup -S artisync-daily && adduser -S artisync-daily -G artisync-daily
USER artisync-daily
WORKDIR /home/artisync-daily

COPY --from=musicmash/artisync-builder --chown=artisync-daily:artisync-daily /var/artisync/migrations /var/artisync/migrations
COPY --from=musicmash/artisync-builder --chown=artisync-daily:artisync-daily /usr/local/bin/artisync-daily /usr/local/bin/artisync-daily

ENTRYPOINT ["/usr/local/bin/artisync-daily"]
CMD []