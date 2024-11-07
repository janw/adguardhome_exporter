FROM golang:1.23-alpine AS build
ARG VERSION=dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o adguardhome_exporter

FROM alpine:3

RUN apk add -U --no-cache \
     ca-certificates \
     tini

COPY --from=build /src/adguardhome_exporter /usr/bin/

EXPOSE 9311

ENTRYPOINT [ "tini", "--", "adguardhome_exporter" ]
