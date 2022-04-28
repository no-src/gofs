FROM golang:latest AS build

WORKDIR /gowork

ARG GOPROXY

COPY . .

RUN echo $(git rev-parse main) >version/commit

RUN go build -v -tags netgo -o . ./...

FROM alpine:3.15.4

WORKDIR /app

ENV PATH=$PATH:/app

EXPOSE 80 443 8105

COPY --from=build /gowork/gofs .

WORKDIR /workspace