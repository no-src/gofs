FROM alpine:3.15.4

WORKDIR /app

ENV PATH=$PATH:/app

COPY ./gofs .