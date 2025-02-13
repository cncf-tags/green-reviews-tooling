# Dockerfile is used to build an image for the dagger container with tools that
# are needed for the pipeline. 
FROM alpine:3.21

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk update

RUN apk add ca-certificates flux kubectl --no-cache
