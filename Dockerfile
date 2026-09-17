# syntax=docker/dockerfile:1

ARG GO_VERSION="1.27"
FROM golang:${GO_VERSION} AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/server ./cmd

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --gid 10001 appuser \
    && useradd --uid 10001 --gid 10001 --home /nonexistent --shell /usr/sbin/nologin appuser

COPY --from=build /bin/server /bin/server

EXPOSE 8080
USER appuser

ENTRYPOINT ["/bin/server"]