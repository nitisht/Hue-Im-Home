FROM alpine:latest

RUN apk add nmap

ADD hue-im-home_amd64 /hue-im-home

RUN chmod +x /hue-im-home && mkdir /config

ENTRYPOINT ["/hue-im-home"]