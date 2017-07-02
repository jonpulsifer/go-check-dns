FROM alpine:3.6
LABEL maintainer "Jonathan Pulsifer <jonathan@pulsifer.ca>"

RUN addgroup -S go-check-dns && adduser -S -G go-check-dns go-check-dns \
 && apk add --no-cache tini

COPY go-check-dns /usr/bin/go-check-dns
COPY crontab /var/spool/cron/crontabs/go-check-dns

# need root for now, drop when kubernetes cronjobs are available
# USER go-check-dns
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/sbin/crond", "-d", "7", "-f"]
