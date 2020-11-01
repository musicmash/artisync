FROM alpine:latest as artisync-sync

RUN addgroup -S artisync-sync && adduser -S artisync-sync -G artisync-sync
USER artisync-sync
WORKDIR /home/artisync-sync

COPY --from=musicmash/artisync-builder --chown=artisync-sync:artisync-sync /var/artisync/migrations /var/artisync/migrations
COPY --from=musicmash/artisync-builder --chown=artisync-sync:artisync-sync /usr/local/bin/artisync-sync /usr/local/bin/artisync-sync

ENTRYPOINT ["/usr/local/bin/artisync-sync"]
CMD []