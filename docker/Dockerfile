FROM golang:1.12.0 AS builder
WORKDIR /go/src/github.com/se1exin/hue-im-home
COPY ./ .
RUN make build-go-linux

FROM alpine:latest
RUN apk add nmap
COPY --from=builder /go/src/github.com/se1exin/hue-im-home/hue-im-home_amd64 /hue-im-home
RUN chmod +x /hue-im-home && mkdir /config
ENTRYPOINT ["/hue-im-home"]