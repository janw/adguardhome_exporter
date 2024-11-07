FROM golang:1.22-alpine AS build
ARG VERSION=dev
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN \
     GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
     go build -o adguardhome_exporter \
     -ldflags "-X main.version=${VERSION}"

FROM alpine:3

RUN apk add -U --no-cache \
     ca-certificates \
     tini

COPY --from=build /src/adguardhome_exporter /usr/bin/

EXPOSE 9311

ENTRYPOINT [ "tini", "--", "adguardhome_exporter" ]
