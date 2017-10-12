FROM golang:1.9.1-alpine as builder
RUN apk add --no-cache git
WORKDIR /go/src/github.com/JonPulsifer/go-check-dns
COPY . .
RUN go install -v

FROM alpine:3.6
LABEL maintainer "Jonathan Pulsifer <jonathan@pulsifer.ca>"

RUN addgroup -S go-check-dns && adduser -S -G go-check-dns go-check-dns \
 && apk add --no-cache tini

COPY --from=builder /go/bin/go-check-dns /usr/bin/go-check-dns

USER go-check-dns
ENTRYPOINT ["/usr/bin/go-check-dns"]
