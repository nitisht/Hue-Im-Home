FROM golang:1.12.0 AS builder
WORKDIR /go/src/github.com/se1exin/hue-im-home
COPY ./ .
RUN make build-go-arm
# While we're here in amd64, download the qemu-arm-static binary for the arm image in the next build step
RUN curl -L https://github.com/balena-io/qemu/releases/download/v3.0.0%2Bresin/qemu-3.0.0+resin-arm.tar.gz | tar zxvf - -C . && mv qemu-3.0.0+resin-arm/qemu-arm-static .

FROM arm32v7/alpine:latest
# Copy across the qemu binary that was downloaded in the previous build step
COPY --from=builder /go/src/github.com/se1exin/hue-im-home/qemu-arm-static /usr/bin
RUN apk add nmap
COPY --from=builder /go/src/github.com/se1exin/hue-im-home/hue-im-home_arm /hue-im-home
RUN chmod +x /hue-im-home && mkdir /config
ENTRYPOINT ["/hue-im-home"]