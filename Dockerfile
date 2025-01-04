FROM alpine:3.21

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk update

RUN apk add ca-certificates flux kubectl --no-cache
