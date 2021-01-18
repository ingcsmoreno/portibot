FROM alpine:latest as alpine

RUN apk --no-cache add tzdata zip ca-certificates

WORKDIR /usr/share/zoneinfo

# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -q -r -0 /zoneinfo.zip .

FROM golang:1.15-alpine as builder

ARG SHA1VER

COPY cmd/tweet-bot-r /go/src

WORKDIR /go/src

RUN CGO_ENABLED=0 go build -ldflags="-w -s -X main.sha1ver=${SHA1VER} -X main.buildTime=`date +'%Y-%m-%d_%T'` -X main.version=v0.2.0" -o /go/bin/tweet-bot

FROM scratch

# the timezone data:
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /

# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/bin/tweet-bot /go/bin/tweet-bot

ENTRYPOINT ["/go/bin/tweet-bot"]